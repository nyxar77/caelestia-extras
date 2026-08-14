# caelestia-extras

Small Linux integrations for [Caelestia](https://github.com/caelestia-dots/shell).
It keeps the cursor, GTK, Hyprtoolkit, pavucontrol, qBittorrent, and XDG portals in step
with Caelestia's active scheme.

The current compositor backend is Hyprland. The shared integrations are kept
separate so another compositor can be added without rewriting them.

## Installation

### Home Manager

Add the flake and module:

```nix
inputs.caelestia-extras.url = "github:nyxar77/caelestia-extras";

imports = [ inputs.caelestia-extras.homeModules.default ];

programs.caelestia-extras = {
  enable = true;
  cursor.enable = true;
  gtk.enable = true;
  hyprtoolkit.enable = true;
  pavucontrol.enable = true;
  prismlauncher.enable = true;
  qbittorrent.enable = true;
  portal.enable = true;
};
```

All integrations are optional. The module installs the package, writes the
configuration, and creates the user services and path watchers.

### Manual install

You need Go 1.25 or newer, Caelestia, and the tools used by the integrations
you enable. The script does not install distro packages or overwrite your
config.

```sh
./scripts/install.sh
```

It creates `~/.config/caelestia-extras/config.toml` if it does not exist,
installs the managed templates and units, and places the binary in
`~/.local/bin`. Configure an integration, then enable the ones you want:

```sh
./scripts/install.sh --enable cursor,gtk,hyprtoolkit,portal
```

It stages files before applying them, so `Ctrl-C` can safely cancel the build.

Run the same command after pulling an update. More detail is in the
[manual-install guide](docs/manual-install.md).

## Quick start

```sh
caelestia-extras cursor sync
caelestia-extras gtk sync
caelestia-extras portal sync
```

Check the setup before running an integration:

```sh
caelestia-extras config validate
```

More commands, help, and shell completion are in the [CLI guide](docs/cli.md).

## Docs

- [Configuration](docs/configuration.md)
- [Integrations and systemd](docs/integrations.md)
- [CLI and shell completion](docs/cli.md)
- [Manual install and updates](docs/manual-install.md)
- [Development](docs/development.md)

## License

[GPL-3.0-only](LICENSE)
