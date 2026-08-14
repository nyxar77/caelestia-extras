# Manual install

The script is for a normal source checkout. It does not install packages, edit
an existing config, or enable anything without `--enable`.

## Requirements

Install Go 1.25 or newer and Caelestia first. Add the tools for the parts you
use:

- Cursor: `hyprcursor-util` and `hyprctl`
- GTK: `dconf`
- XCursor fallback: `cbmp` and `ctgen`
- pavucontrol: `pavucontrol-qt`

## Install or update

From the checkout:

```sh
./scripts/install.sh
```

It builds the current checkout with `-trimpath` and `-buildvcs=false`, installs
the binary in `$XDG_BIN_HOME` or `~/.local/bin`, and creates a starter config
only when one does not already exist.

Managed templates, portal files, and systemd units are stored under
`$XDG_CONFIG_HOME/caelestia-extras/managed/`. Active files are symlinks to that
directory. Rerunning the script updates managed files. If you replace a managed
symlink with your own file, the script leaves your file alone.

The build and file rendering happen in a temporary directory first. Press
`Ctrl-C` before the brief apply step to cancel without changing installed files.
Once that step starts, interruption is ignored so the managed set is not left
half-applied.

Run the same command after updating the checkout:

```sh
git pull
./scripts/install.sh update
```

## Enable integrations

Edit the config first, then check it:

```sh
~/.local/bin/caelestia-extras config validate
```

Enable only the integrations you configured:

```sh
./scripts/install.sh --enable cursor,gtk
```

Valid names are `cursor`, `gtk`, `hyprtoolkit`, and `portal`. Use `all` for all
four. The script enables the matching path units and runs an initial sync.

`pavucontrol` has no background service; run `caelestia-extras pavucontrol`
when you want it.

Portal service drop-ins take effect after the next login. The script does not
restart portal processes for you.

## XDG paths

The script respects `XDG_CONFIG_HOME`, `XDG_DATA_HOME`, `XDG_STATE_HOME`, and
`XDG_BIN_HOME`. The generated units watch the matching default Caelestia paths.

If your config overrides `scheme.file`, `hyprtoolkit.theme_dir`, or
`portal.theme_dir`, pass the same paths while installing and updating:

```sh
CAELESTIA_EXTRAS_SCHEME_FILE=/path/to/scheme.json \
CAELESTIA_EXTRAS_THEME_DIR=/path/to/theme \
./scripts/install.sh update
```

The manual portal files use `Caelestia-Portal` by default. If
`portal.theme_name` differs, set it on every install or update too:

```sh
CAELESTIA_EXTRAS_PORTAL_THEME_NAME=My-Portal-Theme ./scripts/install.sh update
```
