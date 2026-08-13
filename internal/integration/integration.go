package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

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

func run(name string, arguments ...string) error {
	if err := exec.Command(name, arguments...).Run(); err != nil {
		return fmt.Errorf("run %s: %w", name, err)
	}
	return nil
}
