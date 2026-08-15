package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExecuteHelpDoesNotRequireConfig(t *testing.T) {
	var output bytes.Buffer
	if err := execute([]string{"help", "cursor"}, &output, &output); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "Usage: caelestia-extras cursor <action>") {
		t.Fatalf("unexpected help output: %s", output.String())
	}
}

func TestExecuteSyncHelpDoesNotRequireConfig(t *testing.T) {
	var output bytes.Buffer
	if err := execute([]string{"help", "sync"}, &output, &output); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "Applies all enabled integrations concurrently") {
		t.Fatalf("unexpected help output: %s", output.String())
	}
}

func TestExecuteFlagHelpSucceeds(t *testing.T) {
	var output bytes.Buffer
	if err := execute([]string{"--help"}, &output, &output); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "Usage:") {
		t.Fatalf("unexpected help output: %s", output.String())
	}
}

func TestExecuteRejectsUnknownCommandWithUsageError(t *testing.T) {
	err := execute([]string{"not-a-command"}, &bytes.Buffer{}, &bytes.Buffer{})
	var usageErr *usageError
	if !errors.As(err, &usageErr) {
		t.Fatalf("expected usage error, got %T: %v", err, err)
	}
	if !strings.Contains(err.Error(), "unknown command") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecuteValidatesCommandBeforeLoadingConfig(t *testing.T) {
	err := execute([]string{"cursor", "nope"}, &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "unknown cursor action") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecuteWritesCompletion(t *testing.T) {
	var output bytes.Buffer
	if err := execute([]string{"completion", "fish"}, &output, &output); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "complete -c caelestia-extras") {
		t.Fatalf("unexpected completion output: %s", output.String())
	}
}

func TestExecuteRejectsUnknownCompletionShell(t *testing.T) {
	err := execute([]string{"completion", "powershell"}, &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "unsupported shell") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecuteConfigValidate(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte("[hyprtoolkit]\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := execute([]string{"--config", path, "config", "validate"}, &output, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "Configuration is valid: "+path) {
		t.Fatalf("validation did not confirm success: %q", output.String())
	}
	if strings.Contains(output.String(), "\x1b[") {
		t.Fatalf("redirected output contains terminal colours: %q", output.String())
	}
}

func TestExecuteConfigValidateRejectsUnsupportedBackend(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	contents := "[compositor]\nbackend = \"niri\"\n\n[hyprtoolkit]\n"
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
	err := execute([]string{"--config", path, "config", "validate"}, &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "unsupported compositor backend") {
		t.Fatalf("unexpected error: %v", err)
	}
}
