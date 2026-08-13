# caelestia-extras

Optional live integrations for [Caelestia](https://github.com/caelestia-dots/shell). It fills gaps around applications and cursor handling; it does not replace Caelestia or manage a whole desktop.

## Features

- Dynamic Bibata cursor from Caelestia's active scheme, with direct Hyprcursor output and an optional XCursor fallback.
- `pavucontrol-qt` launcher that uses Caelestia's generated stylesheet.
- GTK light/dark preference sync.
- Hyprtoolkit configuration generated from Caelestia's active scheme.
- XDG portal theming: a private GTK theme for `xdg-desktop-portal-gtk`, isolated Qt styling for `xdg-desktop-portal-hyprland`, and generated GTK/qt6ct assets.

The cursor's instant apply step is Hyprland-specific. The rest is normal Linux user-session tooling.

## Install

The core is a Go binary. Build it from this checkout:

```sh
go build ./cmd/caelestia-extras
```

The Nix flake exposes a package and a Home Manager module:

```nix
inputs.caelestia-extras.url = "path:/path/to/caelestia-extras";

imports = [inputs.caelestia-extras.homeModules.default];

programs.caelestia-extras = {
  enable = true;
  cursor.enable = true;
  gtk.enable = true;
  hyprtoolkit.enable = true;
  pavucontrol.enable = true;
  portal.enable = true;
};
```

The portal module adds Caelestia's template files, portal-specific service
drop-ins, and a path unit that syncs generated output after every theme change.
Its defaults reproduce the setup used by this project:

```nix
programs.caelestia-extras.portal = {
  themeName = "Caelestia-Portal";
  iconTheme = "Papirus-Dark";
  applyGlobalGtk = true;
};
```

`themeDir`, `configHome`, and `dataHome` are also configurable when Caelestia
uses non-default XDG locations.

Some GTK applications are D-Bus activated from an app launcher and retain an
older process context. `gtk.directLaunch` can override those desktop entries
without baking any application into Extras:

```nix
programs.caelestia-extras.gtk.directLaunch."org.gnome.Nautilus" = {
  name = "Files";
  exec = "nautilus --new-window %U";
  icon = "org.gnome.Nautilus";
};
```

The attribute name is the desktop ID without `.desktop`; Extras writes the
override with `DBusActivatable=false`.

## Non-Nix configuration

Create `~/.config/caelestia-extras/config.toml` and point the cursor section at a Bibata source checkout:

```toml
[cursor]
source = "/path/to/Bibata_Cursor/svg/modern"
build_config = "/path/to/Bibata_Cursor/configs/normal/x.build.toml"
xcursor_fallback = true

[gtk]

[hyprtoolkit]

[pavucontrol]

[portal]
theme_name = "Caelestia-Portal"
apply_global_gtk = true
```

`cursor sync` needs `hyprcursor-util`, `hyprctl`, and, when GTK cursor updates are enabled, `dconf`. `cursor sync-xcursor` additionally needs `cbmp` and `ctgen`.

```sh
caelestia-extras cursor sync
caelestia-extras gtk sync
caelestia-extras hyprtoolkit sync
caelestia-extras portal sync
caelestia-extras pavucontrol
```

For non-Nix setups, create the static portal templates and service drop-ins in
your own configuration, then copy the units in `systemd/` to
`~/.config/systemd/user/`. Enable the cursor and portal path units if you are
not using the Home Manager module.

## Development

```sh
go test ./...
nix flake check
```
