package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nyxar77/caelestia-extras/internal/config"
)

func TestSyncQtCopiesGeneratedFilesForQt5AndQt6(t *testing.T) {
	root := t.TempDir()
	theme := filepath.Join(root, "theme")
	configHome := filepath.Join(root, "config")
	dataHome := filepath.Join(root, "data")
	write(t, filepath.Join(theme, "qt-caelestia.conf"), "palette")
	write(t, filepath.Join(theme, "breeze-caelestia.colors"), "colours")
	write(t, filepath.Join(configHome, "kdeglobals"), "[General]\nBrowserApplication=firefox\n\n[Other]\nValue=kept\n")

	if err := SyncQt(config.Qt{ThemeDir: theme, ConfigHome: configHome, DataHome: dataHome}); err != nil {
		t.Fatal(err)
	}

	for _, version := range []string{"qt5ct", "qt6ct"} {
		for path, want := range map[string]string{
			filepath.Join(configHome, version, "colors", "caelestia.conf"): "palette",
		} {
			got, err := os.ReadFile(path)
			if err != nil || string(got) != want {
				t.Fatalf("%s = %q, %v; want %q", path, got, err, want)
			}
		}
	}
	for path, want := range map[string]string{
		filepath.Join(dataHome, "color-schemes", "Caelestia.colors"): "colours",
		filepath.Join(configHome, "kdeglobals"):                      "[General]\nBrowserApplication=firefox\nColorScheme=Caelestia\n\n[Other]\nValue=kept\n",
	} {
		got, err := os.ReadFile(path)
		if err != nil || string(got) != want {
			t.Fatalf("%s = %q, %v; want %q", path, got, err, want)
		}
	}
}
