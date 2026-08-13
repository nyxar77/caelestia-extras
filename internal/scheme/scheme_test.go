package scheme

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRead(t *testing.T) {
	path := filepath.Join(t.TempDir(), "scheme.json")
	if err := os.WriteFile(path, []byte(`{"colours":{"primary":"8EC07C"},"mode":"dark"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	value, err := Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if value.Colours.Primary != "8EC07C" || value.Mode != "dark" {
		t.Fatalf("unexpected scheme: %#v", value)
	}
}

func TestReadRejectsInvalidColour(t *testing.T) {
	path := filepath.Join(t.TempDir(), "scheme.json")
	if err := os.WriteFile(path, []byte(`{"colours":{"primary":"green"},"mode":"dark"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Read(path); err == nil {
		t.Fatal("expected invalid colour to fail")
	}
}
