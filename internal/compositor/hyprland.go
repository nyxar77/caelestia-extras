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
	if err := run("hyprctl", "setcursor", cursor.Theme, strconv.Itoa(cursor.Size)); err != nil {
		return err
	}
	return run("hyprctl", "eval", "hl.dsp.force_renderer_reload()")
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
