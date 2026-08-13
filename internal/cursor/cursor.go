package cursor

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"syscall"

	"github.com/nyxar77/caelestia-extras/internal/config"
	"github.com/nyxar77/caelestia-extras/internal/scheme"
	"github.com/pelletier/go-toml/v2"
)

type (
	buildFile struct {
		Cursors map[string]cursorSpec `toml:"cursors"`
	}
	cursorSpec struct {
		PNG         string   `toml:"png"`
		X11Name     string   `toml:"x11_name"`
		X11Sizes    []int    `toml:"x11_sizes"`
		XHotspot    *float64 `toml:"x_hotspot"`
		YHotspot    *float64 `toml:"y_hotspot"`
		X11Delay    *int     `toml:"x11_delay"`
		X11Symlinks []string `toml:"x11_symlinks"`
	}
)

var viewbox = regexp.MustCompile(`viewBox="0 0 ([0-9.]+) ([0-9.]+)"`)

func Sync(cursor config.Cursor, schemeFile string) error {
	if err := os.MkdirAll(cursor.IconDir, 0o755); err != nil {
		return err
	}
	lock, err := lock(cursor)
	if err != nil {
		return err
	}
	defer lock.Close()

	for range 8 {
		active, err := scheme.Read(schemeFile)
		if err != nil {
			return err
		}
		work, err := os.MkdirTemp(cursor.IconDir, "."+cursor.Theme+".")
		if err != nil {
			return err
		}
		defer os.RemoveAll(work)
		theme := filepath.Join(work, "theme")
		if err := buildTheme(cursor, active, theme); err != nil {
			return err
		}
		compiled := filepath.Join(work, "compiled")
		if err := os.MkdirAll(compiled, 0o755); err != nil {
			return err
		}
		if err := runQuiet("hyprcursor-util", "--create", theme, "--output", compiled); err != nil {
			return err
		}
		generated := filepath.Join(compiled, "theme_"+cursor.Theme)
		if _, err := os.Stat(filepath.Join(generated, "manifest.hl")); err != nil {
			return fmt.Errorf("hyprcursor-util did not create %s", cursor.Theme)
		}

		if err := syscall.Flock(int(lock.Fd()), syscall.LOCK_EX); err != nil {
			return err
		}
		current, err := scheme.Read(schemeFile)
		if err != nil {
			syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)
			return err
		}
		if current != active {
			syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)
			continue
		}
		if err := install(cursor, generated); err != nil {
			syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)
			return err
		}
		if cursor.UpdateGTK {
			if err := updateGTK(cursor); err != nil {
				syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)
				return err
			}
		}
		if err := refreshHyprland(cursor); err != nil {
			syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)
			return err
		}
		syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)
		if cursor.XCursorFallback {
			_ = exec.Command("systemctl", "--user", "start", "--no-block", "caelestia-extras-xcursor.service").Run()
		}
		return nil
	}
	return errors.New("scheme changed too often; cursor was not applied")
}

func SyncXCursor(cursor config.Cursor, schemeFile string) error {
	active, err := scheme.Read(schemeFile)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(cursor.IconDir, 0o755); err != nil {
		return err
	}
	work, err := os.MkdirTemp(cursor.IconDir, "."+cursor.Theme+".xcursor.")
	if err != nil {
		return err
	}
	defer os.RemoveAll(work)
	outline := "111111"
	if active.Mode == "dark" {
		outline = "ffffff"
	}
	bitmaps := filepath.Join(work, "bitmaps", cursor.Theme)
	if err := runQuiet("nice", "-n", "19", "cbmp", "-d", cursor.Source, "-o", bitmaps, "-bc", "#"+active.Colours.Primary, "-oc", "#"+outline, "-wc", "#000000"); err != nil {
		return err
	}
	x11 := filepath.Join(work, "x11")
	sizes := make([]string, len(cursor.XCursorSizes))
	for index, size := range cursor.XCursorSizes {
		sizes[index] = strconv.Itoa(size)
	}
	args := append([]string{"-n", "19", "ctgen", cursor.BuildConfig, "-d", bitmaps, "-s"}, sizes...)
	args = append(args, "-p", "x11", "-o", x11, "-n", cursor.Theme, "-c", "Caelestia dynamic cursor")
	if err := runQuiet("nice", args...); err != nil {
		return err
	}

	lock, err := lock(cursor)
	if err != nil {
		return err
	}
	defer lock.Close()
	if err := syscall.Flock(int(lock.Fd()), syscall.LOCK_EX); err != nil {
		return err
	}
	defer syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)
	current, err := scheme.Read(schemeFile)
	if err != nil {
		return err
	}
	if current != active {
		return nil
	}
	target := filepath.Join(cursor.IconDir, cursor.Theme)
	if _, err := os.Stat(filepath.Join(target, "hyprcursors")); err != nil {
		return nil
	}
	return replaceDirectory(filepath.Join(x11, cursor.Theme, "cursors"), filepath.Join(target, "cursors"))
}

func buildTheme(cursor config.Cursor, active scheme.Scheme, output string) error {
	data, err := os.ReadFile(cursor.BuildConfig)
	if err != nil {
		return err
	}
	var build buildFile
	if err := toml.Unmarshal(data, &build); err != nil {
		return fmt.Errorf("parse cursor build config: %w", err)
	}
	fallback, ok := build.Cursors["fallback_settings"]
	if !ok {
		return errors.New("cursor build config has no fallback_settings")
	}
	if err := os.MkdirAll(output, 0o755); err != nil {
		return err
	}
	manifest := "name = " + cursor.Theme + "\ndescription = Caelestia dynamic cursor\nversion = 0.1\ncursors_directory = hyprcursors\n"
	if err := os.WriteFile(filepath.Join(output, "manifest.hl"), []byte(manifest), 0o644); err != nil {
		return err
	}
	index := "[Icon Theme]\nName=" + cursor.Theme + "\nComment=Caelestia dynamic cursor\nInherits=hicolor\n"
	if err := os.WriteFile(filepath.Join(output, "index.theme"), []byte(index), 0o644); err != nil {
		return err
	}

	names := make([]string, 0, len(build.Cursors))
	for name := range build.Cursors {
		if name != "fallback_settings" {
			names = append(names, name)
		}
	}
	slices.Sort(names)
	for _, key := range names {
		spec := build.Cursors[key]
		if spec.X11Name == "" || spec.PNG == "" {
			return fmt.Errorf("invalid cursor %s", key)
		}
		pattern := strings.TrimSuffix(spec.PNG, ".png") + ".svg"
		sources, err := filepath.Glob(filepath.Join(cursor.Source, pattern))
		if err != nil {
			return err
		}
		if len(sources) == 0 {
			sources, err = filepath.Glob(filepath.Join(cursor.Source, spec.X11Name, pattern))
			if err != nil {
				return err
			}
		}
		if len(sources) == 0 {
			return fmt.Errorf("no SVG source for %s", spec.X11Name)
		}
		slices.Sort(sources)
		first, err := os.ReadFile(sources[0])
		if err != nil {
			return err
		}
		match := viewbox.FindStringSubmatch(string(first))
		if len(match) != 3 {
			return fmt.Errorf("no viewBox in %s", sources[0])
		}
		width, _ := strconv.ParseFloat(match[1], 64)
		height, _ := strconv.ParseFloat(match[2], 64)
		x := value(spec.XHotspot, fallback.XHotspot)
		y := value(spec.YHotspot, fallback.YHotspot)
		delay := intValue(spec.X11Delay, fallback.X11Delay)
		shape := filepath.Join(output, "hyprcursors", spec.X11Name)
		if err := os.MkdirAll(shape, 0o755); err != nil {
			return err
		}
		meta := []string{"resize_algorithm = none", fmt.Sprintf("hotspot_x = %.8f", x/width), fmt.Sprintf("hotspot_y = %.8f", y/height), ""}
		for _, source := range sources {
			svg, err := os.ReadFile(source)
			if err != nil {
				return err
			}
			recolored := recolor(string(svg), active)
			filename := filepath.Base(source)
			if err := os.WriteFile(filepath.Join(shape, filename), []byte(recolored), 0o644); err != nil {
				return err
			}
			meta = append(meta, fmt.Sprintf("define_size = 24, %s, %d", filename, delay))
		}
		for _, alias := range spec.X11Symlinks {
			meta = append(meta, "define_override = "+alias)
		}
		if err := os.WriteFile(filepath.Join(shape, "meta.hl"), []byte(strings.Join(meta, "\n")+"\n"), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func value(value, fallback *float64) float64 {
	if value != nil {
		return *value
	}
	if fallback != nil {
		return *fallback
	}
	return 0
}

func intValue(value, fallback *int) int {
	if value != nil {
		return *value
	}
	if fallback != nil {
		return *fallback
	}
	return 0
}

func recolor(svg string, active scheme.Scheme) string {
	outline := "#111111"
	if active.Mode == "dark" {
		outline = "#ffffff"
	}
	return strings.NewReplacer(
		"#00FF00", "#"+strings.ToLower(active.Colours.Primary),
		"#0000FF", outline,
		"#FF0000", "#000000",
	).Replace(svg)
}

func lock(cursor config.Cursor) (*os.File, error) {
	return os.OpenFile(filepath.Join(cursor.IconDir, "."+cursor.Theme+".lock"), os.O_CREATE|os.O_RDWR, 0o600)
}

func install(cursor config.Cursor, generated string) error {
	target := filepath.Join(cursor.IconDir, cursor.Theme)
	old := filepath.Join(cursor.IconDir, "."+cursor.Theme+".old."+strconv.Itoa(os.Getpid()))
	if err := os.RemoveAll(old); err != nil {
		return err
	}
	if _, err := os.Stat(target); err == nil {
		if err := os.Rename(target, old); err != nil {
			return err
		}
		if _, err := os.Stat(filepath.Join(old, "cursors")); err == nil {
			if err := copyDirectory(filepath.Join(old, "cursors"), filepath.Join(generated, "cursors")); err != nil {
				return err
			}
		}
	}
	if err := os.Rename(generated, target); err != nil {
		return err
	}
	return os.RemoveAll(old)
}

func replaceDirectory(source, target string) error {
	next := target + ".new." + strconv.Itoa(os.Getpid())
	old := target + ".old." + strconv.Itoa(os.Getpid())
	if err := os.RemoveAll(next); err != nil {
		return err
	}
	if err := os.RemoveAll(old); err != nil {
		return err
	}
	if err := os.Rename(source, next); err != nil {
		return err
	}
	if _, err := os.Stat(target); err == nil {
		if err := os.Rename(target, old); err != nil {
			return err
		}
	}
	if err := os.Rename(next, target); err != nil {
		return err
	}
	return os.RemoveAll(old)
}

func copyDirectory(source, target string) error {
	return filepath.WalkDir(source, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		destination := filepath.Join(target, relative)
		if entry.IsDir() {
			return os.MkdirAll(destination, 0o755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(destination, data, 0o644)
	})
}

func updateGTK(cursor config.Cursor) error {
	if err := runQuiet("dconf", "write", "/org/gnome/desktop/interface/cursor-theme", "'"+cursor.Theme+"'"); err != nil {
		return err
	}
	return runQuiet("dconf", "write", "/org/gnome/desktop/interface/cursor-size", strconv.Itoa(cursor.Size))
}

func refreshHyprland(cursor config.Cursor) error {
	if err := runQuiet("hyprctl", "setcursor", cursor.Theme, strconv.Itoa(cursor.Size)); err != nil {
		return err
	}
	return runQuiet("hyprctl", "eval", "hl.dsp.force_renderer_reload()")
}

func runQuiet(name string, arguments ...string) error {
	command := exec.Command(name, arguments...)
	command.Stdout = nil
	command.Stderr = nil
	if err := command.Run(); err != nil {
		return fmt.Errorf("run %s: %w", name, err)
	}
	return nil
}
