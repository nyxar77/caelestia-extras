# CLI

Global options must come before the command.

```text
caelestia-extras [options] <command> [arguments]
```

## Commands

```text
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
version                  Print the version
```

Commands report the action they start and confirm completion. Diagnostics use
red for errors, yellow for watcher warnings, cyan for progress, and green for
success when output is a terminal. Redirected output and `NO_COLOR` remain
plain text. A successful validation ends with `Configuration is valid` and the
path that was checked.

Use `--help` or `help <command>` for details:

```sh
caelestia-extras --help
caelestia-extras help sync
caelestia-extras help cursor
caelestia-extras portal --help
```

Use another config file with `--config`:

```sh
caelestia-extras --config ~/work/caelestia-extras.toml gtk sync
```

## Shell completion

The program prints a completion script. The shell loads it and handles Tab.

```sh
# Bash
source <(caelestia-extras completion bash)

# Zsh
caelestia-extras completion zsh > ~/.zfunc/_caelestia-extras

# Fish
caelestia-extras completion fish | source
```
