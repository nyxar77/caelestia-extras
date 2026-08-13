package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nyxar77/caelestia-extras/internal/config"
)

func TestSyncHyprtoolkitCopiesGeneratedConfig(t *testing.T) {
	root := t.TempDir()
	theme := filepath.Join(root, "theme")
	destination := filepath.Join(root, "config", "hypr", "hyprtoolkit.conf")
	write(t, filepath.Join(theme, "hyprtoolkit.conf"), "accent = blue\n")

	if err := SyncHyprtoolkit(config.Hyprtoolkit{ThemeDir: theme, ConfigFile: destination}); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(destination)
	if err != nil || string(content) != "accent = blue\n" {
		t.Fatalf("config = %q, %v", content, err)
	}
}
