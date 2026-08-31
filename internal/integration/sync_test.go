package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nyxar77/caelestia-extras/internal/config"
)

func TestSyncAllAppliesIndependentIntegrations(t *testing.T) {
	root := t.TempDir()
	theme := filepath.Join(root, "theme")
	configHome := filepath.Join(root, "config")
	dataHome := filepath.Join(root, "data")
	write(t, filepath.Join(theme, "hyprtoolkit.conf"), "accent = blue\n")
	write(t, filepath.Join(theme, "qt-caelestia.conf"), "palette")
	write(t, filepath.Join(theme, "breeze-caelestia.colors"), "colours")
	write(t, filepath.Join(theme, "gtk-portal.css"), "portal")

	configuration := config.Config{
		Hyprtoolkit: &config.Hyprtoolkit{
			ThemeDir:   theme,
			ConfigFile: filepath.Join(configHome, "hypr", "hyprtoolkit.conf"),
		},
		Qt: &config.Qt{
			ThemeDir:   theme,
			ConfigHome: configHome,
			DataHome:   dataHome,
		},
		QBittorrent: &config.QBittorrent{
			ConfigFile: filepath.Join(configHome, "qBittorrent", "qBittorrent.conf"),
		},
		Portal: &config.Portal{
			ThemeDir:   theme,
			ConfigHome: configHome,
			DataHome:   dataHome,
			ThemeName:  "Caelestia-Portal",
		},
	}

	if err := SyncAll(configuration, false, false); err != nil {
		t.Fatal(err)
	}

	checks := map[string]string{
		filepath.Join(configHome, "hypr", "hyprtoolkit.conf"):                       "accent = blue\n",
		filepath.Join(configHome, "qt5ct", "colors", "caelestia.conf"):              "palette",
		filepath.Join(configHome, "portal-qt", "qt6ct", "colors", "caelestia.conf"): "palette",
		filepath.Join(dataHome, "color-schemes", "Caelestia.colors"):                "colours",
		filepath.Join(dataHome, "themes", "Caelestia-Portal", "gtk-3.0", "gtk.css"): "portal",
		filepath.Join(dataHome, "themes", "Caelestia-Portal", "gtk-4.0", "gtk.css"): "portal",
	}
	for path, expected := range checks {
		value, err := os.ReadFile(path)
		if err != nil || string(value) != expected {
			t.Errorf("%s = %q, %v; want %q", path, value, err, expected)
		}
	}
	qbittorrent, err := os.ReadFile(configuration.QBittorrent.ConfigFile)
	if err != nil || !strings.Contains(string(qbittorrent), "General\\UseCustomUITheme=false") {
		t.Errorf("qBittorrent config = %q, %v", qbittorrent, err)
	}
}

func TestSyncAllReportsEveryFailedIntegration(t *testing.T) {
	root := t.TempDir()
	err := SyncAll(config.Config{
		Hyprtoolkit: &config.Hyprtoolkit{
			ThemeDir:   filepath.Join(root, "missing-hyprtoolkit"),
			ConfigFile: filepath.Join(root, "hyprtoolkit.conf"),
		},
		Qt: &config.Qt{
			ThemeDir:   filepath.Join(root, "missing-qt"),
			ConfigHome: filepath.Join(root, "config"),
			DataHome:   filepath.Join(root, "data"),
		},
	}, false, false)
	if err == nil {
		t.Fatal("expected aggregate sync to fail")
	}
	for _, integration := range []string{"Hyprtoolkit", "Qt"} {
		if !strings.Contains(err.Error(), integration) {
			t.Errorf("aggregate error does not mention %s: %v", integration, err)
		}
	}
}
