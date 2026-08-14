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

func TestSyncQBittorrentIgnoresMissingResourceCompiler(t *testing.T) {
	root := t.TempDir()
	bin := filepath.Join(root, "bin")
	if err := os.MkdirAll(bin, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(bin, "qbittorrent"), []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin)
	configFile := filepath.Join(root, "config", "qBittorrent.conf")
	if err := SyncQBittorrent(config.QBittorrent{Command: "qbittorrent", RCCCommand: "missing-rcc", ConfigFile: configFile}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(configFile); !os.IsNotExist(err) {
		t.Fatalf("config file = %v, want no file", err)
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
		"rcc":         "#!/bin/sh\nset -eu\nout=\nqrc=\nwhile [ \"$#\" -gt 0 ]; do\n  if [ \"$1\" = -o ]; then out=$2; shift; else qrc=$1; fi\n  shift\ndone\n: > \"$QRC_CAPTURE\"\nwhile IFS= read -r line; do printf '%s\\n' \"$line\" >> \"$QRC_CAPTURE\"; done < \"$qrc\"\nprintf bundle > \"$out\"\n",
	} {
		path := filepath.Join(bin, name)
		if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("PATH", bin)
	qrcCapture := filepath.Join(root, "resources.qrc")
	t.Setenv("QRC_CAPTURE", qrcCapture)
	themeDir := filepath.Join(root, "theme")
	write(t, filepath.Join(themeDir, "qbittorrent.qss"), "QWidget {}")
	write(t, filepath.Join(themeDir, "qbittorrent.json"), `{ "colors": { "Palette.Window": "#091613" } }`)
	themeFile := filepath.Join(root, "data", "qBittorrent", "themes", "Caelestia.qbtheme")
	configFile := filepath.Join(root, "config", "qBittorrent.conf")
	if err := SyncQBittorrent(config.QBittorrent{Command: "qbittorrent", RCCCommand: "rcc", ThemeDir: themeDir, ThemeFile: themeFile, ConfigFile: configFile}); err != nil {
		t.Fatal(err)
	}
	if data, err := os.ReadFile(themeFile); err != nil || string(data) != "bundle" {
		t.Fatalf("theme = %q, %v", data, err)
	}
	configData, err := os.ReadFile(configFile)
	if err != nil || !strings.Contains(string(configData), "General\\UseCustomUITheme=true") || !strings.Contains(string(configData), "General\\CustomUIThemePath="+themeFile) || !strings.Contains(string(configData), "Advanced\\useSystemIconTheme=false") {
		t.Fatalf("config = %q, %v", configData, err)
	}
	qrcData, err := os.ReadFile(qrcCapture)
	if err != nil || !strings.Contains(string(qrcData), "icons/downloading.svg") || !strings.Contains(string(qrcData), "icons/dark/downloading.svg") || !strings.Contains(string(qrcData), "icons/dark/tracker-error.svg") {
		t.Fatalf("resources = %q, %v", qrcData, err)
	}
}

func TestQBittorrentStatusColoursStaySemanticInBothModes(t *testing.T) {
	for surface, want := range map[string]map[string]string{
		"#091613": {"success": "#9fc8a4", "warning": "#c9b983", "error": "#d79a94"},
		"#f8f8f8": {"success": "#3e7758", "warning": "#7f6a38", "error": "#9b504c"},
	} {
		got, err := qbittorrentStatusColours(surface)
		if err != nil {
			t.Fatal(err)
		}
		for name, colour := range want {
			if got[name] != colour {
				t.Fatalf("%s on %s = %s, want %s", name, surface, got[name], colour)
			}
		}
	}
}
