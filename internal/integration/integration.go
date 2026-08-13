package integration

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

func SyncPortal(portal config.Portal) error {
	theme := portal.ThemeDir
	portalGTK := filepath.Join(portal.DataHome, "themes", portal.ThemeName, "gtk-3.0", "gtk.css")
	if exists(filepath.Join(theme, "gtk-portal.css")) {
		if err := copyFile(filepath.Join(theme, "gtk-portal.css"), portalGTK); err != nil {
			return err
		}
	} else if exists(filepath.Join(theme, "gtk.css")) && sameFile(filepath.Join(theme, "gtk.css"), portalGTK) {
		if err := os.Remove(portalGTK); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	gtk4 := filepath.Join(portal.DataHome, "themes", portal.ThemeName, "gtk-4.0", "gtk.css")
	if err := os.Remove(gtk4); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := removeEmpty(filepath.Dir(gtk4)); err != nil {
		return err
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
