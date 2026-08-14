package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

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
		fmt.Fprintln(os.Stderr, "error:", err)
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
	case "cursor":
		if configFile.Cursor == nil {
			return fmt.Errorf("cursor is disabled in %q; add a [cursor] section", *configPath)
		}
		if len(args) != 2 {
			return usage("cursor expects one of: sync, sync-xcursor")
		}
		switch args[1] {
		case "sync":
			return runIntegration("cursor sync", func() error {
				return cursor.Sync(*configFile.Cursor, configFile.Scheme.File)
			})
		case "sync-xcursor":
			return runIntegration("cursor XCursor sync", func() error {
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
		return runIntegration("GTK sync", func() error {
			return integration.SyncGTK(*configFile.GTK, configFile.Scheme.File)
		})
	case "hyprtoolkit":
		if configFile.Hyprtoolkit == nil {
			return fmt.Errorf("Hyprtoolkit sync is disabled in %q; add a [hyprtoolkit] section", *configPath)
		}
		if len(args) != 2 || args[1] != "sync" {
			return usage("hyprtoolkit expects: sync")
		}
		return runIntegration("Hyprtoolkit sync", func() error {
			return integration.SyncHyprtoolkit(*configFile.Hyprtoolkit)
		})
	case "pavucontrol":
		if configFile.Pavucontrol == nil {
			return fmt.Errorf("pavucontrol is disabled in %q; add a [pavucontrol] section", *configPath)
		}
		return runIntegration("launch pavucontrol", func() error {
			return integration.LaunchPavucontrol(*configFile.Pavucontrol, args[1:])
		})
	case "portal":
		if configFile.Portal == nil {
			return fmt.Errorf("portal sync is disabled in %q; add a [portal] section", *configPath)
		}
		if len(args) != 2 || args[1] != "sync" {
			return usage("portal expects: sync")
		}
		return runIntegration("portal sync", func() error {
			return integration.SyncPortal(*configFile.Portal)
		})
	default:
		return usage(fmt.Sprintf("unknown command %q; expected cursor, gtk, hyprtoolkit, pavucontrol, portal, completion, help, or version", args[0]))
	}
}

func validateCommand(args []string) error {
	switch args[0] {
	case "cursor":
		if len(args) != 2 {
			return usage("cursor expects one of: sync, sync-xcursor")
		}
		if args[1] != "sync" && args[1] != "sync-xcursor" {
			return usage(fmt.Sprintf("unknown cursor action %q; expected sync or sync-xcursor", args[1]))
		}
	case "gtk", "hyprtoolkit", "portal":
		if len(args) != 2 || args[1] != "sync" {
			return usage(fmt.Sprintf("%s expects: sync", args[0]))
		}
	case "pavucontrol":
		return nil
	default:
		return usage(fmt.Sprintf("unknown command %q; expected cursor, gtk, hyprtoolkit, pavucontrol, portal, completion, help, or version", args[0]))
	}
	return nil
}

func runIntegration(name string, run func() error) error {
	if err := run(); err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	return nil
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
	case "cursor":
		help = cursorHelp
	case "gtk":
		help = gtkHelp
	case "hyprtoolkit":
		help = hyprtoolkitHelp
	case "pavucontrol":
		help = pavucontrolHelp
	case "portal":
		help = portalHelp
	case "completion":
		help = completionHelp
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
  cursor sync              Build and apply the active Hyprcursor theme
  cursor sync-xcursor      Build the XCursor fallback
  gtk sync                 Sync the GTK theme and colour preference
  hyprtoolkit sync         Apply generated Hyprtoolkit configuration
  pavucontrol [arguments]  Launch the configured volume mixer
  portal sync              Apply generated portal themes
  completion <shell>       Print shell completion for bash, zsh, or fish
  help [command]           Show general or command-specific help
  version                  Print the version

Options:
  -config PATH             Use PATH instead of the default config file
  -version                 Print the version
  -h, --help               Show this help

Examples:
  caelestia-extras cursor sync
  caelestia-extras --config ~/.config/caelestia-extras/work.toml gtk sync
  caelestia-extras completion zsh > ~/.zfunc/_caelestia-extras

Configuration:
  Default: $XDG_CONFIG_HOME/caelestia-extras/config.toml
  Help:    https://github.com/nyxar77/caelestia-extras#manual-configuration
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

const portalHelp = `Usage: caelestia-extras portal sync

Copies generated GTK and Qt assets into the isolated portal configuration.
Requires a [portal] section in the configuration.
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
    COMPREPLY=($(compgen -W "cursor gtk hyprtoolkit pavucontrol portal completion help version" -- "$cur"))
  elif [[ $COMP_CWORD == 2 && $command == cursor ]]; then
    COMPREPLY=($(compgen -W "sync sync-xcursor" -- "$cur"))
  elif [[ $COMP_CWORD == 2 && $command == completion ]]; then
    COMPREPLY=($(compgen -W "bash zsh fish" -- "$cur"))
  elif [[ $COMP_CWORD == 2 && $command == help ]]; then
    COMPREPLY=($(compgen -W "cursor gtk hyprtoolkit pavucontrol portal completion" -- "$cur"))
  fi
}
complete -F _caelestia_extras caelestia-extras
`

const zshCompletion = `#compdef caelestia-extras

_arguments \
  '(-config 1)'{-config,-config}'[configuration file]:file:_files' \
  '(-version 1)'--version'[show version]' \
  '1:command:(cursor gtk hyprtoolkit pavucontrol portal completion help version)' \
  '*::argument:->args'

case $words[2] in
  cursor) _arguments '1:action:(sync sync-xcursor)' ;;
  completion) _arguments '1:shell:(bash zsh fish)' ;;
  help) _arguments '1:command:(cursor gtk hyprtoolkit pavucontrol portal completion)' ;;
esac
`

const fishCompletion = `complete -c caelestia-extras -f -n '__fish_use_subcommand' -a 'cursor gtk hyprtoolkit pavucontrol portal completion help version'
complete -c caelestia-extras -l config -r -d 'configuration file'
complete -c caelestia-extras -l version -d 'show version'
complete -c caelestia-extras -n '__fish_seen_subcommand_from cursor' -a 'sync sync-xcursor'
complete -c caelestia-extras -n '__fish_seen_subcommand_from completion' -a 'bash zsh fish'
complete -c caelestia-extras -n '__fish_seen_subcommand_from help' -a 'cursor gtk hyprtoolkit pavucontrol portal completion'
`
