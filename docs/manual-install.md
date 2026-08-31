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
- Shared Qt theme: `qt5ct`, `qt6ct`, and Breeze for Qt 6
- qBittorrent: `qbittorrent`
- Portal sync and the watcher: a systemd user session with `systemctl`

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

The GTK 3 and GTK 4 user stylesheets are symlinked directly to Caelestia's
generated `gtk.css` and `gtk4.css`, respectively. Existing user stylesheets are
preserved instead of being overwritten.

Caelestia CLI can also write those two user stylesheets. If you enable the GTK
integration here, set `theme.enableGtk` to `false` in Caelestia's CLI settings
so it keeps rendering Extras' templates without replacing the installed
symlinks.

The installer also exposes the generated `Caelestia-GTK3` named theme for GTK
3 applications whose stock theme does not consume compatibility colour names.
It imports Adwaita and changes colour properties only.

The build and file rendering happen in a temporary directory first. Press
`Ctrl-C` before the brief apply step to cancel without changing installed
files. Interruption is ignored only while the staged files are copied into
place, then restored before service management and the initial sync.

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

Enable the unified watcher after configuring the integrations you want:

```sh
./scripts/install.sh --enable all
```

The script enables one watcher and runs `caelestia-extras sync` once. Disabled
config sections are skipped.

`pavucontrol` has no background service; run `caelestia-extras pavucontrol`
when you want it. qBittorrent's native palette is owned by the shared Qt sync.
Restart qBittorrent once after changing its appearance setting; later palette
updates are propagated by qt6ct.

The initial sync restarts the portal backends so their scoped GTK and Qt themes
take effect immediately.

The `qt` integration also writes an `environment.d` file and a generated
Caelestia KDE colour scheme. Log out and back in before launching Qt
applications so they receive the Qt5/Qt6 platform theme.

## XDG paths

The script respects `XDG_CONFIG_HOME`, `XDG_DATA_HOME`, `XDG_STATE_HOME`, and
`XDG_BIN_HOME`. The watcher reads source paths from `config.toml`.

`CAELESTIA_EXTRAS_THEME_DIR` is needed only when the manual PrismLauncher
symlink uses a non-default generated-theme directory:

```sh
CAELESTIA_EXTRAS_THEME_DIR=/path/to/theme \
./scripts/install.sh update
```

The manual portal files use `Caelestia-Portal` by default. If
`portal.theme_name` differs, set it on every install or update too:

```sh
CAELESTIA_EXTRAS_PORTAL_THEME_NAME=My-Portal-Theme ./scripts/install.sh update
```
