package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nyxar77/caelestia-extras/internal/config"
)

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestSyncPortalCopiesGeneratedFiles(t *testing.T) {
	root := t.TempDir()
	theme := filepath.Join(root, "theme")
	configHome := filepath.Join(root, "config")
	dataHome := filepath.Join(root, "data")
	write(t, filepath.Join(theme, "gtk-portal.css"), "portal")
	write(t, filepath.Join(theme, "gtk-global.css"), "global")
	write(t, filepath.Join(theme, "qt6ct-caelestia.conf"), "palette")
	write(t, filepath.Join(theme, "qt6ct-portal.qss"), "stylesheet")

	if err := SyncPortal(config.Portal{ThemeDir: theme, ConfigHome: configHome, DataHome: dataHome, ThemeName: "Caelestia-Portal", ApplyGlobalGTK: true}); err != nil {
		t.Fatal(err)
	}

	checks := map[string]string{
		filepath.Join(dataHome, "themes", "Caelestia-Portal", "gtk-3.0", "gtk.css"): "portal",
		filepath.Join(dataHome, "themes", "Caelestia-Portal", "gtk-4.0", "gtk.css"): "portal",
		filepath.Join(configHome, "gtk-3.0", "gtk.css"):                             "global",
		filepath.Join(configHome, "gtk-4.0", "gtk.css"):                             "global",
		filepath.Join(configHome, "portal-qt", "qt6ct", "colors", "caelestia.conf"): "palette",
		filepath.Join(configHome, "portal-qt", "qt6ct", "qss", "caelestia.qss"):     "stylesheet",
	}
	for path, expected := range checks {
		value, err := os.ReadFile(path)
		if err != nil || string(value) != expected {
			t.Fatalf("%s = %q, %v", path, value, err)
		}
	}
}

func TestSyncPortalRemovesStalePortalOverride(t *testing.T) {
	root := t.TempDir()
	theme := filepath.Join(root, "theme")
	configHome := filepath.Join(root, "config")
	dataHome := filepath.Join(root, "data")
	write(t, filepath.Join(theme, "gtk.css"), "generic")
	for _, version := range []string{"gtk-3.0", "gtk-4.0"} {
		write(t, filepath.Join(dataHome, "themes", "Caelestia-Portal", version, "gtk.css"), "generic")
	}

	if err := SyncPortal(config.Portal{ThemeDir: theme, ConfigHome: configHome, DataHome: dataHome, ThemeName: "Caelestia-Portal"}); err != nil {
		t.Fatal(err)
	}
	for _, version := range []string{"gtk-3.0", "gtk-4.0"} {
		portalGTK := filepath.Join(dataHome, "themes", "Caelestia-Portal", version, "gtk.css")
		if _, err := os.Stat(portalGTK); !os.IsNotExist(err) {
			t.Fatalf("stale portal override still exists: %v", err)
		}
	}
}

func TestSyncPortalRestoresGlobalGTKThemeWithoutRemovingUserStyles(t *testing.T) {
	root := t.TempDir()
	theme := filepath.Join(root, "theme")
	configHome := filepath.Join(root, "config")
	dataHome := filepath.Join(root, "data")
	write(t, filepath.Join(theme, "gtk-portal.css"), "portal")
	write(t, filepath.Join(configHome, "gtk-3.0", "gtk.css"), managedGlobalGTKHeader+"\nmanaged")
	write(t, filepath.Join(configHome, "gtk-4.0", "gtk.css"), "/* user stylesheet */\n")

	if err := SyncPortal(config.Portal{ThemeDir: theme, ConfigHome: configHome, DataHome: dataHome, ThemeName: "Caelestia-Portal"}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(configHome, "gtk-3.0", "gtk.css")); !os.IsNotExist(err) {
		t.Fatalf("managed GTK stylesheet = %v, want removed", err)
	}
	data, err := os.ReadFile(filepath.Join(configHome, "gtk-4.0", "gtk.css"))
	if err != nil || string(data) != "/* user stylesheet */\n" {
		t.Fatalf("user GTK stylesheet = %q, %v", data, err)
	}
}
