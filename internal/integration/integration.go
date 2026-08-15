package integration

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/nyxar77/caelestia-extras/internal/compositor"
	"github.com/nyxar77/caelestia-extras/internal/config"
	"github.com/nyxar77/caelestia-extras/internal/cursor"
	"github.com/nyxar77/caelestia-extras/internal/scheme"
)

const (
	settleDelay  = 300 * time.Millisecond
	xcursorDelay = 10 * time.Second
)

// SyncAll applies every enabled integration. Independent integrations run in
// parallel so one slow backend does not hold up the rest of the desktop.
func SyncAll(configuration config.Config, includeXCursor, reloadPortals bool) error {
	type syncJob struct {
		name string
		run  func() error
	}
	var jobs []syncJob
	if configuration.Cursor != nil {
		backend, err := compositor.New(configuration.Compositor.Backend)
		if err != nil {
			return err
		}
		jobs = append(jobs, syncJob{"cursor", func() error {
			return cursor.Sync(*configuration.Cursor, configuration.Scheme.File, backend)
		}})
	}
	if configuration.GTK != nil {
		jobs = append(jobs, syncJob{"GTK", func() error {
			return SyncGTK(*configuration.GTK, configuration.Scheme.File)
		}})
	}
	if configuration.Hyprtoolkit != nil {
		jobs = append(jobs, syncJob{"Hyprtoolkit", func() error { return SyncHyprtoolkit(*configuration.Hyprtoolkit) }})
	}
	if configuration.Qt != nil {
		jobs = append(jobs, syncJob{"Qt", func() error { return SyncQt(*configuration.Qt) }})
	}
	if configuration.QBittorrent != nil {
		jobs = append(jobs, syncJob{"qBittorrent", func() error { return SyncQBittorrent(*configuration.QBittorrent) }})
	}
	if configuration.Portal != nil {
		jobs = append(jobs, syncJob{"portal", func() error {
			if err := SyncPortal(*configuration.Portal); err != nil {
				return err
			}
			if reloadPortals {
				return reloadPortalServices()
			}
			return nil
		}})
	}

	errorsChannel := make(chan error, len(jobs))
	var group sync.WaitGroup
	for _, job := range jobs {
		group.Add(1)
		go func(job syncJob) {
			defer group.Done()
			if err := job.run(); err != nil {
				errorsChannel <- fmt.Errorf("%s: %w", job.name, err)
			}
		}(job)
	}
	group.Wait()
	close(errorsChannel)
	var joined []error
	for err := range errorsChannel {
		joined = append(joined, err)
	}
	if err := errors.Join(joined...); err != nil {
		return err
	}
	if includeXCursor && configuration.Cursor != nil && configuration.Cursor.XCursorFallback {
		return cursor.SyncXCursor(*configuration.Cursor, configuration.Scheme.File)
	}
	return nil
}

func reloadPortalServices() error {
	command := exec.Command(
		"systemctl", "--user", "try-restart", "--no-block",
		"xdg-desktop-portal-gtk.service",
		"xdg-desktop-portal-hyprland.service",
	)
	if output, err := command.CombinedOutput(); err != nil {
		return fmt.Errorf("reload portal services: %w: %s", err, output)
	}
	return nil
}

// Watch collapses a wallpaper burst into one synchronization of the newest
// completed theme. XCursor generation has a separate quiet-period worker
// because it is intentionally much more expensive than palette propagation.
func Watch(configuration config.Config, progress func(string), warn func(error)) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	events, watcherErrors, closeWatcher, err := watchFiles(ctx, watchedFiles(configuration))
	if err != nil {
		return err
	}
	defer closeWatcher()

	syncRequests := make(chan struct{}, 1)
	xcursorRequests := make(chan struct{}, 1)
	workerErrors := make(chan error, 2)
	go syncWorker(ctx, configuration, syncRequests, workerErrors, progress)
	go xcursorWorker(ctx, configuration, xcursorRequests, workerErrors, progress)

	request(syncRequests)
	var settleTimer *time.Timer
	var xcursorTimer *time.Timer
	var settle <-chan time.Time
	var xcursorQuiet <-chan time.Time
	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-watcherErrors:
			if err != nil {
				return err
			}
		case _, open := <-events:
			if !open {
				if ctx.Err() != nil {
					return nil
				}
				return errors.New("theme watcher stopped unexpectedly")
			}
			settleTimer, settle = resetTimer(settleTimer, settleDelay)
			xcursorTimer, xcursorQuiet = resetTimer(xcursorTimer, xcursorDelay)
		case <-settle:
			settle = nil
			request(syncRequests)
		case <-xcursorQuiet:
			xcursorQuiet = nil
			request(xcursorRequests)
		case err := <-workerErrors:
			if err != nil && warn != nil {
				warn(err)
			}
		}
	}
}

func syncWorker(ctx context.Context, configuration config.Config, requests <-chan struct{}, failures chan<- error, progress func(string)) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-requests:
			report(progress, "Applying the newest generated desktop theme")
			err := SyncAll(configuration, false, true)
			if err == nil {
				report(progress, "Desktop theme updated")
			}
			failures <- err
		}
	}
}

func xcursorWorker(ctx context.Context, configuration config.Config, requests <-chan struct{}, failures chan<- error, progress func(string)) {
	if configuration.Cursor == nil || !configuration.Cursor.XCursorFallback {
		return
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-requests:
			report(progress, "Refreshing the XCursor fallback after the wallpaper settled")
			err := cursor.SyncXCursor(*configuration.Cursor, configuration.Scheme.File)
			if err == nil {
				report(progress, "XCursor fallback updated")
			}
			failures <- err
		}
	}
}

func report(progress func(string), message string) {
	if progress != nil {
		progress(message)
	}
}

func request(channel chan<- struct{}) {
	select {
	case channel <- struct{}{}:
	default:
	}
}

func resetTimer(timer *time.Timer, delay time.Duration) (*time.Timer, <-chan time.Time) {
	if timer == nil {
		timer = time.NewTimer(delay)
		return timer, timer.C
	}
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}
	timer.Reset(delay)
	return timer, timer.C
}

func watchedFiles(configuration config.Config) []string {
	seen := make(map[string]bool)
	var files []string
	add := func(path string) {
		path = filepath.Clean(path)
		if path != "." && !seen[path] {
			seen[path] = true
			files = append(files, path)
		}
	}
	if configuration.Cursor != nil || configuration.GTK != nil {
		add(configuration.Scheme.File)
	}
	if configuration.Hyprtoolkit != nil {
		add(filepath.Join(configuration.Hyprtoolkit.ThemeDir, "hyprtoolkit.conf"))
	}
	if configuration.Qt != nil {
		add(filepath.Join(configuration.Qt.ThemeDir, "qt-caelestia.conf"))
		add(filepath.Join(configuration.Qt.ThemeDir, "breeze-caelestia.colors"))
	}
	if configuration.Portal != nil {
		for _, name := range []string{"gtk-portal.css", "qt-caelestia.conf", "qt6ct-portal.qss"} {
			add(filepath.Join(configuration.Portal.ThemeDir, name))
		}
	}
	return files
}

func watchFiles(ctx context.Context, files []string) (<-chan struct{}, <-chan error, func(), error) {
	fd, err := syscall.InotifyInit1(syscall.IN_CLOEXEC)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("initialize theme watcher: %w", err)
	}
	closeWatcher := func() { _ = syscall.Close(fd) }
	watches := make(map[int]string)
	for _, file := range files {
		directory := filepath.Dir(file)
		alreadyWatched := false
		for _, existing := range watches {
			if existing == directory {
				alreadyWatched = true
				break
			}
		}
		if alreadyWatched {
			continue
		}
		watch, addErr := syscall.InotifyAddWatch(fd, directory, syscall.IN_CLOSE_WRITE|syscall.IN_MOVED_TO|syscall.IN_CREATE|syscall.IN_DELETE)
		if addErr != nil {
			closeWatcher()
			return nil, nil, nil, fmt.Errorf("watch theme directory %q: %w", directory, addErr)
		}
		watches[watch] = directory
	}

	wanted := make(map[string]bool, len(files))
	for _, file := range files {
		wanted[filepath.Clean(file)] = true
	}
	events := make(chan struct{}, 1)
	failures := make(chan error, 1)
	go func() {
		defer close(events)
		buffer := make([]byte, 16*1024)
		for {
			length, readErr := syscall.Read(fd, buffer)
			if readErr != nil {
				if ctx.Err() != nil || errors.Is(readErr, syscall.EBADF) {
					return
				}
				failures <- fmt.Errorf("read theme watcher: %w", readErr)
				return
			}
			for offset := 0; offset+syscall.SizeofInotifyEvent <= length; {
				event := (*syscall.InotifyEvent)(unsafe.Pointer(&buffer[offset]))
				nameBytes := buffer[offset+syscall.SizeofInotifyEvent : offset+syscall.SizeofInotifyEvent+int(event.Len)]
				name := string(bytes.TrimRight(nameBytes, "\x00"))
				if wanted[filepath.Join(watches[int(event.Wd)], name)] {
					request(events)
				}
				offset += syscall.SizeofInotifyEvent + int(event.Len)
			}
		}
	}()
	return events, failures, closeWatcher, nil
}

func SyncGTK(gtk config.GTK, schemeFile string) error {
	active, err := scheme.Read(schemeFile)
	if err != nil {
		return err
	}
	preference, theme := "'prefer-light'", gtk.LightTheme
	if active.Mode == "dark" {
		preference, theme = "'prefer-dark'", gtk.DarkTheme
	}
	if err := run("dconf", "write", "/org/gnome/desktop/interface/color-scheme", preference); err != nil {
		return err
	}
	return run("dconf", "write", "/org/gnome/desktop/interface/gtk-theme", "'"+theme+"'")
}

func SyncHyprtoolkit(hyprtoolkit config.Hyprtoolkit) error {
	source := filepath.Join(hyprtoolkit.ThemeDir, "hyprtoolkit.conf")
	if !exists(source) {
		return fmt.Errorf("generated theme file %q is not available", source)
	}
	return copyFile(source, hyprtoolkit.ConfigFile)
}

func LaunchPavucontrol(pavucontrol config.Pavucontrol, arguments []string) error {
	if _, err := exec.LookPath(pavucontrol.Command); err != nil {
		return fmt.Errorf("find command %q: %w", pavucontrol.Command, err)
	}
	stateHome := os.Getenv("XDG_STATE_HOME")
	if stateHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		stateHome = filepath.Join(home, ".local", "state")
	}
	command := []string{pavucontrol.Command, "-style", "Breeze"}
	stylesheet := filepath.Join(stateHome, "caelestia", "theme", "pavucontrol-qt.qss")
	if _, err := os.Stat(stylesheet); err == nil {
		command = append(command, "-stylesheet", stylesheet)
	}
	return syscallExec(append(command, arguments...))
}

func SyncQt(qt config.Qt) error {
	palette := filepath.Join(qt.ThemeDir, "qt-caelestia.conf")
	colourScheme := filepath.Join(qt.ThemeDir, "breeze-caelestia.colors")
	if !exists(palette) {
		return fmt.Errorf("generated Qt palette %q is not available", palette)
	}
	if !exists(colourScheme) {
		return fmt.Errorf("generated Breeze colour scheme %q is not available", colourScheme)
	}
	assets := map[string]string{
		"qt-caelestia.conf": "colors/caelestia.conf",
	}
	for _, version := range []string{"qt5ct", "qt6ct"} {
		for source, destination := range assets {
			path := filepath.Join(qt.ThemeDir, source)
			if !exists(path) {
				continue
			}
			if err := copyFile(path, filepath.Join(qt.ConfigHome, version, destination)); err != nil {
				return err
			}
		}
	}
	if err := copyFile(colourScheme, filepath.Join(qt.DataHome, "color-schemes", "Caelestia.colors")); err != nil {
		return err
	}
	return setKDEColourScheme(filepath.Join(qt.ConfigHome, "kdeglobals"), "Caelestia")
}

// setKDEColourScheme changes only the active colour scheme and keeps unrelated
// KDE settings intact. Breeze uses this file when resolving its colour scheme.
func setKDEColourScheme(path, name string) error {
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	lines := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
	if len(lines) == 1 && lines[0] == "" {
		lines = nil
	}
	start, end := -1, len(lines)
	for index, line := range lines {
		if strings.TrimSpace(line) != "[General]" {
			continue
		}
		start = index
		for next := index + 1; next < len(lines); next++ {
			if strings.HasPrefix(strings.TrimSpace(lines[next]), "[") {
				end = next
				break
			}
		}
		break
	}
	if start == -1 {
		if len(lines) > 0 {
			lines = append(lines, "")
		}
		lines = append(lines, "[General]", "ColorScheme="+name)
	} else {
		updated := false
		for index := start + 1; index < end; index++ {
			if strings.HasPrefix(strings.TrimSpace(lines[index]), "ColorScheme=") {
				lines[index] = "ColorScheme=" + name
				updated = true
				break
			}
		}
		if !updated {
			insertAt := end
			for insertAt > start+1 && strings.TrimSpace(lines[insertAt-1]) == "" {
				insertAt--
			}
			lines = append(lines[:insertAt], append([]string{"ColorScheme=" + name}, lines[insertAt:]...)...)
		}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o644)
}

func SyncQBittorrent(qbittorrent config.QBittorrent) error {
	return writeQBittorrentPreferences(qbittorrent.ConfigFile)
}

// qBittorrent custom themes replace only QPalette::Normal (the active window
// group). On focus-follows-mouse compositors that leaves inactive windows on
// the system palette and causes the whole interface to switch colour. Keep the
// application on its native Qt palette and let qt6ct and Breeze own all states.
func writeQBittorrentPreferences(path string) error {
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	lines := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
	if len(lines) == 1 && lines[0] == "" {
		lines = nil
	}
	settings := []struct{ section, key, value string }{
		{"Appearance", "ColorScheme", "System"},
		{"Appearance", "Style", "system"},
		{"Preferences", "Advanced\\useSystemIconTheme", "true"},
		{"Preferences", "General\\UseCustomUITheme", "false"},
	}
	changed := false
	for _, setting := range settings {
		var settingChanged bool
		lines, settingChanged = setINIValue(lines, setting.section, setting.key, setting.value)
		changed = changed || settingChanged
	}
	if !changed {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	file, err := os.CreateTemp(filepath.Dir(path), ".qBittorrent.conf.")
	if err != nil {
		return err
	}
	temporary := file.Name()
	defer os.Remove(temporary)
	if err := file.Chmod(0o644); err != nil {
		file.Close()
		return err
	}
	if _, err := file.WriteString(strings.Join(lines, "\n") + "\n"); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(temporary, path)
}

func setINIValue(lines []string, section, key, value string) ([]string, bool) {
	sectionStart, sectionEnd := -1, len(lines)
	for index, line := range lines {
		if strings.TrimSpace(line) != "["+section+"]" {
			continue
		}
		sectionStart = index
		for next := index + 1; next < len(lines); next++ {
			if strings.HasPrefix(strings.TrimSpace(lines[next]), "[") {
				sectionEnd = next
				break
			}
		}
		break
	}
	if sectionStart == -1 {
		if len(lines) > 0 {
			lines = append(lines, "")
		}
		return append(lines, "["+section+"]", key+"="+value), true
	}
	for index := sectionStart + 1; index < sectionEnd; index++ {
		if !strings.HasPrefix(lines[index], key+"=") {
			continue
		}
		updated := key + "=" + value
		if lines[index] == updated {
			return lines, false
		}
		lines[index] = updated
		return lines, true
	}
	lines = append(lines[:sectionEnd], append([]string{key + "=" + value}, lines[sectionEnd:]...)...)
	return lines, true
}

func SyncPortal(portal config.Portal) error {
	theme := portal.ThemeDir
	portalCSS := filepath.Join(theme, "gtk-portal.css")
	for _, version := range []string{"gtk-3.0", "gtk-4.0"} {
		destination := filepath.Join(portal.DataHome, "themes", portal.ThemeName, version, "gtk.css")
		if exists(portalCSS) {
			if err := copyFile(portalCSS, destination); err != nil {
				return err
			}
		} else if exists(filepath.Join(theme, "gtk.css")) && sameFile(filepath.Join(theme, "gtk.css"), destination) {
			if err := os.Remove(destination); err != nil && !os.IsNotExist(err) {
				return err
			}
		}
	}

	portalQt := filepath.Join(portal.ConfigHome, "portal-qt", "qt6ct")
	if exists(filepath.Join(theme, "qt-caelestia.conf")) {
		destination := filepath.Join(portalQt, "colors", "caelestia.conf")
		if err := copyFile(filepath.Join(theme, "qt-caelestia.conf"), destination); err != nil {
			return err
		}
	}

	portalQSS := filepath.Join(portalQt, "qss", "caelestia.qss")
	if exists(filepath.Join(theme, "qt6ct-portal.qss")) {
		if err := copyFile(filepath.Join(theme, "qt6ct-portal.qss"), portalQSS); err != nil {
			return err
		}
	}
	return nil
}

func exists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

func sameFile(first, second string) bool {
	left, leftErr := os.ReadFile(first)
	right, rightErr := os.ReadFile(second)
	return leftErr == nil && rightErr == nil && bytes.Equal(left, right)
}

func copyFile(source, destination string) error {
	data, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	if current, readErr := os.ReadFile(destination); readErr == nil && bytes.Equal(current, data) {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return err
	}
	file, err := os.CreateTemp(filepath.Dir(destination), ".caelestia-extras.")
	if err != nil {
		return err
	}
	temporary := file.Name()
	defer os.Remove(temporary)
	if err := file.Chmod(0o644); err != nil {
		file.Close()
		return err
	}
	if _, err := file.Write(data); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(temporary, destination)
}

func removeEmpty(path string) error {
	err := os.Remove(path)
	if err == nil || os.IsNotExist(err) {
		return nil
	}
	if errors.Is(err, syscall.ENOTEMPTY) {
		return nil
	}
	return err
}

func run(name string, arguments ...string) error {
	if err := exec.Command(name, arguments...).Run(); err != nil {
		return fmt.Errorf("run %s: %w", name, err)
	}
	return nil
}
