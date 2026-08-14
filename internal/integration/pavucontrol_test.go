package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nyxar77/caelestia-extras/internal/config"
)

func TestLaunchPavucontrolIgnoresMissingCommand(t *testing.T) {
	if err := LaunchPavucontrol(config.Pavucontrol{Command: "caelestia-extras-test-missing-pavucontrol"}, nil); err != nil {
		t.Fatal(err)
	}
}

func TestSyncQBittorrentIgnoresMissingCommand(t *testing.T) {
	if err := SyncQBittorrent(config.QBittorrent{Command: "caelestia-extras-test-missing-qbittorrent"}); err != nil {
		t.Fatal(err)
	}
}

func TestSyncQBittorrentBuildsAndSelectsTheme(t *testing.T) {
	root := t.TempDir()
	bin := filepath.Join(root, "bin")
	if err := os.MkdirAll(bin, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, script := range map[string]string{
		"qbittorrent": "#!/bin/sh\nexit 0\n",
		"rcc":         "#!/bin/sh\nset -eu\nout=\nwhile [ \"$#\" -gt 0 ]; do\n  if [ \"$1\" = -o ]; then out=$2; shift; fi\n  shift\ndone\nprintf bundle > \"$out\"\n",
	} {
		path := filepath.Join(bin, name)
		if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("PATH", bin)
	themeDir := filepath.Join(root, "theme")
	write(t, filepath.Join(themeDir, "qbittorrent.qss"), "QWidget {}")
	write(t, filepath.Join(themeDir, "qbittorrent.json"), `{ "colors": {} }`)
	themeFile := filepath.Join(root, "data", "caelestia.qbtheme")
	configFile := filepath.Join(root, "config", "qBittorrent.conf")
	if err := SyncQBittorrent(config.QBittorrent{Command: "qbittorrent", RCCCommand: "rcc", ThemeDir: themeDir, ThemeFile: themeFile, ConfigFile: configFile}); err != nil {
		t.Fatal(err)
	}
	if data, err := os.ReadFile(themeFile); err != nil || string(data) != "bundle" {
		t.Fatalf("theme = %q, %v", data, err)
	}
	configData, err := os.ReadFile(configFile)
	if err != nil || !strings.Contains(string(configData), "General\\UseCustomUITheme=true") || !strings.Contains(string(configData), "General\\CustomUIThemePath="+themeFile) {
		t.Fatalf("config = %q, %v", configData, err)
	}
}
