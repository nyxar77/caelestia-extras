package config

import (
	"encoding/json"
	"os"
	"path/filepath"
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
