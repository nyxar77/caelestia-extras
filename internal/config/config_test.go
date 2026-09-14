package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestSchemaIsValidJSON(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "config", "caelestia-extras.schema.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !json.Valid(data) {
		t.Fatal("configuration schema is not valid JSON")
	}
	var schema struct {
		Properties map[string]schemaNode `json:"properties"`
	}
	if err := json.Unmarshal(data, &schema); err != nil {
		t.Fatal(err)
	}
	assertSchemaMatches(t, reflect.TypeOf(Config{}), schema.Properties, "")
}

type schemaNode struct {
	Properties map[string]schemaNode `json:"properties"`
}

func assertSchemaMatches(t *testing.T, configType reflect.Type, properties map[string]schemaNode, path string) {
	t.Helper()
	if configType.Kind() == reflect.Pointer {
		configType = configType.Elem()
	}
	fields := make(map[string]reflect.StructField)
	for index := range configType.NumField() {
		field := configType.Field(index)
		fields[field.Tag.Get("toml")] = field
	}
	for name, field := range fields {
		node, ok := properties[name]
		if !ok {
			t.Errorf("schema is missing %s%s", path, name)
			continue
		}
		fieldType := field.Type
		if fieldType.Kind() == reflect.Pointer {
			fieldType = fieldType.Elem()
		}
		if fieldType.Kind() == reflect.Struct {
			assertSchemaMatches(t, fieldType, node.Properties, path+name+".")
		}
	}
	for name := range properties {
		if _, ok := fields[name]; !ok {
			t.Errorf("schema contains unknown property %s%s", path, name)
		}
	}
}

func TestLoadDefaultsToHyprland(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte("[gtk]\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	config, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if config.Compositor.Backend != "hyprland" {
		t.Fatalf("backend = %q", config.Compositor.Backend)
	}
}

func TestManualStarterConfigLoads(t *testing.T) {
	path := filepath.Join("..", "..", "assets", "manual", "config.toml")
	configuration, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if configuration.Compositor.Backend != "hyprland" {
		t.Fatalf("backend = %q", configuration.Compositor.Backend)
	}
}

func TestLoadSetsQtDataHome(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", filepath.Join(t.TempDir(), "data"))
	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte("[qt]\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	config, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if config.Qt.DataHome != os.Getenv("XDG_DATA_HOME") {
		t.Fatalf("data home = %q", config.Qt.DataHome)
	}
	if config.Qt.WidgetStyle != "Breeze" {
		t.Fatalf("widget style = %q", config.Qt.WidgetStyle)
	}
	if config.Qt.IconTheme != "Papirus-Dark" {
		t.Fatalf("icon theme = %q", config.Qt.IconTheme)
	}
}

func TestLoadSetsDefaultQBittorrentConfig(t *testing.T) {
	configHome := filepath.Join(t.TempDir(), "config")
	t.Setenv("XDG_CONFIG_HOME", configHome)
	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte("[qbittorrent]\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	config, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(configHome, "qBittorrent", "qBittorrent.conf")
	if config.QBittorrent.ConfigFile != want {
		t.Fatalf("config file = %q, want %q", config.QBittorrent.ConfigFile, want)
	}
}

func TestLoadSetsDefaultBloomThemeDir(t *testing.T) {
	configHome := filepath.Join(t.TempDir(), "config")
	t.Setenv("XDG_CONFIG_HOME", configHome)
	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte("[bloom]\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	configuration, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(configHome, "spicetify", "Themes", "Bloom")
	if configuration.Bloom.ThemeDir != want {
		t.Fatalf("theme dir = %q, want %q", configuration.Bloom.ThemeDir, want)
	}
}

func TestLoadRejectsUnknownFields(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte("[gtk]\ndark_thme = \"typo\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil || !strings.Contains(err.Error(), "dark_thme") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadRejectsNonPositiveXCursorSize(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	contents := "[cursor]\nsource = \"/tmp/source\"\nbuild_config = \"/tmp/build.toml\"\nxcursor_sizes = [24, 0]\n"
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil || !strings.Contains(err.Error(), "xcursor_sizes") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateAcceptsGeneratedOutputThatDoesNotExistYet(t *testing.T) {
	config := Config{Hyprtoolkit: &Hyprtoolkit{}}
	if err := config.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestValidateReportsMissingRequirements(t *testing.T) {
	config := Config{
		Scheme: Scheme{File: filepath.Join(t.TempDir(), "scheme.json")},
		Cursor: &Cursor{
			Source:      filepath.Join(t.TempDir(), "cursor"),
			BuildConfig: filepath.Join(t.TempDir(), "cursor.toml"),
		},
	}
	err := config.Validate()
	if err == nil {
		t.Fatal("expected validation to fail")
	}
	for _, expected := range []string{"scheme file", "cursor source", "cursor build config", "hyprcursor-util"} {
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("validation error does not mention %q: %v", expected, err)
		}
	}
}

func TestValidateReportsMissingQtPlatformThemes(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	err := Config{Qt: &Qt{}}.Validate()
	if err == nil || !strings.Contains(err.Error(), "qt5ct") || !strings.Contains(err.Error(), "qt6ct") {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestValidateReportsMissingPortalServiceManager(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	err := Config{Portal: &Portal{}}.Validate()
	if err == nil || !strings.Contains(err.Error(), "systemctl") {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestValidateBloomRequiresScheme(t *testing.T) {
	err := Config{
		Scheme: Scheme{File: filepath.Join(t.TempDir(), "missing-scheme.json")},
		Bloom:  &Bloom{},
	}.Validate()
	if err == nil || !strings.Contains(err.Error(), "scheme file") || strings.Contains(err.Error(), "no integrations") {
		t.Fatalf("unexpected validation error: %v", err)
	}
}
