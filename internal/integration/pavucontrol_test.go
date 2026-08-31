package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nyxar77/caelestia-extras/internal/config"
)

func TestLaunchPavucontrolReportsMissingCommand(t *testing.T) {
	err := LaunchPavucontrol(config.Pavucontrol{Command: "caelestia-extras-test-missing-pavucontrol"}, nil)
	if err == nil || !strings.Contains(err.Error(), "find command") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSyncQBittorrentSelectsNativeQtTheme(t *testing.T) {
	root := t.TempDir()
	configFile := filepath.Join(root, "qBittorrent", "qBittorrent.conf")
	write(t, configFile, "[Appearance]\nColorScheme=Dark\nStyle=Fusion\n\n[Preferences]\nAdvanced\\useSystemIconTheme=false\nGeneral\\CustomUIThemePath=/old/theme.qbtheme\nGeneral\\UseCustomUITheme=true\nGeneral\\Locale=en\n")

	if err := SyncQBittorrent(config.QBittorrent{ConfigFile: configFile}); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatal(err)
	}
	for _, setting := range []string{
		"ColorScheme=System",
		"Style=system",
		"Advanced\\useSystemIconTheme=true",
		"General\\UseCustomUITheme=false",
		"General\\CustomUIThemePath=/old/theme.qbtheme",
		"General\\Locale=en",
	} {
		if !strings.Contains(string(data), setting) {
			t.Fatalf("config does not contain %q:\n%s", setting, data)
		}
	}
}

func TestSyncQBittorrentCreatesMinimalConfig(t *testing.T) {
	configFile := filepath.Join(t.TempDir(), "qBittorrent", "qBittorrent.conf")
	if err := SyncQBittorrent(config.QBittorrent{ConfigFile: configFile}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatal(err)
	}
	want := "[Appearance]\nColorScheme=System\nStyle=system\n\n[Preferences]\nAdvanced\\useSystemIconTheme=true\nGeneral\\UseCustomUITheme=false\n"
	if string(data) != want {
		t.Fatalf("config = %q, want %q", data, want)
	}
}

func TestSyncQBittorrentDoesNotRewriteUnchangedConfig(t *testing.T) {
	configFile := filepath.Join(t.TempDir(), "qBittorrent.conf")
	write(t, configFile, "[Appearance]\nColorScheme=System\nStyle=system\n\n[Preferences]\nAdvanced\\useSystemIconTheme=true\nGeneral\\UseCustomUITheme=false\n")
	before, err := os.Stat(configFile)
	if err != nil {
		t.Fatal(err)
	}
	if err := SyncQBittorrent(config.QBittorrent{ConfigFile: configFile}); err != nil {
		t.Fatal(err)
	}
	after, err := os.Stat(configFile)
	if err != nil {
		t.Fatal(err)
	}
	if !os.SameFile(before, after) {
		t.Fatal("unchanged qBittorrent config was replaced")
	}
}

func TestSyncQBittorrentPreservesConfigPermissions(t *testing.T) {
	configFile := filepath.Join(t.TempDir(), "qBittorrent.conf")
	write(t, configFile, "[Appearance]\nColorScheme=Dark\n")
	if err := os.Chmod(configFile, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := SyncQBittorrent(config.QBittorrent{ConfigFile: configFile}); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(configFile)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("config mode = %o, want 600", info.Mode().Perm())
	}
}

func TestSyncQBittorrentDoesNotReplaceSymlink(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "managed.conf")
	configFile := filepath.Join(root, "qBittorrent.conf")
	write(t, target, "[Appearance]\nColorScheme=Dark\n")
	if err := os.Symlink(target, configFile); err != nil {
		t.Fatal(err)
	}
	err := SyncQBittorrent(config.QBittorrent{ConfigFile: configFile})
	if err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("unexpected error: %v", err)
	}
	info, statErr := os.Lstat(configFile)
	if statErr != nil {
		t.Fatal(statErr)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Fatal("qBittorrent config symlink was replaced")
	}
}
