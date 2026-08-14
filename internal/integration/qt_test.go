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
	write(t, filepath.Join(theme, "qt-caelestia.conf"), "palette")
	write(t, filepath.Join(theme, "qt-caelestia.qss"), "stylesheet")

	if err := SyncQt(config.Qt{ThemeDir: theme, ConfigHome: configHome}); err != nil {
		t.Fatal(err)
	}

	for _, version := range []string{"qt5ct", "qt6ct"} {
		for path, want := range map[string]string{
			filepath.Join(configHome, version, "colors", "caelestia.conf"): "palette",
			filepath.Join(configHome, version, "qss", "caelestia.qss"):     "stylesheet",
		} {
			got, err := os.ReadFile(path)
			if err != nil || string(got) != want {
				t.Fatalf("%s = %q, %v; want %q", path, got, err, want)
			}
		}
	}
}
