# Home Manager module

The flake exports `homeManagerModules.default`. `homeModules.default` remains
as a compatibility alias. The module generates the CLI configuration, installs
the required templates, and can run one watcher service. Individual
integrations remain disabled until their `enable` option is set.

## Add the module

Add the input to your flake:

```nix
inputs.caelestia-extras = {
  url = "github:nyxar77/caelestia-extras";
  inputs.nixpkgs.follows = "nixpkgs";
  inputs.home-manager.follows = "home-manager";
};
```

Import the module from a Home Manager module where `inputs` is available:

```nix
{inputs, pkgs, ...}: {
  imports = [inputs.caelestia-extras.homeManagerModules.default];

  programs.caelestia-extras = {
    enable = true;
    autoEnable = true;

    # Disable integrations you do not want:
    # localsend.enable = false;
  };

  # These four integration flags configure launchers and themes; they do not
  # install the applications.
  home.packages = with pkgs; [
    imv
    pavucontrol-qt
    prismlauncher
    qbittorrent
  ];
}
```

If your modules do not already receive `inputs`, pass it with
`home-manager.extraSpecialArgs` or import the module directly from the
`caelestia-extras` output in your flake's `modules` list.

The default paths expect Caelestia to write:

- `$XDG_STATE_HOME/caelestia/scheme.json` for cursor and GTK updates;
- generated integration files under `$XDG_STATE_HOME/caelestia/theme`.

Change `schemeFile` and the relevant `themeDir` options together if Caelestia
uses another location.

`autoEnable` defaults to `true`, like Stylix. Each integration's `enable`
option inherits that value but can be overridden explicitly. To opt into a
small set instead, use:

```nix
programs.caelestia-extras = {
  enable = true;
  autoEnable = false;
  gtk.enable = true;
  imv.enable = true;
};
```

MPV theming is explicitly opt-in because it requires ModernZ in the final MPV
package:

```nix
{
  programs.caelestia-extras.mpv.enable = true;
  programs.mpv = {
    enable = true;
    scripts = [ pkgs.mpvScripts.modernz ];
  };
}
```

### Caelestia CLI GTK ownership

Caelestia CLI can write `gtk-3.0/gtk.css` and `gtk-4.0/gtk.css` itself. The GTK
integration in this module owns those same paths as links to generated
templates. When using Caelestia's Home Manager module, disable its GTK writer:

```nix
programs.caelestia.cli.settings.theme.enableGtk = false;
```

This does not disable Caelestia's template renderer. Files under
`$XDG_CONFIG_HOME/caelestia/templates` are still rendered into
`$XDG_STATE_HOME/caelestia/theme`, which is the output this module consumes.

## What activation changes

With `programs.caelestia-extras.enable = true`, the module:

- installs a wrapped `caelestia-extras` command with its runtime tools;
- writes `$XDG_CONFIG_HOME/caelestia-extras/config.toml`;
- installs only the templates and desktop entries required by enabled
  integrations;
- runs `caelestia-extras sync` during Home Manager activation when
  `syncOnActivation` is enabled and a user D-Bus session is available;
- creates `caelestia-extras-watch.service` when `systemd.enable` is true and at
  least one enabled integration has generated files to watch.

With the default `syncOnActivation = true`, the first graphical activation
needs Caelestia's scheme and the generated files used by enabled integrations.
Set it to `false` when another activation step creates those files later. The
watcher performs an aggregate sync when it starts.

Check the installed service with:

```sh
systemctl --user status caelestia-extras-watch.service
journalctl --user -u caelestia-extras-watch.service -b
caelestia-extras config validate
```

## Avoid conflicting owners

Some integrations deliberately own standard desktop configuration paths. Do
not generate the same paths from another Home Manager module or from Stylix.

| Integration | Paths or settings owned |
| --- | --- |
| GTK | `gtk-3.0/gtk.css`, `gtk-4.0/gtk.css`, the named GTK 3 theme, and GNOME theme settings in dconf |
| Qt | `qt5ct/qt5ct.conf`, `qt6ct/qt6ct.conf`, `environment.d/10-caelestia-qt.conf`, Qt session variables, and the colour scheme, widget style, and icon theme in `kdeglobals` |
| Portal | portal-specific qt6ct configuration and systemd drop-ins for the GTK, GNOME, and Hyprland portal services |
| imv | `imv/config` |
| MPV | `mpv/script-opts/modernz.conf` |
| pavucontrol | the `pavucontrol-qt.desktop` entry |
| qBittorrent | `org.qbittorrent.qBittorrent.desktop` and the appearance keys in `qBittorrent.conf` |
| PrismLauncher | the selected theme directory under `$XDG_DATA_HOME/PrismLauncher/themes` |
| `gtk.directLaunch` | the desktop IDs used as keys in that option |

Caelestia remains the owner of the source templates and generated files under
its own config and state directories. `caelestia-extras` copies or links the
finished files to the consumers above.

The qBittorrent preferences file must be a writable regular file. Sync refuses
to replace a symlink, including a Home Manager-managed `qBittorrent.conf`.

## Option reference

The top-level module defaults to disabled. `autoEnable` and `systemd.enable`
default to `true`, but neither has an effect until the top-level module is
enabled. Individual integration values override `autoEnable`.

### Core

| Option | Default | Meaning |
| --- | --- | --- |
| `programs.caelestia-extras.enable` | `false` | Install the CLI, generated config, and resources for selected integrations |
| `autoEnable` | `true` | Enable integrations by default; individual `enable` values take precedence |
| `package` | repository package | Package installed and wrapped with the required runtime tools |
| `syncOnActivation` | `true` | Run an aggregate sync during activation when user D-Bus is available |
| `systemd.enable` | `true` | Create the watcher when an enabled integration has files to watch |
| `systemd.target` | `config.wayland.systemd.target` | User target that starts and stops the watcher |
| `systemd.environment` | `[]` | Additional `NAME=value` assignments for the watcher service |
| `compositor.backend` | `"hyprland"` | Backend used for compositor-specific actions; `hyprland` is the only accepted value |
| `schemeFile` | `${config.xdg.stateHome}/caelestia/scheme.json` | Caelestia scheme JSON read by cursor and GTK sync |

### Cursor

| Option | Default | Meaning |
| --- | --- | --- |
| `cursor.enable` | `autoEnable` | Enable dynamic Bibata cursor generation |
| `cursor.source` | Bibata's `svg/modern` directory from `pkgs.bibata-cursors.src` | SVG source directory |
| `cursor.buildConfig` | Bibata's `configs/normal/x.build.toml` | Build metadata used for cursor names, hotspots, and animation |
| `cursor.iconDir` | `${config.xdg.dataHome}/icons` | Directory where the generated theme is installed |
| `cursor.theme` | `"Bibata-Caelestia"` | Generated cursor theme name |
| `cursor.size` | `20` | Size passed to `hyprctl setcursor` |
| `cursor.xcursorSizes` | `[20 24 32]` | Sizes generated for the XCursor fallback |
| `cursor.xcursorFallback` | `true` | Build and maintain the XCursor fallback after theme changes |
| `cursor.updateGtk` | `true` | Set the GTK cursor theme and size through dconf |

The Home Manager defaults enable both GTK cursor settings and the XCursor
fallback. The plain TOML defaults are `false`, because a manual installation
must opt into those side effects.

### GTK

| Option | Default | Meaning |
| --- | --- | --- |
| `gtk.enable` | `autoEnable` | Enable GTK light/dark preference and theme synchronization |
| `gtk.darkTheme` | `"adw-gtk3-dark"` | GTK theme selected for a dark Caelestia scheme |
| `gtk.lightTheme` | `"adw-gtk3"` | GTK theme selected for a light Caelestia scheme |
| `gtk.themeDir` | `${config.xdg.stateHome}/caelestia/theme` | Directory containing generated GTK files |
| `gtk.gtk3ThemeName` | `"Caelestia-GTK3"` | Name of the generated stock-Adwaita GTK 3 theme |
| `gtk.directLaunch` | `{}` | Desktop-entry overrides for applications that otherwise reuse a D-Bus-activated process |

Each `gtk.directLaunch.<desktop-id>` entry accepts:

| Field | Required | Meaning |
| --- | --- | --- |
| `name` | yes | Display name |
| `exec` | yes | Command written to the desktop entry |
| `icon` | no | Icon name or path |
| `comment` | no | Desktop-entry comment |
| `genericName` | no | Generic application name |
| `categories` | no | List of desktop categories |
| `mimeType` | no | List of MIME types and URI handlers |
| `startupNotify` | no | Startup-notification setting |
| `settings` | no | Additional string-valued desktop-entry keys |

The module always sets `Terminal=false` and `DBusActivatable=false` for these
overrides. Use the application's existing desktop ID as the attribute name so
the override replaces its normal launcher.

```nix
programs.caelestia-extras.gtk.directLaunch."org.gnome.Nautilus" = {
  name = "Files";
  exec = "nautilus --new-window %U";
  icon = "org.gnome.Nautilus";
  categories = ["GNOME" "GTK" "Utility" "Core"];
  mimeType = ["inode/directory"];
};
```

### Hyprtoolkit and application integrations

| Option | Default | Meaning |
| --- | --- | --- |
| `hyprtoolkit.enable` | `autoEnable` | Install the Caelestia template and copy its generated config when it changes |
| `hyprtoolkit.themeDir` | `${config.xdg.stateHome}/caelestia/theme` | Directory containing generated `hyprtoolkit.conf` |
| `hyprtoolkit.configFile` | `${config.xdg.configHome}/hypr/hyprtoolkit.conf` | Active Hyprtoolkit config destination |
| `imv.enable` | `autoEnable` | Install the generated imv theme as imv's active configuration |
| `imv.themeDir` | `${config.xdg.stateHome}/caelestia/theme` | Directory containing generated `imv.conf` |
| `mpv.enable` | `false` | Install the generated ModernZ theme after checking the final MPV package |
| `mpv.themeDir` | `${config.xdg.stateHome}/caelestia/theme` | Directory containing generated `modernz.conf` |
| `pavucontrol.enable` | `autoEnable` | Install the themed pavucontrol desktop entry and TOML section |
| `pavucontrol.command` | `"pavucontrol-qt"` | Command launched by `caelestia-extras pavucontrol` |
| `localsend.enable` | `autoEnable` | Install a LocalSend wrapper whose GTK host window uses the named GTK 3 theme |
| `localsend.package` | `pkgs.localsend` | LocalSend package to wrap |
| `prismlauncher.enable` | `autoEnable` | Install the generated PrismLauncher theme |
| `prismlauncher.themeDir` | `${config.xdg.stateHome}/caelestia/theme` | Directory containing generated PrismLauncher files |
| `prismlauncher.themeName` | `"caelestia-breeze"` | Directory name under PrismLauncher's custom themes |
| `qbittorrent.enable` | `autoEnable` | Enable qBittorrent preference sync and install its scoped desktop entry |
| `qbittorrent.command` | `"qbittorrent"` | Executable used by the desktop-entry wrapper |
| `qbittorrent.configFile` | `${config.xdg.configHome}/qBittorrent/qBittorrent.conf` | Preferences file updated by the sync |
| `bloom.enable` | `autoEnable` | Generate Bloom's Spicetify palette and keep it synchronized |
| `bloom.themeDir` | `${config.xdg.configHome}/spicetify/Themes/Bloom` | Bloom theme directory containing `color.ini` |

`localsend.enable` requires `gtk.enable`. LocalSend is installed from
`localsend.package`; the imv, MPV, pavucontrol, qBittorrent, and PrismLauncher
flags do not install those applications. `mpv.enable` also requires
`programs.mpv.enable` and a final MPV package containing ModernZ. The imv
integration owns its complete `imv/config` file because imv does not support
including a generated theme fragment; do not also set `programs.imv.settings`
or manage that path elsewhere.

### Qt

| Option | Default | Meaning |
| --- | --- | --- |
| `qt.enable` | `autoEnable` | Enable the shared Qt 5/6 palette and Breeze configuration |
| `qt.themeDir` | `${config.xdg.stateHome}/caelestia/theme` | Directory containing the generated palette and Breeze colour scheme |
| `qt.configHome` | `config.xdg.configHome` | Base configuration directory used by qt5ct and qt6ct |
| `qt.dataHome` | `config.xdg.dataHome` | Base data directory for the generated KDE colour scheme |
| `qt.widgetStyle` | `"Breeze"` | Widget style used by qt6ct and KDE Frameworks applications |
| `qt.iconTheme` | `"Papirus-Dark"` | Icon theme used by qt5ct, qt6ct, and KDE Frameworks applications |

The module installs qt5ct, qt6ct, and Qt 6 Breeze when this integration is
enabled. Qt 6 and KDE Frameworks applications use Breeze; Qt 5 uses Fusion
because the required Qt 5 Breeze plugin is not supplied by the pinned Nixpkgs
package set.

### XDG portals

| Option | Default | Meaning |
| --- | --- | --- |
| `portal.enable` | `autoEnable` | Enable isolated GTK and Qt themes for the portal backends |
| `portal.themeDir` | `${config.xdg.stateHome}/caelestia/theme` | Directory containing generated portal theme files |
| `portal.configHome` | `config.xdg.configHome` | Base directory for portal-specific qt6ct configuration |
| `portal.dataHome` | `config.xdg.dataHome` | Base directory for the private GTK portal theme |
| `portal.themeName` | `"Caelestia-Portal"` | GTK theme name assigned only to the GTK and GNOME portal services |
| `portal.iconTheme` | `"Papirus-Dark"` | Icon theme used by the portal file choosers |

This option configures existing `xdg-desktop-portal-gtk`,
`xdg-desktop-portal-gnome`, and `xdg-desktop-portal-hyprland` services. It does
not install or select the portal backends themselves.
