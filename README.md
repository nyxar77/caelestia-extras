# caelestia-extras

Optional live integrations for [Caelestia](https://github.com/caelestia-dots/shell). It fills gaps around applications and cursor handling; it does not replace Caelestia or manage a whole desktop.

## Features

- Dynamic Bibata cursor from Caelestia's active scheme, with direct Hyprcursor output and an optional XCursor fallback.
- `pavucontrol-qt` launcher that uses Caelestia's generated stylesheet.
- GTK light/dark preference sync.

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
  pavucontrol.enable = true;
};
```

## Non-Nix configuration

Create `~/.config/caelestia-extras/config.toml` and point the cursor section at a Bibata source checkout:

```toml
[cursor]
source = "/path/to/Bibata_Cursor/svg/modern"
build_config = "/path/to/Bibata_Cursor/configs/normal/x.build.toml"
xcursor_fallback = true

[gtk]

[pavucontrol]
```

`cursor sync` needs `hyprcursor-util`, `hyprctl`, and, when GTK cursor updates are enabled, `dconf`. `cursor sync-xcursor` additionally needs `cbmp` and `ctgen`.

```sh
caelestia-extras cursor sync
caelestia-extras gtk sync
caelestia-extras pavucontrol
```

Copy the units in `systemd/` to `~/.config/systemd/user/` and enable the cursor path unit if you are not using the Home Manager module.

## Development

```sh
go test ./...
nix flake check
```
