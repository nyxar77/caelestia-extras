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
		filepath.Join(configHome, "kdeglobals"):                      "[General]\nBrowserApplication=firefox\nColorScheme=Caelestia\n\n[Other]\nValue=kept\n\n[UiSettings]\nColorScheme=Caelestia\n\n[KDE]\nwidgetStyle=Breeze\n\n[Icons]\nTheme=Papirus-Dark\n",
	} {
		got, err := os.ReadFile(path)
		if err != nil || string(got) != want {
			t.Fatalf("%s = %q, %v; want %q", path, got, err, want)
		}
	}
}

func TestSyncQtUpdatesKDEAppearance(t *testing.T) {
	root := t.TempDir()
	theme := filepath.Join(root, "theme")
	configHome := filepath.Join(root, "config")
	write(t, filepath.Join(theme, "qt-caelestia.conf"), "palette")
	write(t, filepath.Join(theme, "breeze-caelestia.colors"), "colours")
	write(t, filepath.Join(configHome, "kdeglobals"), "[General]\nColorScheme=Old\n\n[UiSettings]\nColorScheme=OldId\nOther=kept\n\n[KDE]\nwidgetStyle=Fusion\n\n[Icons]\nTheme=breeze\n")

	if err := SyncQt(config.Qt{
		ThemeDir:    theme,
		ConfigHome:  configHome,
		DataHome:    filepath.Join(root, "data"),
		WidgetStyle: "Breeze",
		IconTheme:   "Papirus-Dark",
	}); err != nil {
		t.Fatal(err)
	}

	want := "[General]\nColorScheme=Caelestia\n\n[UiSettings]\nColorScheme=Caelestia\nOther=kept\n\n[KDE]\nwidgetStyle=Breeze\n\n[Icons]\nTheme=Papirus-Dark\n"
	got, err := os.ReadFile(filepath.Join(configHome, "kdeglobals"))
	if err != nil || string(got) != want {
		t.Fatalf("kdeglobals = %q, %v; want %q", got, err, want)
	}
}
