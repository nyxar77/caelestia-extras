package compositor

import (
	"fmt"
	"os/exec"
	"strconv"
)

// Hyprland applies generated cursor themes through hyprctl.
type Hyprland struct{}

func (Hyprland) Name() string { return "hyprland" }

func (Hyprland) ApplyCursor(cursor Cursor) error {
	// Bump size briefly so Hyprland invalidates its cached cursor surface and
	// re-uploads the active cursor shape immediately without waiting for hover.
	if err := run("hyprctl", "setcursor", cursor.Theme, strconv.Itoa(cursor.Size+1)); err != nil {
		return err
	}
	return run("hyprctl", "setcursor", cursor.Theme, strconv.Itoa(cursor.Size))
}

func (Hyprland) Validate() error {
	if _, err := exec.LookPath("hyprctl"); err != nil {
		return fmt.Errorf("required command \"hyprctl\" was not found in PATH")
	}
	return nil
}

func run(name string, arguments ...string) error {
	command := exec.Command(name, arguments...)
	command.Stdout = nil
	command.Stderr = nil
	if err := command.Run(); err != nil {
		return fmt.Errorf("run %s: %w", name, err)
	}
	return nil
}
