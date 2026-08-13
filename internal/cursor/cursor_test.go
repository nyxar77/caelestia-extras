package cursor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nyxar77/caelestia-extras/internal/config"
	"github.com/nyxar77/caelestia-extras/internal/scheme"
)

func TestBuildTheme(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "modern")
	if err := os.MkdirAll(filepath.Join(source, "wait"), 0o755); err != nil {
		t.Fatal(err)
	}
	svg := `<svg viewBox="0 0 256 256"><path fill="#00FF00" stroke="#0000FF"/></svg>`
	if err := os.WriteFile(filepath.Join(source, "left_ptr.svg"), []byte(svg), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "wait", "wait-01.svg"), []byte(svg), 0o644); err != nil {
		t.Fatal(err)
	}
	build := filepath.Join(root, "x.build.toml")
	contents := "[cursors.fallback_settings]\n" +
		"x_hotspot = 128\n" +
		"y_hotspot = 128\n" +
		"x11_delay = 40\n\n" +
		"[cursors.left_ptr]\n" +
		"png = 'left_ptr.png'\n" +
		"x_hotspot = 55\n" +
		"y_hotspot = 17\n" +
		"x11_name = 'left_ptr'\n" +
		"x11_symlinks = ['arrow']\n\n" +
		"[cursors.wait]\n" +
		"png = 'wait-*.png'\n" +
		"x11_name = 'wait'\n"
	if err := os.WriteFile(build, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(root, "theme")
	active := scheme.Scheme{Mode: "dark"}
	active.Colours.Primary = "8EC07C"
	if err := buildTheme(config.Cursor{Source: source, BuildConfig: build, Theme: "Test"}, active, output); err != nil {
		t.Fatal(err)
	}
	value, err := os.ReadFile(filepath.Join(output, "hyprcursors", "left_ptr", "left_ptr.svg"))
	if err != nil || !strings.Contains(string(value), "#8ec07c") {
		t.Fatalf("cursor not recoloured: %v %s", err, value)
	}
	if _, err := os.Stat(filepath.Join(output, "hyprcursors", "wait", "wait-01.svg")); err != nil {
		t.Fatal(err)
	}
}
