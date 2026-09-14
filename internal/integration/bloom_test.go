package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nyxar77/caelestia-extras/internal/config"
)

func TestSyncBloomWritesMaterialPalette(t *testing.T) {
	root := t.TempDir()
	schemeFile := filepath.Join(root, "scheme.json")
	themeDir := filepath.Join(root, "Bloom")
	schemeJSON := `{"colours":{"primary":"111111","onSurface":"222222","onSurfaceVariant":"333333","surface":"444444","surfaceContainerLowest":"555555","surfaceContainer":"666666","surfaceContainerHigh":"777777","surfaceContainerHighest":"888888","primaryContainer":"999999","surfaceVariant":"AAAAAA","secondary":"BBBBBB","error":"CCCCCC","tertiary":"DDDDDD","shadow":"EEEEEE"},"mode":"dark"}`
	write(t, schemeFile, schemeJSON)

	if err := SyncBloom(config.Bloom{ThemeDir: themeDir}, schemeFile); err != nil {
		t.Fatal(err)
	}

	want := "[Bloom]\ntext = 222222\nsubtext = 333333\nmain = 444444\nsidebar = 555555\nplayer = 666666\ncard = 777777\nmain-elevated = 777777\nhighlight-elevated = 888888\nhighlight = 999999\nselected-row = 111111\nbutton = 111111\nbutton-active = 111111\nbutton-disabled = AAAAAA\nnotification = BBBBBB\nnotification-error = CCCCCC\nmisc = DDDDDD\nshadow = EEEEEE\ntab-active = 888888\n"
	got, err := os.ReadFile(filepath.Join(themeDir, "color.ini"))
	if err != nil || string(got) != want {
		t.Fatalf("color.ini = %q, %v; want %q", got, err, want)
	}
}
