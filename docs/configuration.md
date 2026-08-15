# Configuration

The default file is:

```text
$XDG_CONFIG_HOME/caelestia-extras/config.toml
```

If `XDG_CONFIG_HOME` is not set, this is `~/.config/caelestia-extras/config.toml`.

## Example

```toml
[compositor]
backend = "hyprland"

[scheme]
file = "/path/to/caelestia/scheme.json"

[cursor]
source = "/path/to/Bibata_Cursor/svg/modern"
build_config = "/path/to/Bibata_Cursor/configs/normal/x.build.toml"
xcursor_fallback = true
update_gtk = true

[gtk]

[hyprtoolkit]

[pavucontrol]

[qt]

[qbittorrent]

[portal]
theme_name = "Caelestia-Portal"
```

The `scheme.file` setting defaults to `$XDG_STATE_HOME/caelestia/scheme.json`.
An integration is enabled by adding its section. Empty sections use the
defaults below.

## Options

### `[cursor]`

- `source` — Bibata SVG directory. Required.
- `build_config` — Bibata build configuration. Required.
- `icon_dir` — install directory. Default: `$XDG_DATA_HOME/icons`.
- `theme` — generated theme name. Default: `Bibata-Caelestia`.
- `size` — cursor size. Default: `20`.
- `xcursor_sizes` — XCursor sizes. Default: `[20, 24, 32]`.
- `xcursor_fallback` — refresh the XCursor fallback. Default: `false`.
- `update_gtk` — update the GTK cursor setting. Default: `false`.

### `[compositor]`

- `backend` — compositor-specific backend. Default: `hyprland`.

Keep this explicit. User services do not always receive enough session state
to detect the compositor reliably.

### `[gtk]`

- `dark_theme` — Default: `adw-gtk3-dark`.
- `light_theme` — Default: `adw-gtk3`.

Home Manager and the manual installer link the generated `gtk.css` into the
GTK 3 and GTK 4 user-style locations. The stylesheet changes only GTK and
libadwaita colour tokens; the toolkit still owns widget layout and interaction
states.

### `[hyprtoolkit]`

- `theme_dir` — generated theme directory. Default:
  `$XDG_STATE_HOME/caelestia/theme`.
- `config_file` — active Hyprtoolkit config. Default:
  `$XDG_CONFIG_HOME/hypr/hyprtoolkit.conf`.

### `[pavucontrol]`

- `command` — mixer command. Default: `pavucontrol-qt`.

Enabling this integration does not install the application. Validation and the
launcher report a missing command as an error instead of silently doing
nothing.

### `[qt]`

- `theme_dir` — generated palette and Breeze colour-scheme directory. Default:
  `$XDG_STATE_HOME/caelestia/theme`.
- `config_home` — XDG config directory. Default: `$XDG_CONFIG_HOME`.
- `data_home` — XDG data directory. Default: `$XDG_DATA_HOME`.

The integration copies Caelestia's palette to qt5ct and qt6ct, installs a
generated KDE colour scheme, and selects it in `kdeglobals`. Qt 6 uses Breeze;
Qt 5 keeps Fusion because this Nixpkgs revision does not ship a Qt 5 Breeze
plugin. Home Manager configures the platform theme and plugin path.

### `[qbittorrent]`

- `command` — qBittorrent command. Default: `qbittorrent`.
- `config_file` — qBittorrent config file. Default:
  `$XDG_CONFIG_HOME/qBittorrent/qBittorrent.conf`.

The sync disables qBittorrent's custom UI theme and selects its native system
style. Home Manager also supplies a desktop-entry wrapper with a scoped qt6ct
and Breeze plugin environment. The wrapper applies the sync before every start,
so qBittorrent cannot restore a stale custom-theme setting on shutdown. It also
loads a top-bar-only stylesheet that references Qt palette roles; the shared
Breeze palette still owns all actual colours.

### `[portal]`

- `theme_dir` — generated portal theme directory. Default:
  `$XDG_STATE_HOME/caelestia/theme`.
- `config_home` — XDG config directory. Default: `$XDG_CONFIG_HOME`.
- `data_home` — XDG data directory. Default: `$XDG_DATA_HOME`.
- `theme_name` — private GTK theme name. Default: `Caelestia-Portal`.

## Validation and editor support

Run this before debugging an integration:

```sh
caelestia-extras config validate
```

It checks the files and commands needed by the enabled integrations. It does
not require generated theme output to exist yet.

The repository includes [`../config/caelestia-extras.schema.json`](../config/caelestia-extras.schema.json).
Associate it with this TOML file in your TOML language server for key
completion, type checking, and inline documentation.
