package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Compositor  Compositor   `toml:"compositor"`
	Scheme      Scheme       `toml:"scheme"`
	Cursor      *Cursor      `toml:"cursor"`
	GTK         *GTK         `toml:"gtk"`
	Hyprtoolkit *Hyprtoolkit `toml:"hyprtoolkit"`
	Pavucontrol *Pavucontrol `toml:"pavucontrol"`
	Qt          *Qt          `toml:"qt"`
	QBittorrent *QBittorrent `toml:"qbittorrent"`
	Portal      *Portal      `toml:"portal"`
}

type Compositor struct {
	Backend string `toml:"backend"`
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

type Hyprtoolkit struct {
	ThemeDir   string `toml:"theme_dir"`
	ConfigFile string `toml:"config_file"`
}

type Pavucontrol struct {
	Command string `toml:"command"`
}

type Qt struct {
	ThemeDir   string `toml:"theme_dir"`
	ConfigHome string `toml:"config_home"`
}

type QBittorrent struct {
	Command    string `toml:"command"`
	RCCCommand string `toml:"rcc_command"`
	ThemeDir   string `toml:"theme_dir"`
	ThemeFile  string `toml:"theme_file"`
	ConfigFile string `toml:"config_file"`
}

type Portal struct {
	ThemeDir       string `toml:"theme_dir"`
	ConfigHome     string `toml:"config_home"`
	DataHome       string `toml:"data_home"`
	ThemeName      string `toml:"theme_name"`
	ApplyGlobalGTK bool   `toml:"apply_global_gtk"`
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
	if config.Compositor.Backend == "" {
		config.Compositor.Backend = "hyprland"
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
	if config.Hyprtoolkit != nil {
		if config.Hyprtoolkit.ThemeDir == "" {
			config.Hyprtoolkit.ThemeDir = filepath.Join(xdg("XDG_STATE_HOME", ".local/state"), "caelestia", "theme")
		}
		if config.Hyprtoolkit.ConfigFile == "" {
			config.Hyprtoolkit.ConfigFile = filepath.Join(xdg("XDG_CONFIG_HOME", ".config"), "hypr", "hyprtoolkit.conf")
		}
	}
	if config.Pavucontrol != nil && config.Pavucontrol.Command == "" {
		config.Pavucontrol.Command = "pavucontrol-qt"
	}
	if config.Qt != nil {
		if config.Qt.ThemeDir == "" {
			config.Qt.ThemeDir = filepath.Join(xdg("XDG_STATE_HOME", ".local/state"), "caelestia", "theme")
		}
		if config.Qt.ConfigHome == "" {
			config.Qt.ConfigHome = xdg("XDG_CONFIG_HOME", ".config")
		}
	}
	if config.QBittorrent != nil {
		if config.QBittorrent.Command == "" {
			config.QBittorrent.Command = "qbittorrent"
		}
		if config.QBittorrent.RCCCommand == "" {
			config.QBittorrent.RCCCommand = "rcc"
		}
		if config.QBittorrent.ThemeDir == "" {
			config.QBittorrent.ThemeDir = filepath.Join(xdg("XDG_STATE_HOME", ".local/state"), "caelestia", "theme")
		}
		if config.QBittorrent.ThemeFile == "" {
			config.QBittorrent.ThemeFile = filepath.Join(xdg("XDG_DATA_HOME", ".local/share"), "caelestia-extras", "qbittorrent", "caelestia.qbtheme")
		}
		if config.QBittorrent.ConfigFile == "" {
			config.QBittorrent.ConfigFile = filepath.Join(xdg("XDG_CONFIG_HOME", ".config"), "qBittorrent", "qBittorrent.conf")
		}
	}
	if config.Portal != nil {
		if config.Portal.ThemeDir == "" {
			config.Portal.ThemeDir = filepath.Join(xdg("XDG_STATE_HOME", ".local/state"), "caelestia", "theme")
		}
		if config.Portal.ConfigHome == "" {
			config.Portal.ConfigHome = xdg("XDG_CONFIG_HOME", ".config")
		}
		if config.Portal.DataHome == "" {
			config.Portal.DataHome = xdg("XDG_DATA_HOME", ".local/share")
		}
		if config.Portal.ThemeName == "" {
			config.Portal.ThemeName = "Caelestia-Portal"
		}
	}
	return config, nil
}

// Validate checks the files and external commands needed by the enabled
// integrations. It does not require generated theme output to exist yet.
func (c Config) Validate() error {
	var problems []string
	enabled := 0

	needsScheme := c.Cursor != nil || c.GTK != nil
	if needsScheme {
		if err := regularFile(c.Scheme.File, "scheme file"); err != nil {
			problems = append(problems, err.Error())
		}
	}

	if c.Cursor != nil {
		enabled++
		if err := directory(c.Cursor.Source, "cursor source"); err != nil {
			problems = append(problems, err.Error())
		}
		if err := regularFile(c.Cursor.BuildConfig, "cursor build config"); err != nil {
			problems = append(problems, err.Error())
		}
		for _, command := range []string{"hyprcursor-util"} {
			if err := commandAvailable(command); err != nil {
				problems = append(problems, err.Error())
			}
		}
		if c.Cursor.UpdateGTK {
			if err := commandAvailable("dconf"); err != nil {
				problems = append(problems, err.Error())
			}
		}
		if c.Cursor.XCursorFallback {
			for _, command := range []string{"cbmp", "ctgen"} {
				if err := commandAvailable(command); err != nil {
					problems = append(problems, err.Error())
				}
			}
		}
	}

	if c.GTK != nil {
		enabled++
		if err := commandAvailable("dconf"); err != nil {
			problems = append(problems, err.Error())
		}
	}
	if c.Hyprtoolkit != nil {
		enabled++
	}
	if c.Pavucontrol != nil {
		enabled++
	}
	if c.Qt != nil {
		enabled++
		for _, command := range []string{"qt5ct", "qt6ct"} {
			if err := commandAvailable(command); err != nil {
				problems = append(problems, err.Error())
			}
		}
	}
	if c.QBittorrent != nil {
		enabled++
	}
	if c.Portal != nil {
		enabled++
	}
	if enabled == 0 {
		problems = append(problems, "no integrations are enabled")
	}
	if len(problems) == 0 {
		return nil
	}
	return errors.New("configuration validation failed:\n- " + strings.Join(problems, "\n- "))
}

func regularFile(path, label string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("%s %q is not available: %w", label, path, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("%s %q is not a regular file", label, path)
	}
	return nil
}

func directory(path, label string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("%s %q is not available: %w", label, path, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%s %q is not a directory", label, path)
	}
	return nil
}

func commandAvailable(command string) error {
	if _, err := exec.LookPath(command); err != nil {
		return fmt.Errorf("required command %q was not found in PATH", command)
	}
	return nil
}

func xdg(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, fallback)
}
