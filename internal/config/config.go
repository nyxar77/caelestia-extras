package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Scheme      Scheme       `toml:"scheme"`
	Cursor      *Cursor      `toml:"cursor"`
	GTK         *GTK         `toml:"gtk"`
	Pavucontrol *Pavucontrol `toml:"pavucontrol"`
}

type Scheme struct {
	File string `toml:"file"`
}

type Cursor struct {
	Source          string `toml:"source"`
	BuildConfig     string `toml:"build_config"`
	IconDir         string `toml:"icon_dir"`
	Theme           string `toml:"theme"`
	Size            int    `toml:"size"`
	XCursorSizes    []int  `toml:"xcursor_sizes"`
	XCursorFallback bool   `toml:"xcursor_fallback"`
	UpdateGTK       bool   `toml:"update_gtk"`
}

type GTK struct {
	DarkTheme  string `toml:"dark_theme"`
	LightTheme string `toml:"light_theme"`
}

type Pavucontrol struct {
	Command string `toml:"command"`
}

func DefaultPath() string {
	return filepath.Join(xdg("XDG_CONFIG_HOME", ".config"), "caelestia-extras", "config.toml")
}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}
	var config Config
	if err := toml.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}
	if config.Scheme.File == "" {
		config.Scheme.File = filepath.Join(xdg("XDG_STATE_HOME", ".local/state"), "caelestia", "scheme.json")
	}
	if config.Cursor != nil {
		cursor := config.Cursor
		if cursor.Source == "" || cursor.BuildConfig == "" {
			return Config{}, fmt.Errorf("cursor requires source and build_config")
		}
		if cursor.IconDir == "" {
			cursor.IconDir = filepath.Join(xdg("XDG_DATA_HOME", ".local/share"), "icons")
		}
		if cursor.Theme == "" {
			cursor.Theme = "Bibata-Caelestia"
		}
		if cursor.Size == 0 {
			cursor.Size = 20
		}
		if len(cursor.XCursorSizes) == 0 {
			cursor.XCursorSizes = []int{20, 24, 32}
		}
		if cursor.Size < 1 {
			return Config{}, fmt.Errorf("cursor size must be positive")
		}
	}
	if config.GTK != nil {
		if config.GTK.DarkTheme == "" {
			config.GTK.DarkTheme = "adw-gtk3-dark"
		}
		if config.GTK.LightTheme == "" {
			config.GTK.LightTheme = "adw-gtk3"
		}
	}
	if config.Pavucontrol != nil && config.Pavucontrol.Command == "" {
		config.Pavucontrol.Command = "pavucontrol-qt"
	}
	return config, nil
}

func xdg(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, fallback)
}
