package integration

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	command := []string{pavucontrol.Command, "-style", "Fusion"}
	stylesheet := filepath.Join(stateHome, "caelestia", "theme", "pavucontrol-qt.qss")
	if _, err := os.Stat(stylesheet); err == nil {
		command = append(command, "-stylesheet", stylesheet)
	}
	return syscallExec(append(command, arguments...))
}

func SyncQt(qt config.Qt) error {
	assets := map[string]string{
		"qt-caelestia.conf": "colors/caelestia.conf",
		"qt-caelestia.qss":  "qss/caelestia.qss",
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
	return nil
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
	if err := copyFile(colors, filepath.Join(work, "config.json")); err != nil {
		return err
	}
	resources := "<RCC>\n  <qresource prefix=\"/\">\n    <file>stylesheet.qss</file>\n    <file>config.json</file>\n  </qresource>\n</RCC>\n"
	qrc := filepath.Join(work, "resources.qrc")
	if err := os.WriteFile(qrc, []byte(resources), 0o644); err != nil {
		return err
	}
	compiled := filepath.Join(work, "caelestia.qbtheme")
	if output, err := exec.Command(rcc, "-binary", "-o", compiled, qrc).CombinedOutput(); err != nil {
		return fmt.Errorf("compile qBittorrent theme: %w: %s", err, strings.TrimSpace(string(output)))
	}
	if err := os.Rename(compiled, qbittorrent.ThemeFile); err != nil {
		return err
	}
	return writeQBittorrentPreferences(qbittorrent.ConfigFile, qbittorrent.ThemeFile)
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

	if portal.ApplyGlobalGTK && exists(filepath.Join(theme, "gtk-global.css")) {
		for _, version := range []string{"gtk-3.0", "gtk-4.0"} {
			destination := filepath.Join(portal.ConfigHome, version, "gtk.css")
			if err := copyFile(filepath.Join(theme, "gtk-global.css"), destination); err != nil {
				return err
			}
			if err := os.Remove(filepath.Join(portal.ConfigHome, version, "thunar.css")); err != nil && !os.IsNotExist(err) {
				return err
			}
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
