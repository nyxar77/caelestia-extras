package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/nyxar77/caelestia-extras/internal/compositor"
	"github.com/nyxar77/caelestia-extras/internal/config"
	"github.com/nyxar77/caelestia-extras/internal/cursor"
	"github.com/nyxar77/caelestia-extras/internal/integration"
)

const version = "0.1.0"

type usageError struct {
	err error
}

func (e *usageError) Error() string { return e.err.Error() }
func (e *usageError) Unwrap() error { return e.err }

func main() {
	if err := run(os.Args[1:]); err != nil {
		var usage *usageError
		printError(os.Stderr, err)
		if errors.As(err, &usage) {
			fmt.Fprintln(os.Stderr, "Run 'caelestia-extras --help' for usage.")
		}
		os.Exit(1)
	}
}

func run(arguments []string) error {
	return execute(arguments, os.Stdout, os.Stderr)
}

func execute(arguments []string, stdout, stderr io.Writer) error {
	ui := newTerminalUI(stdout, stderr)
	flags := flag.NewFlagSet("caelestia-extras", flag.ContinueOnError)
	flags.SetOutput(stderr)
	flags.Usage = func() { _, _ = io.WriteString(stderr, generalHelp) }
	configPath := flags.String("config", config.DefaultPath(), "configuration file")
	showVersion := flags.Bool("version", false, "show version")
	if err := flags.Parse(arguments); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return &usageError{err: err}
	}
	if *showVersion {
		_, err := fmt.Fprintln(stdout, version)
		return err
	}

	args := flags.Args()
	if len(args) == 0 {
		return &usageError{err: errors.New("a command is required")}
	}
	if args[0] == "--help" || args[0] == "-h" || args[0] == "help" {
		return writeHelp(args[1:], stdout)
	}
	if len(args) > 1 && (args[1] == "--help" || args[1] == "-h") {
		return writeHelp(args[:1], stdout)
	}
	if args[0] == "version" {
		if len(args) != 1 {
			return usage("version does not accept arguments")
		}
		_, err := fmt.Fprintln(stdout, version)
		return err
	}
	if args[0] == "completion" {
		return writeCompletion(args[1:], stdout)
	}
	if err := validateCommand(args); err != nil {
		return err
	}

	configFile, err := config.Load(*configPath)
	if err != nil {
		return fmt.Errorf("load config %q: %w", *configPath, err)
	}

	switch args[0] {
	case "sync":
		return ui.run("Synchronizing desktop integrations", "Desktop integrations synchronized", "desktop sync", func() error {
			return integration.SyncAll(configFile, true, true)
		})
	case "watch":
		return ui.run("Watching Caelestia theme files", "Theme watcher stopped", "theme watcher", func() error {
			return integration.Watch(configFile, ui.progress, ui.warning)
		})
	case "cursor":
		if configFile.Cursor == nil {
			return fmt.Errorf("cursor is disabled in %q; add a [cursor] section", *configPath)
		}
		backend, err := compositor.New(configFile.Compositor.Backend)
		if err != nil {
			return err
		}
		if len(args) != 2 {
			return usage("cursor expects one of: sync, sync-xcursor")
		}
		switch args[1] {
		case "sync":
			return ui.run("Building and applying the cursor theme", "Cursor theme synchronized", "cursor sync", func() error {
				return cursor.Sync(*configFile.Cursor, configFile.Scheme.File, backend)
			})
		case "sync-xcursor":
			return ui.run("Building the XCursor fallback", "XCursor fallback synchronized", "cursor XCursor sync", func() error {
				return cursor.SyncXCursor(*configFile.Cursor, configFile.Scheme.File)
			})
		default:
			return usage(fmt.Sprintf("unknown cursor action %q; expected sync or sync-xcursor", args[1]))
		}
	case "gtk":
		if configFile.GTK == nil {
			return fmt.Errorf("GTK sync is disabled in %q; add a [gtk] section", *configPath)
		}
		if len(args) != 2 || args[1] != "sync" {
			return usage("gtk expects: sync")
		}
		return ui.run("Synchronizing GTK preferences", "GTK preferences synchronized", "GTK sync", func() error {
			return integration.SyncGTK(*configFile.GTK, configFile.Scheme.File)
		})
	case "hyprtoolkit":
		if configFile.Hyprtoolkit == nil {
			return fmt.Errorf("Hyprtoolkit sync is disabled in %q; add a [hyprtoolkit] section", *configPath)
		}
		if len(args) != 2 || args[1] != "sync" {
			return usage("hyprtoolkit expects: sync")
		}
		return ui.run("Applying the Hyprtoolkit theme", "Hyprtoolkit theme synchronized", "Hyprtoolkit sync", func() error {
			return integration.SyncHyprtoolkit(*configFile.Hyprtoolkit)
		})
	case "pavucontrol":
		if configFile.Pavucontrol == nil {
			return fmt.Errorf("pavucontrol is disabled in %q; add a [pavucontrol] section", *configPath)
		}
		return ui.run("Launching pavucontrol", "pavucontrol exited", "launch pavucontrol", func() error {
			return integration.LaunchPavucontrol(*configFile.Pavucontrol, args[1:])
		})
	case "qt":
		if configFile.Qt == nil {
			return fmt.Errorf("Qt sync is disabled in %q; add a [qt] section", *configPath)
		}
		if len(args) != 2 || args[1] != "sync" {
			return usage("qt expects: sync")
		}
		return ui.run("Synchronizing the Qt palette", "Qt palette synchronized", "Qt sync", func() error {
			return integration.SyncQt(*configFile.Qt)
		})
	case "qbittorrent":
		if configFile.QBittorrent == nil {
			return fmt.Errorf("qBittorrent sync is disabled in %q; add a [qbittorrent] section", *configPath)
		}
		if len(args) != 2 || args[1] != "sync" {
			return usage("qbittorrent expects: sync")
		}
		return ui.run("Selecting native Breeze for qBittorrent", "qBittorrent will use native Breeze on its next start", "qBittorrent sync", func() error {
			return integration.SyncQBittorrent(*configFile.QBittorrent)
		})
	case "config":
		return ui.run("Validating "+*configPath, "Configuration is valid: "+*configPath, "config validation", func() error {
			if err := configFile.Validate(); err != nil {
				return err
			}
			backend, err := compositor.New(configFile.Compositor.Backend)
			if err != nil {
				return err
			}
			if configFile.Cursor != nil {
				return backend.Validate()
			}
			return nil
		})
	case "portal":
		if configFile.Portal == nil {
			return fmt.Errorf("portal sync is disabled in %q; add a [portal] section", *configPath)
		}
		if len(args) != 2 || args[1] != "sync" {
			return usage("portal expects: sync")
		}
		return ui.run("Synchronizing desktop portals", "Desktop portal themes synchronized", "portal sync", func() error {
			return integration.SyncPortal(*configFile.Portal)
		})
	default:
		return usage(fmt.Sprintf("unknown command %q; expected sync, watch, cursor, gtk, hyprtoolkit, pavucontrol, qt, qbittorrent, portal, config, completion, help, or version", args[0]))
	}
}

func validateCommand(args []string) error {
	switch args[0] {
	case "sync", "watch":
		if len(args) != 1 {
			return usage(fmt.Sprintf("%s does not accept arguments", args[0]))
		}
	case "cursor":
		if len(args) != 2 {
			return usage("cursor expects one of: sync, sync-xcursor")
		}
		if args[1] != "sync" && args[1] != "sync-xcursor" {
			return usage(fmt.Sprintf("unknown cursor action %q; expected sync or sync-xcursor", args[1]))
		}
	case "gtk", "hyprtoolkit", "qt", "portal", "qbittorrent":
		if len(args) != 2 || args[1] != "sync" {
			return usage(fmt.Sprintf("%s expects: sync", args[0]))
		}
	case "pavucontrol":
		return nil
	case "config":
		if len(args) != 2 || args[1] != "validate" {
			return usage("config expects: validate")
		}
	default:
		return usage(fmt.Sprintf("unknown command %q; expected sync, watch, cursor, gtk, hyprtoolkit, pavucontrol, qt, qbittorrent, portal, config, completion, help, or version", args[0]))
	}
	return nil
}

type terminalUI struct {
	stdout       io.Writer
	stderr       io.Writer
	stdoutColour bool
	stderrColour bool
}

func newTerminalUI(stdout, stderr io.Writer) terminalUI {
	return terminalUI{
		stdout:       stdout,
		stderr:       stderr,
		stdoutColour: supportsColour(stdout),
		stderrColour: supportsColour(stderr),
	}
}

func (ui terminalUI) run(action, completed, name string, run func() error) error {
	ui.status("\x1b[36m", "→", action+"...")
	if err := run(); err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	ui.status("\x1b[32m", "✓", completed+".")
	return nil
}

func (ui terminalUI) warning(err error) {
	prefix, reset := "", ""
	if ui.stderrColour {
		prefix, reset = "\x1b[33m", "\x1b[0m"
	}
	fmt.Fprintf(ui.stderr, "%s!%s Sync warning: %v\n", prefix, reset, err)
}

func (ui terminalUI) progress(message string) {
	ui.status("\x1b[36m", "→", message+"...")
}

func (ui terminalUI) status(colour, marker, message string) {
	prefix, reset := "", ""
	if ui.stdoutColour {
		prefix, reset = colour, "\x1b[0m"
	}
	fmt.Fprintf(ui.stdout, "%s%s%s %s\n", prefix, marker, reset, message)
}

func printError(stderr io.Writer, err error) {
	prefix, reset := "", ""
	if supportsColour(stderr) {
		prefix, reset = "\x1b[31m", "\x1b[0m"
	}
	fmt.Fprintf(stderr, "%serror:%s %v\n", prefix, reset, err)
}

func supportsColour(writer io.Writer) bool {
	if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		return false
	}
	file, ok := writer.(*os.File)
	if !ok {
		return false
	}
	info, err := file.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

func usage(message string) error {
	return &usageError{err: errors.New(message)}
}

func writeHelp(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		_, err := io.WriteString(stdout, generalHelp)
		return err
	}
	if len(args) != 1 {
		return usage("help accepts at most one command")
	}

	var help string
	switch args[0] {
	case "sync":
		help = syncHelp
	case "watch":
		help = watchHelp
	case "cursor":
		help = cursorHelp
	case "gtk":
		help = gtkHelp
	case "hyprtoolkit":
		help = hyprtoolkitHelp
	case "pavucontrol":
		help = pavucontrolHelp
	case "qt":
		help = qtHelp
	case "qbittorrent":
		help = qbittorrentHelp
	case "portal":
		help = portalHelp
	case "completion":
		help = completionHelp
	case "config":
		help = configHelp
	default:
		return usage(fmt.Sprintf("no help available for %q", args[0]))
	}
	_, err := io.WriteString(stdout, help)
	return err
}

func writeCompletion(args []string, stdout io.Writer) error {
	if len(args) != 1 {
		return usage("completion expects one shell: bash, zsh, or fish")
	}
	var completion string
	switch args[0] {
	case "bash":
		completion = bashCompletion
	case "zsh":
		completion = zshCompletion
	case "fish":
		completion = fishCompletion
	default:
		return usage(fmt.Sprintf("unsupported shell %q; expected bash, zsh, or fish", args[0]))
	}
	_, err := io.WriteString(stdout, completion)
	return err
}

const generalHelp = `caelestia-extras keeps desktop integrations in sync with Caelestia.

Usage:
  caelestia-extras [options] <command> [arguments]

Commands:
  sync                     Apply every enabled integration now
  watch                    Watch and coalesce runtime theme changes
  cursor sync              Build and apply the active Hyprcursor theme
  cursor sync-xcursor      Build the XCursor fallback
  gtk sync                 Sync the GTK theme and colour preference
  hyprtoolkit sync         Apply generated Hyprtoolkit configuration
  pavucontrol [arguments]  Launch the configured volume mixer
  qt sync                  Apply the shared Qt palette and Breeze colour scheme
  qbittorrent sync         Select native Qt/Breeze theming for qBittorrent
  portal sync              Apply generated portal themes
  config validate          Check files, tools, and enabled integrations
  completion <shell>       Print shell completion for bash, zsh, or fish
  help [command]           Show general or command-specific help
  version                  Print the version

Options:
  -config PATH             Use PATH instead of the default config file
  -version                 Print the version
  -h, --help               Show this help

Examples:
  caelestia-extras sync
  caelestia-extras --config ~/.config/caelestia-extras/work.toml gtk sync
  caelestia-extras completion zsh > ~/.zfunc/_caelestia-extras

Configuration:
  Default: $XDG_CONFIG_HOME/caelestia-extras/config.toml
  Help:    https://github.com/nyxar77/caelestia-extras#manual-configuration
`

const syncHelp = `Usage: caelestia-extras sync

Applies all enabled integrations concurrently, including the XCursor fallback,
then reloads the portal backends so their new palette is used immediately.
`

const watchHelp = `Usage: caelestia-extras watch

Watches Caelestia's scheme and generated theme files. Rapid wallpaper changes
are collapsed into one update of the newest completed theme. The expensive
XCursor fallback runs only after wallpaper changes have been quiet for ten
seconds.
`

const cursorHelp = `Usage: caelestia-extras cursor <action>

Builds a cursor theme from the active Caelestia scheme.

Actions:
  sync              Build and apply the Hyprcursor theme
  sync-xcursor      Build the XCursor fallback
`

const gtkHelp = `Usage: caelestia-extras gtk sync

Sets the GTK theme and GNOME colour preference for the active scheme.
Requires dconf and a [gtk] section in the configuration.
`

const hyprtoolkitHelp = `Usage: caelestia-extras hyprtoolkit sync

Copies Caelestia's generated Hyprtoolkit configuration to its active path.
Requires a [hyprtoolkit] section in the configuration.
`

const pavucontrolHelp = `Usage: caelestia-extras pavucontrol [arguments]

Launches the configured pavucontrol-qt command with the generated stylesheet.
Arguments after pavucontrol are passed to pavucontrol-qt.
`

const qtHelp = `Usage: caelestia-extras qt sync

Copies Caelestia's generated Qt palette into qt5ct and qt6ct, then installs
the matching Breeze colour scheme.
Requires a [qt] section in the configuration.
`

const qbittorrentHelp = `Usage: caelestia-extras qbittorrent sync

Disables qBittorrent's custom UI theme so the native Qt platform palette and
widget style remain consistent for active, inactive, and disabled windows.
Restart qBittorrent after changing this setting.
`

const portalHelp = `Usage: caelestia-extras portal sync

Copies generated GTK and Qt assets into the isolated portal configuration.
Requires a [portal] section in the configuration.
`

const configHelp = `Usage: caelestia-extras config validate

Checks the configured files, required commands, and enabled integrations.
Prints an explicit confirmation and the checked path when validation succeeds.
`

const completionHelp = `Usage: caelestia-extras completion <shell>

Prints completion scripts for bash, zsh, or fish.

Examples:
  source <(caelestia-extras completion bash)
  caelestia-extras completion zsh > ~/.zfunc/_caelestia-extras
  caelestia-extras completion fish | source
`

const bashCompletion = `# bash completion for caelestia-extras
_caelestia_extras() {
  local cur prev command
  cur="${COMP_WORDS[COMP_CWORD]}"
  prev="${COMP_WORDS[COMP_CWORD-1]}"
  command="${COMP_WORDS[1]}"
  if [[ $COMP_CWORD == 1 ]]; then
    COMPREPLY=($(compgen -W "sync watch cursor gtk hyprtoolkit pavucontrol qt qbittorrent portal config completion help version" -- "$cur"))
  elif [[ $COMP_CWORD == 2 && $command == cursor ]]; then
    COMPREPLY=($(compgen -W "sync sync-xcursor" -- "$cur"))
  elif [[ $COMP_CWORD == 2 && $command == completion ]]; then
    COMPREPLY=($(compgen -W "bash zsh fish" -- "$cur"))
  elif [[ $COMP_CWORD == 2 && $command == help ]]; then
    COMPREPLY=($(compgen -W "sync watch cursor gtk hyprtoolkit pavucontrol qt qbittorrent portal config completion" -- "$cur"))
  fi
}
complete -F _caelestia_extras caelestia-extras
`

const zshCompletion = `#compdef caelestia-extras

_arguments \
  '(-config 1)'{-config,-config}'[configuration file]:file:_files' \
  '(-version 1)'--version'[show version]' \
  '1:command:(sync watch cursor gtk hyprtoolkit pavucontrol qt qbittorrent portal config completion help version)' \
  '*::argument:->args'

case $words[2] in
  cursor) _arguments '1:action:(sync sync-xcursor)' ;;
  qt|qbittorrent) _arguments '1:action:(sync)' ;;
  completion) _arguments '1:shell:(bash zsh fish)' ;;
  help) _arguments '1:command:(sync watch cursor gtk hyprtoolkit pavucontrol qt qbittorrent portal config completion)' ;;
  config) _arguments '1:action:(validate)' ;;
esac
`

const fishCompletion = `complete -c caelestia-extras -f -n '__fish_use_subcommand' -a 'sync watch cursor gtk hyprtoolkit pavucontrol qt qbittorrent portal config completion help version'
complete -c caelestia-extras -l config -r -d 'configuration file'
complete -c caelestia-extras -l version -d 'show version'
complete -c caelestia-extras -n '__fish_seen_subcommand_from cursor' -a 'sync sync-xcursor'
complete -c caelestia-extras -n '__fish_seen_subcommand_from qt' -a 'sync'
complete -c caelestia-extras -n '__fish_seen_subcommand_from qbittorrent' -a 'sync'
complete -c caelestia-extras -n '__fish_seen_subcommand_from completion' -a 'bash zsh fish'
complete -c caelestia-extras -n '__fish_seen_subcommand_from help' -a 'sync watch cursor gtk hyprtoolkit pavucontrol qbittorrent portal config completion'
complete -c caelestia-extras -n '__fish_seen_subcommand_from config' -a 'validate'
`
