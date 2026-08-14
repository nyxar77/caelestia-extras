# caelestia-extras

Small Linux integrations for [Caelestia](https://github.com/caelestia-dots/shell). It keeps applications and desktop components in sync with Caelestia's active colour scheme without replacing Caelestia or managing the rest of the desktop.

## Integrations

- **Cursor** — builds a dynamic Bibata Hyprcursor theme and can maintain an XCursor fallback.
- **GTK** — updates the GNOME colour preference and GTK theme when the scheme changes.
- **Hyprtoolkit** — copies Caelestia's generated Hyprtoolkit configuration into its active config path.
- **pavucontrol-qt** — launches the volume mixer with Caelestia's generated stylesheet.
- **XDG portals** — gives GTK and Hyprland portals isolated Caelestia themes and Qt settings.

The live cursor refresh uses Hyprland tools. The other integrations are regular Linux user-session services.

## Requirements

For a Nix installation, the flake provides the package and runtime dependencies used by the Home Manager module.

For a manual installation, install:

- Go 1.25 or newer to build the binary
- Caelestia and its generated scheme/theme files
- `hyprcursor-util`, `hyprctl`, and `dconf` for the relevant integrations
- `cbmp` and `ctgen` for the XCursor fallback

## Home Manager

Add the flake and module to your Home Manager configuration:

```nix
inputs.caelestia-extras.url = "github:nyxar77/caelestia-extras";

imports = [ inputs.caelestia-extras.homeModules.default ];

programs.caelestia-extras = {
  enable = true;
  cursor.enable = true;
  gtk.enable = true;
  hyprtoolkit.enable = true;
  pavucontrol.enable = true;
  portal.enable = true;
};
```

Each integration is optional. The module writes the configuration, installs the enabled user services and path watchers, and runs the initial sync for Hyprtoolkit and portal output.

Useful options include:

- `schemeFile` — Caelestia scheme JSON; defaults to `$XDG_STATE_HOME/caelestia/scheme.json`.
- `cursor.source` and `cursor.buildConfig` — Bibata SVG and build-config paths. Nix defaults to the packaged Bibata source.
- `cursor.xcursorFallback` and `cursor.updateGtk` — control the fallback and GTK cursor updates.
- `gtk.darkTheme` and `gtk.lightTheme` — themes selected for each mode.
- `hyprtoolkit.themeDir` and `hyprtoolkit.configFile` — generated and active config paths.
- `portal.themeName`, `portal.iconTheme`, and `portal.applyGlobalGtk` — portal theme settings.

GTK applications that keep a stale D-Bus-activated process can be launched directly with a desktop-entry override:

```nix
programs.caelestia-extras.gtk.directLaunch."org.gnome.Nautilus" = {
  name = "Files";
  exec = "nautilus --new-window %U";
  icon = "org.gnome.Nautilus";
};
```

The attribute is the desktop ID without `.desktop`. The generated entry sets `DBusActivatable=false`.

## Manual configuration

Build the binary from a checkout:

```sh
go build -o caelestia-extras ./cmd/caelestia-extras
```

Create `~/.config/caelestia-extras/config.toml`:

```toml
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

[portal]
theme_name = "Caelestia-Portal"
apply_global_gtk = true
```

The `[scheme]` section is optional when Caelestia uses its default path. Other sections are enabled by adding them to the file. Run an integration with:

```sh
caelestia-extras cursor sync
caelestia-extras cursor sync-xcursor
caelestia-extras gtk sync
caelestia-extras hyprtoolkit sync
caelestia-extras pavucontrol
caelestia-extras portal sync
```

Use `--config PATH` to select another configuration file and `--version` to print the binary version.

Run `caelestia-extras --help` for the full command list. Every integration has
command-specific help, for example:

```sh
caelestia-extras help cursor
caelestia-extras portal --help
```

Generate completions for the shell you use:

```sh
# Bash
source <(caelestia-extras completion bash)

# Zsh
caelestia-extras completion zsh > ~/.zfunc/_caelestia-extras

# Fish
caelestia-extras completion fish | source
```

The `systemd/` directory contains user service and path-unit files for the cursor, GTK, Hyprtoolkit, and portal integrations. Copy the units you enable to `~/.config/systemd/user/`, then enable the corresponding `.path` units. The binary must be available in the user service `PATH`.

## License

[GPL-3.0-only](LICENSE)
