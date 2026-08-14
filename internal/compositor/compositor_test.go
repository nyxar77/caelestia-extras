package compositor

import "testing"

func TestNewHyprland(t *testing.T) {
	backend, err := New("hyprland")
	if err != nil {
		t.Fatal(err)
	}
	if backend.Name() != "hyprland" {
		t.Fatalf("backend = %q", backend.Name())
	}
}

func TestNewRejectsUnknownBackend(t *testing.T) {
	if _, err := New("niri"); err == nil {
		t.Fatal("expected unsupported backend to fail")
	}
}
