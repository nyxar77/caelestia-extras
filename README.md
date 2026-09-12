# caelestia-extras

`caelestia-extras` applies [Caelestia](https://github.com/caelestia-dots/shell)
colours to cursors, toolkits, portals, and applications. It uses Caelestia's
scheme and template output, then keeps the configured consumers synchronized
after a wallpaper or scheme change.

The supported compositor backend is Hyprland. The program is Linux-only and
the supplied service expects a systemd user session.

## Integrations

| Integration | What it changes |
| --- | --- |
| Cursor | Builds and applies a recoloured Bibata Hyprcursor theme, with an optional XCursor fallback |
| GTK | Updates the light/dark preference and GTK theme; installs Caelestia colour stylesheets |
| Hyprtoolkit | Copies Caelestia's generated Hyprtoolkit configuration into place |
| [imv](https://sr.ht/~exec64/imv/) | Colours the image canvas and information overlay from Caelestia's generated palette |
| MPV | Colours the ModernZ on-screen controller from Caelestia's generated palette |
| Qt | Updates qt5ct, qt6ct, KDE/Breeze colours, and the Qt session environment |
| pavucontrol | Provides a `pavucontrol-qt` launcher using the generated stylesheet |
| qBittorrent | Restores native Qt/Breeze theming and provides a scoped launcher under Home Manager |
| PrismLauncher | Installs Caelestia's generated PrismLauncher theme under Home Manager or the manual installer |
| LocalSend | Wraps LocalSend's GTK host window with the generated GTK 3 theme under Home Manager |
| XDG portals | Gives GTK/GNOME choosers and the Hyprland screencast picker isolated GTK/Qt themes |

See [Integrations](docs/integrations.md) for the exact files and settings each
integration owns.

## Home Manager

Add the flake input and import its module:

```nix
{
  inputs.caelestia-extras = {
    url = "github:nyxar77/caelestia-extras";
    inputs.nixpkgs.follows = "nixpkgs";
    inputs.home-manager.follows = "home-manager";
  };
}
```

```nix
{inputs, ...}: {
  imports = [inputs.caelestia-extras.homeManagerModules.default];

  programs.caelestia-extras = {
    enable = true;
    autoEnable = true;

    # Optional per-integration opt-out:
    # localsend.enable = false;
  };
}
```

`autoEnable` defaults to `true`. Every integration inherits it, and an explicit
`integration.enable` value wins, so set unwanted integrations to `false`. Set
`autoEnable = false` to return to an opt-in list instead. The MPV integration is
opt-in because it requires ModernZ. LocalSend requires the GTK integration;
imv, MPV, pavucontrol, qBittorrent, and PrismLauncher are not installed by their
integration flags.

If Caelestia's own Home Manager module manages the CLI, disable its built-in
GTK writer before enabling `gtk` here:

```nix
programs.caelestia.cli.settings.theme.enableGtk = false;
```

Both integrations otherwise write `gtk-3.0/gtk.css` and `gtk-4.0/gtk.css`.

The module installs the CLI, writes its TOML configuration, runs an initial
sync during activation, and starts a watcher when an enabled integration has
generated files to monitor. Read the [Home Manager guide](docs/home-manager.md)
before enabling it in an existing desktop configuration: GTK, Qt, portal, and
desktop-entry files need a single declarative owner.

## Manual installation

Manual installation requires Go 1.25 or newer. The installer builds the
current checkout and preserves an existing configuration:

```sh
./scripts/install.sh
```

Edit `~/.config/caelestia-extras/config.toml`, validate it, then enable the
watcher and run the first sync:

```sh
~/.local/bin/caelestia-extras config validate
./scripts/install.sh update --enable all
```

The script does not install distribution packages. Required commands and the
files managed by the script are listed in the
[manual installation guide](docs/manual-install.md).

## Manual configuration

The default configuration is
`$XDG_CONFIG_HOME/caelestia-extras/config.toml`, or
`~/.config/caelestia-extras/config.toml` when `XDG_CONFIG_HOME` is unset. A TOML
section enables the corresponding runtime integration.

Run this after editing the file:

```sh
caelestia-extras config validate
caelestia-extras sync
```

See [Configuration](docs/configuration.md) for a complete example, defaults,
and the included editor schema.

## Documentation

- [Home Manager module](docs/home-manager.md)
- [Configuration](docs/configuration.md)
- [Integrations and systemd](docs/integrations.md)
- [CLI and shell completion](docs/cli.md)
- [Manual installation and updates](docs/manual-install.md)
- [Development and release checks](docs/development.md)

## License

[GPL-3.0-only](LICENSE)
