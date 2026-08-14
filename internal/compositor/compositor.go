package compositor

import "fmt"

// Cursor is the compositor-facing part of a generated cursor theme.
type Cursor struct {
	Theme string
	Size  int
}

// Backend owns compositor-specific actions. Theme generation and installation
// stay outside this package so they can be shared by every backend.
type Backend interface {
	Name() string
	ApplyCursor(Cursor) error
	Validate() error
}

// New returns the configured compositor backend.
func New(name string) (Backend, error) {
	switch name {
	case "hyprland":
		return Hyprland{}, nil
	default:
		return nil, fmt.Errorf("unsupported compositor backend %q", name)
	}
}
