package integration

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

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

func TestWatchedFilesDeduplicatesSharedThemeOutputs(t *testing.T) {
	configuration := config.Config{
		Scheme: config.Scheme{File: "/state/caelestia/scheme.json"},
		GTK:    &config.GTK{},
		Qt:     &config.Qt{ThemeDir: "/state/caelestia/theme"},
		Portal: &config.Portal{ThemeDir: "/state/caelestia/theme"},
	}
	files := watchedFiles(configuration)
	count := 0
	for _, file := range files {
		if file == "/state/caelestia/theme/qt-caelestia.conf" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("shared Qt palette watched %d times: %#v", count, files)
	}
}

func TestWatchFilesReportsAtomicReplacement(t *testing.T) {
	directory := t.TempDir()
	target := filepath.Join(directory, "theme.conf")
	write(t, target, "old")
	ctx, cancel := context.WithCancel(context.Background())
	events, failures, closeWatcher, err := watchFiles(ctx, []string{target})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		cancel()
		closeWatcher()
	}()
	temporary := filepath.Join(directory, ".theme.conf.new")
	write(t, temporary, "new")
	if err := os.Rename(temporary, target); err != nil {
		t.Fatal(err)
	}
	select {
	case <-events:
	case err := <-failures:
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("atomic replacement did not produce a watcher event")
	}
}

func TestSyncPortalCopiesGeneratedFiles(t *testing.T) {
	root := t.TempDir()
	theme := filepath.Join(root, "theme")
	configHome := filepath.Join(root, "config")
	dataHome := filepath.Join(root, "data")
	write(t, filepath.Join(theme, "gtk-portal.css"), "portal")
	write(t, filepath.Join(theme, "qt-caelestia.conf"), "palette")
	write(t, filepath.Join(theme, "qt6ct-portal.qss"), "stylesheet")

	if err := SyncPortal(config.Portal{ThemeDir: theme, ConfigHome: configHome, DataHome: dataHome, ThemeName: "Caelestia-Portal"}); err != nil {
		t.Fatal(err)
	}

	checks := map[string]string{
		filepath.Join(dataHome, "themes", "Caelestia-Portal", "gtk-3.0", "gtk.css"): "portal",
		filepath.Join(dataHome, "themes", "Caelestia-Portal", "gtk-4.0", "gtk.css"): "portal",
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

func TestSyncPortalNeverRemovesSharedQtFiles(t *testing.T) {
	root := t.TempDir()
	theme := filepath.Join(root, "theme")
	configHome := filepath.Join(root, "config")
	dataHome := filepath.Join(root, "data")
	write(t, filepath.Join(theme, "qt-caelestia.conf"), "palette")
	write(t, filepath.Join(theme, "qt6ct-portal.qss"), "portal stylesheet")
	sharedPalette := filepath.Join(configHome, "qt6ct", "colors", "caelestia.conf")
	sharedStylesheet := filepath.Join(configHome, "qt6ct", "qss", "caelestia.qss")
	write(t, sharedPalette, "palette")
	write(t, sharedStylesheet, "shared stylesheet")

	if err := SyncPortal(config.Portal{ThemeDir: theme, ConfigHome: configHome, DataHome: dataHome, ThemeName: "Caelestia-Portal"}); err != nil {
		t.Fatal(err)
	}
	for path, want := range map[string]string{
		sharedPalette:    "palette",
		sharedStylesheet: "shared stylesheet",
	} {
		data, err := os.ReadFile(path)
		if err != nil || string(data) != want {
			t.Fatalf("shared Qt file %s = %q, %v; want %q", path, data, err, want)
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

func TestSyncPortalDoesNotTouchApplicationGTKStyles(t *testing.T) {
	root := t.TempDir()
	theme := filepath.Join(root, "theme")
	configHome := filepath.Join(root, "config")
	dataHome := filepath.Join(root, "data")
	write(t, filepath.Join(theme, "gtk-portal.css"), "portal")
	write(t, filepath.Join(configHome, "gtk-3.0", "gtk.css"), "/* managed elsewhere */\n")
	write(t, filepath.Join(configHome, "gtk-4.0", "gtk.css"), "/* user stylesheet */\n")

	if err := SyncPortal(config.Portal{ThemeDir: theme, ConfigHome: configHome, DataHome: dataHome, ThemeName: "Caelestia-Portal"}); err != nil {
		t.Fatal(err)
	}
	for path, want := range map[string]string{
		filepath.Join(configHome, "gtk-3.0", "gtk.css"): "/* managed elsewhere */\n",
		filepath.Join(configHome, "gtk-4.0", "gtk.css"): "/* user stylesheet */\n",
	} {
		data, err := os.ReadFile(path)
		if err != nil || string(data) != want {
			t.Fatalf("GTK stylesheet %s = %q, %v; want %q", path, data, err, want)
		}
	}
}
