package integration

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"github.com/nyxar77/caelestia-extras/internal/config"
	"github.com/nyxar77/caelestia-extras/internal/scheme"
)

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
		return nil
	}
	return copyFile(source, hyprtoolkit.ConfigFile)
}

func LaunchPavucontrol(pavucontrol config.Pavucontrol, arguments []string) error {
	if _, err := exec.LookPath(pavucontrol.Command); err != nil {
		return nil
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
	colourScheme := filepath.Join(qt.ThemeDir, "breeze-caelestia.colors")
	if exists(colourScheme) {
		if err := copyFile(colourScheme, filepath.Join(qt.DataHome, "color-schemes", "Caelestia.colors")); err != nil {
			return err
		}
		return setKDEColourScheme(filepath.Join(qt.ConfigHome, "kdeglobals"), "Caelestia")
	}
	return nil
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
	if _, err := exec.LookPath(qbittorrent.Command); err != nil {
		return nil
	}
	rcc, err := exec.LookPath(qbittorrent.RCCCommand)
	if err != nil {
		return nil
	}

	stylesheet := filepath.Join(qbittorrent.ThemeDir, "qbittorrent.qss")
	colors := filepath.Join(qbittorrent.ThemeDir, "qbittorrent.json")
	if !exists(stylesheet) || !exists(colors) {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(qbittorrent.ThemeFile), 0o755); err != nil {
		return err
	}
	work, err := os.MkdirTemp(filepath.Dir(qbittorrent.ThemeFile), ".caelestia-qbittorrent.")
	if err != nil {
		return err
	}
	defer os.RemoveAll(work)
	if err := copyFile(stylesheet, filepath.Join(work, "stylesheet.qss")); err != nil {
		return err
	}
	statusColours, err := prepareQBittorrentConfig(colors, filepath.Join(work, "config.json"))
	if err != nil {
		return err
	}
	icons, err := writeQBittorrentSidebarIcons(work, statusColours)
	if err != nil {
		return err
	}
	resources := qbittorrentResources(icons)
	qrc := filepath.Join(work, "resources.qrc")
	if err := os.WriteFile(qrc, []byte(resources), 0o644); err != nil {
		return err
	}
	compiled := filepath.Join(work, "Caelestia.qbtheme")
	if output, err := exec.Command(rcc, "-binary", "-o", compiled, qrc).CombinedOutput(); err != nil {
		return fmt.Errorf("compile qBittorrent theme: %w: %s", err, strings.TrimSpace(string(output)))
	}
	if err := os.Rename(compiled, qbittorrent.ThemeFile); err != nil {
		return err
	}
	return writeQBittorrentPreferences(qbittorrent.ConfigFile, qbittorrent.ThemeFile)
}

type qbittorrentTheme struct {
	Colours map[string]string `json:"colors"`
}

// prepareQBittorrentConfig keeps the Caelestia palette as the base UI theme,
// then adds a small, stable semantic palette for torrent states. Wallpaper
// colours can change freely, but an error should never look like success.
func prepareQBittorrentConfig(source, destination string) (map[string]string, error) {
	data, err := os.ReadFile(source)
	if err != nil {
		return nil, err
	}
	var theme qbittorrentTheme
	if err := json.Unmarshal(data, &theme); err != nil {
		return nil, fmt.Errorf("parse qBittorrent theme colours: %w", err)
	}
	if theme.Colours == nil {
		return nil, errors.New("qBittorrent theme has no colours")
	}
	statusColours, err := qbittorrentStatusColours(theme.Colours["Palette.Window"])
	if err != nil {
		return nil, err
	}
	for key, colour := range qbittorrentSemanticColourMap(statusColours) {
		theme.Colours[key] = colour
	}
	output, err := json.MarshalIndent(theme, "", "  ")
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(destination, append(output, '\n'), 0o644); err != nil {
		return nil, err
	}
	return statusColours, nil
}

func qbittorrentStatusColours(surface string) (map[string]string, error) {
	value := strings.TrimPrefix(surface, "#")
	if len(value) != 6 {
		return nil, fmt.Errorf("qBittorrent Palette.Window must be a six-digit hex colour")
	}
	rgb, err := strconv.ParseUint(value, 16, 32)
	if err != nil {
		return nil, fmt.Errorf("parse qBittorrent Palette.Window: %w", err)
	}
	red, green, blue := float64(rgb>>16), float64((rgb>>8)&0xff), float64(rgb&0xff)
	luminance := (0.2126*red + 0.7152*green + 0.0722*blue) / 255
	if luminance < 0.5 {
		return map[string]string{
			"success":  "#9fc8a4",
			"activity": "#9fbfd2",
			"warning":  "#c9b983",
			"error":    "#d79a94",
			"muted":    "#92aaa2",
			"disabled": "#829a92",
		}, nil
	}
	return map[string]string{
		"success":  "#3e7758",
		"activity": "#47708b",
		"warning":  "#7f6a38",
		"error":    "#9b504c",
		"muted":    "#536660",
		"disabled": "#61736d",
	}, nil
}

func qbittorrentSemanticColourMap(colours map[string]string) map[string]string {
	success, activity := colours["success"], colours["activity"]
	warning, failure := colours["warning"], colours["error"]
	muted, disabled := colours["muted"], colours["disabled"]
	return map[string]string{
		"Palette.WindowTextDisabled":             muted,
		"Palette.TextDisabled":                   muted,
		"Palette.ToolTipTextDisabled":            muted,
		"Palette.BrightTextDisabled":             muted,
		"Palette.HighlightedTextDisabled":        disabled,
		"Palette.ButtonTextDisabled":             muted,
		"Log.TimeStamp":                          muted,
		"Log.Info":                               activity,
		"Log.Warning":                            warning,
		"Log.Critical":                           failure,
		"Log.BannedPeer":                         failure,
		"TransferList.Downloading":               success,
		"TransferList.StalledDownloading":        warning,
		"TransferList.DownloadingMetadata":       success,
		"TransferList.ForcedDownloadingMetadata": success,
		"TransferList.ForcedDownloading":         success,
		"TransferList.Uploading":                 activity,
		"TransferList.StalledUploading":          warning,
		"TransferList.ForcedUploading":           activity,
		"TransferList.QueuedDownloading":         warning,
		"TransferList.QueuedUploading":           warning,
		"TransferList.CheckingDownloading":       activity,
		"TransferList.CheckingUploading":         activity,
		"TransferList.CheckingResumeData":        activity,
		"TransferList.StoppedDownloading":        muted,
		"TransferList.StoppedUploading":          muted,
		"TransferList.PausedDownloading":         muted,
		"TransferList.PausedUploading":           muted,
		"TransferList.Moving":                    activity,
		"TransferList.MissingFiles":              failure,
		"TransferList.Error":                     failure,
	}
}

func writeQBittorrentSidebarIcons(work string, colours map[string]string) ([]string, error) {
	icons := qbittorrentSidebarIcons()
	resources := make([]string, 0, len(icons)*2)
	for _, directory := range []string{"icons", "icons/dark"} {
		for name, icon := range icons {
			resource := filepath.Join(directory, name+".svg")
			path := filepath.Join(work, resource)
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return nil, err
			}
			if err := os.WriteFile(path, []byte(qbittorrentSVG(icon.path, colours[icon.colour])), 0o644); err != nil {
				return nil, err
			}
			resources = append(resources, filepath.ToSlash(resource))
		}
	}
	sort.Strings(resources)
	return resources, nil
}

type qbittorrentIcon struct {
	path   string
	colour string
}

func qbittorrentSidebarIcons() map[string]qbittorrentIcon {
	return map[string]qbittorrentIcon{
		"filter-all":        {"M4 5h16v2H4zm0 6h16v2H4zm0 6h16v2H4z", "muted"},
		"downloading":       {"M5 19h14v-2H5m9-8h3l-5 5-5-5h3V3h4z", "success"},
		"upload":            {"M5 19h14v-2H5m9-1V7h3l-5-5-5 5h3v9z", "activity"},
		"checked-completed": {"m9 16.2-4.2-4.1L3.4 13.5 9 19 21 7l-1.4-1.4z", "success"},
		"torrent-start":     {"M8 5v14l11-7z", "success"},
		"stopped":           {"M6 6h12v12H6z", "muted"},
		"filter-active":     {"M12 4a8 8 0 1 0 8 8h-2a6 6 0 1 1-6-6zM11 7v6l4 2 .8-1.3-3.3-2V7z", "success"},
		"filter-inactive":   {"M12 4a8 8 0 1 0 8 8h-2a6 6 0 1 1-6-6z", "muted"},
		"filter-stalled":    {"M7 3h10v3l-3 3 3 3v3H7v-3l3-3-3-3zm2 2v.2l3 3.1 3-3.1V5zm3 8-3 3.1v.2h6v-.2z", "warning"},
		"stalledUP":         {"M5 19h14v-2H5m9-1V8h3l-5-5-5 5h3v8zm5-4h2v4h-2z", "warning"},
		"stalledDL":         {"M5 19h14v-2H5m9-8h3l-5 5-5-5h3V3h4zm5 7h2v4h-2z", "warning"},
		"force-recheck":     {"M17.7 6.3A8 8 0 1 0 20 12h-2a6 6 0 1 1-1.8-4.3L13 11h7V4z", "activity"},
		"set-location":      {"M12 2a7 7 0 0 0-7 7c0 5.3 7 13 7 13s7-7.7 7-13a7 7 0 0 0-7-7m0 9.5A2.5 2.5 0 1 1 12 6a2.5 2.5 0 0 1 0 5.5", "activity"},
		"error":             {"M12 2 2 20h20zm1 14h-2v-2h2zm0-4h-2v-4h2z", "error"},
		"view-categories":   {"M3 5h7l2 2h9v12H3zm2 4v8h14V9z", "activity"},
		"tags":              {"M3 12V5h7l9 9-7 7zm4-5a2 2 0 1 0 0 4 2 2 0 0 0 0-4", "activity"},
		"trackers":          {"M12 2a7 7 0 0 0-7 7c0 5.3 7 13 7 13s7-7.7 7-13a7 7 0 0 0-7-7m0 9.5A2.5 2.5 0 1 1 12 6a2.5 2.5 0 0 1 0 5.5", "muted"},
		"tracker-warning":   {"M12 2 2 20h20zm1 14h-2v-2h2zm0-4h-2v-4h2z", "warning"},
		"tracker-error":     {"M12 2 2 20h20zm1 14h-2v-2h2zm0-4h-2v-4h2z", "error"},
		"trackerless":       {"M12 2a7 7 0 0 0-7 7c0 5.3 7 13 7 13s7-7.7 7-13a7 7 0 0 0-7-7m0 9.5A2.5 2.5 0 1 1 12 6a2.5 2.5 0 0 1 0 5.5", "muted"},
	}
}

func qbittorrentSVG(path, colour string) string {
	return fmt.Sprintf("<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 24 24\"><path fill=\"%s\" d=\"%s\"/></svg>\n", colour, path)
}

func qbittorrentResources(icons []string) string {
	var resources strings.Builder
	resources.WriteString("<RCC>\n  <qresource prefix=\"/\">\n    <file>stylesheet.qss</file>\n    <file>config.json</file>\n")
	for _, icon := range icons {
		fmt.Fprintf(&resources, "    <file>%s</file>\n", icon)
	}
	resources.WriteString("  </qresource>\n</RCC>\n")
	return resources.String()
}

func writeQBittorrentPreferences(path, themeFile string) error {
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	lines := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
	if len(lines) == 1 && lines[0] == "" {
		lines = nil
	}
	sectionStart, sectionEnd := -1, len(lines)
	for index, line := range lines {
		if strings.TrimSpace(line) != "[Preferences]" {
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
		sectionStart = len(lines)
		lines = append(lines, "[Preferences]")
		sectionEnd = len(lines)
	}
	pairs := []struct{ key, value string }{
		{"General\\UseCustomUITheme", "true"},
		{"General\\CustomUIThemePath", themeFile},
		{"Advanced\\useSystemIconTheme", "false"},
	}
	for _, pair := range pairs {
		found := false
		for index := sectionStart + 1; index < sectionEnd; index++ {
			if strings.HasPrefix(lines[index], pair.key+"=") {
				lines[index] = pair.key + "=" + pair.value
				found = true
				break
			}
		}
		if !found {
			lines = append(lines[:sectionEnd], append([]string{pair.key + "=" + pair.value}, lines[sectionEnd:]...)...)
			sectionEnd++
		}
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

	for _, version := range []string{"gtk-3.0", "gtk-4.0"} {
		destination := filepath.Join(portal.ConfigHome, version, "gtk.css")
		if portal.ApplyGlobalGTK && exists(filepath.Join(theme, "gtk-global.css")) {
			if err := copyFile(filepath.Join(theme, "gtk-global.css"), destination); err != nil {
				return err
			}
			if err := os.Remove(filepath.Join(portal.ConfigHome, version, "thunar.css")); err != nil && !os.IsNotExist(err) {
				return err
			}
		} else if err := removeManagedGlobalGTK(destination); err != nil {
			return err
		}
	}

	portalQt := filepath.Join(portal.ConfigHome, "portal-qt", "qt6ct")
	if exists(filepath.Join(theme, "qt6ct-caelestia.conf")) {
		destination := filepath.Join(portalQt, "colors", "caelestia.conf")
		if err := copyFile(filepath.Join(theme, "qt6ct-caelestia.conf"), destination); err != nil {
			return err
		}
		global := filepath.Join(portal.ConfigHome, "qt6ct", "colors", "caelestia.conf")
		if sameFile(filepath.Join(theme, "qt6ct-caelestia.conf"), global) {
			if err := os.Remove(global); err != nil && !os.IsNotExist(err) {
				return err
			}
		}
	}

	portalQSS := filepath.Join(portalQt, "qss", "caelestia.qss")
	if exists(filepath.Join(theme, "qt6ct-portal.qss")) {
		if err := copyFile(filepath.Join(theme, "qt6ct-portal.qss"), portalQSS); err != nil {
			return err
		}
		global := filepath.Join(portal.ConfigHome, "qt6ct", "qss", "caelestia.qss")
		if sameFile(filepath.Join(theme, "qt6ct-portal.qss"), global) || sameFile(filepath.Join(theme, "qt6ct-caelestia.qss"), global) {
			if err := os.Remove(global); err != nil && !os.IsNotExist(err) {
				return err
			}
		}
	} else if sameFile(filepath.Join(theme, "qt6ct-caelestia.qss"), portalQSS) {
		if err := os.Remove(portalQSS); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

const managedGlobalGTKHeader = "/*\n  Dynamic GTK colours only."

// removeManagedGlobalGTK restores GTK's normal theme only when the file was
// written by this integration. User stylesheets are left untouched.
func removeManagedGlobalGTK(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if !strings.HasPrefix(string(data), managedGlobalGTKHeader) {
		return nil
	}
	return os.Remove(path)
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
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return err
	}
	return os.WriteFile(destination, data, 0o644)
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
