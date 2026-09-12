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
defaults below. Unknown sections and keys are rejected so a misspelling cannot
silently disable part of the configuration.

## Options

### `[compositor]`

- `backend` — compositor-specific backend. Default: `hyprland`.

Only the cursor integration uses the compositor backend. Keep it explicit in
service configurations; a systemd user service may not receive enough session
state for reliable compositor detection.

### `[scheme]`

- `file` — Caelestia's active scheme JSON. Default:
  `$XDG_STATE_HOME/caelestia/scheme.json`.

Cursor and GTK sync require this file. The remaining integrations consume
Caelestia's rendered template files from their `theme_dir`.

### `[cursor]`

- `source` — Bibata SVG directory. Required.
- `build_config` — Bibata build configuration. Required.
- `icon_dir` — install directory. Default: `$XDG_DATA_HOME/icons`.
- `theme` — generated theme name. Default: `Bibata-Caelestia`.
- `size` — cursor size. Default: `20`.
- `xcursor_sizes` — XCursor sizes. Default: `[20, 24, 32]`.
- `xcursor_fallback` — refresh the XCursor fallback. Default: `false`.
- `update_gtk` — update the GTK cursor setting. Default: `false`.

### `[gtk]`

- `dark_theme` — Default: `adw-gtk3-dark`.
- `light_theme` — Default: `adw-gtk3`.

Home Manager and the manual installer link the generated GTK 3 `gtk.css` and
GTK 4 `gtk4.css` into their respective user-style locations. Keeping them
separate prevents GTK 3 from parsing GTK 4-only CSS variables. Both stylesheets
change only colour tokens; the toolkit still owns widget layout and interaction
states.

They also install the generated `Caelestia-GTK3` named theme. It uses stock
Adwaita as its base and is intended for GTK 3 applications that do not consume
the global compatibility colour names.

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

### LocalSend (Home Manager only)

- `programs.caelestia-extras.localsend.enable` — install LocalSend with its
  Linux GTK host window scoped to the generated `Caelestia-GTK3` theme.
- `programs.caelestia-extras.localsend.package` — LocalSend package to wrap.
  Default: `pkgs.localsend`.

LocalSend's Flutter content keeps using its own system colour mode. The wrapper
only themes the native Linux title bar and requires the GTK integration.

### imv (Home Manager and manual installer)

- `programs.caelestia-extras.imv.enable` — install the generated imv theme.
  Defaults to `programs.caelestia-extras.autoEnable`.
- `programs.caelestia-extras.imv.themeDir` — generated theme directory.
  Default: `$XDG_STATE_HOME/caelestia/theme`.

The generated configuration sets imv's canvas, overlay text, and overlay
background colours. imv has no config include mechanism, so the integration
owns `$XDG_CONFIG_HOME/imv/config` rather than adding a fragment to an existing
file. The manual installer preserves an existing user-owned file.

### MPV with ModernZ (Home Manager only)

- `programs.caelestia-extras.mpv.enable` — install the generated ModernZ theme.
  Default: `false`.
- `programs.caelestia-extras.mpv.themeDir` — generated theme directory.
  Default: `$XDG_STATE_HOME/caelestia/theme`.

This integration links the rendered `modernz.conf` into MPV's `script-opts`
directory. It requires `programs.mpv.enable` and verifies that
`programs.mpv.finalPackage` actually loads `modernz.lua`. Add
`pkgs.mpvScripts.modernz` through `programs.mpv.scripts` or through the `scripts`
argument of a custom MPV package override.

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

The preferences file must be writable and must not be a symlink. Sync preserves
the permissions of an existing regular file and refuses to replace a symlink.

### `[portal]`

- `theme_dir` — generated portal theme directory. Default:
  `$XDG_STATE_HOME/caelestia/theme`.
- `config_home` — XDG config directory. Default: `$XDG_CONFIG_HOME`.
- `data_home` — XDG data directory. Default: `$XDG_DATA_HOME`.
- `theme_name` — private GTK theme name. Default: `Caelestia-Portal`.

An aggregate `sync` restarts the GTK, GNOME, and Hyprland portal user services after
copying the generated files. `config validate` therefore requires `systemctl`
when this section is enabled. The direct `portal sync` command only copies the
files.

## Validation and editor support

Run this before debugging an integration:

```sh
caelestia-extras config validate
```

It checks the files and commands needed by the enabled integrations. It does
not require generated theme output to exist yet. Validation checks the runtime
TOML, not only the editor schema.

The repository includes [`../config/caelestia-extras.schema.json`](../config/caelestia-extras.schema.json).
Associate it with this TOML file in your TOML language server for key
completion, type checking, and inline documentation.
