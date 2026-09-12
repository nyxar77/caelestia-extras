# Integrations

## Compositor support

Extras currently supports the upstream Caelestia Hyprland setup. Niri
Caelestia shells are separate community ports, not an upstream backend, so
Extras does not claim Niri support yet.

The code keeps cursor generation and the other shared integrations separate
from compositor actions. A Niri backend can be added when there is a compatible
Caelestia port and a clear way to apply its cursor theme.

## Cursor

`cursor sync` reads Caelestia's active scheme, recolours the Bibata SVGs,
builds a Hyprcursor theme, and applies it through Hyprland. `sync` also builds
the XCursor fallback. The watcher delays that expensive fallback until rapid
wallpaper changes have stopped.

This is the only integration that needs the compositor directly. Today that
backend is Hyprland.

## GTK

`gtk sync` sets the GNOME colour preference and selects the configured light or
dark GTK theme. It uses `dconf`. Home Manager and the manual installer expose
Caelestia's generated `gtk.css` to GTK 3 and `gtk4.css` to GTK 4. Those
stylesheets define the public GTK/libadwaita colour tokens only: no global
`button`, `entry`, `window`, or geometry selectors are used.

GTK 3's built-in Adwaita theme does not use those compatibility names as
palette inputs. The module therefore also exposes `Caelestia-GTK3`, a named
theme for applications that need it. It imports stock Adwaita and overrides
colour properties only; Adwaita still owns geometry, spacing, and widget
structure. Select it for an application with `GTK_THEME=Caelestia-GTK3`.

This is intentionally global at the palette boundary. App-specific styling is
reserved for an application that cannot consume the toolkit tokens correctly.
Already-running GTK processes may need to be reopened after a palette change;
new processes read the generated stylesheet directly.

Caelestia CLI also has a built-in GTK writer. Do not enable both owners: when
using Caelestia's Home Manager module, set
`programs.caelestia.cli.settings.theme.enableGtk = false` before enabling the
GTK integration here. Caelestia continues to render the templates installed by
Extras into its state directory.

## Hyprtoolkit

`hyprtoolkit sync` copies the generated Caelestia Hyprtoolkit config to the
active Hyprtoolkit config path.

## pavucontrol

`pavucontrol` launches the configured `pavucontrol-qt` command. If Caelestia
generated a stylesheet, it is passed to the application.

## imv

Enabling `imv` installs a generated configuration at
`$XDG_CONFIG_HOME/imv/config`. It colours every role exposed by imv: the image
canvas uses `surface`, overlay text uses `onSurface`, and the overlay background
uses `surfaceContainerLowest`. Reopen imv after a scheme change so it reads the
new colours.

imv does not support config includes, so this integration owns the complete
config file. Its command-entry bar is still black and white because imv
hard-codes those colours and exposes no setting for them.

## MPV and ModernZ

Enabling `mpv` installs a Caelestia template for ModernZ and links the rendered
file to `$XDG_CONFIG_HOME/mpv/script-opts/modernz.conf`. MPV reads the new
colours the next time it starts.

This integration does not install MPV or ModernZ. Under Home Manager it checks
`programs.mpv.finalPackage`, so ModernZ must be loaded by the final MPV wrapper.
Use `programs.mpv.scripts = [ pkgs.mpvScripts.modernz ]`, or include ModernZ in
the `scripts` argument of a custom `programs.mpv.package` override.

## LocalSend

LocalSend draws its main interface with Flutter but uses GTK for the native
Linux host window. Its system colour mode already follows the detected desktop
accent. The Home Manager integration keeps that application-owned interface and
wraps only `localsend_app` with the generated `Caelestia-GTK3` named theme, so
the title bar follows Caelestia without enabling a global GTK theme override.

## qBittorrent

`qbittorrent sync` disables qBittorrent's custom UI theme and returns the client
to its native Qt palette, widget style, and icon theme. qBittorrent custom
themes only replace the active Qt palette group, which makes the whole window
change colour on focus-follows-mouse compositors. The Home Manager desktop
entry launches qBittorrent with a scoped qt6ct and Breeze environment, so it
does not depend on stale global session state. It also performs the preference
sync before qBittorrent starts, preventing a previous process from restoring
the broken custom-theme setting during shutdown. A small launcher stylesheet
sets only `QMenuBar` and `QToolBar` foreground/background roles, fixing their
inactive contrast while leaving qBittorrent's tables, forms, dialogs, and
palette under Breeze.

## PrismLauncher

Enabling `prismlauncher` writes a generated Caelestia palette and a small
toolbar-contrast stylesheet into PrismLauncher’s theme directory.
PrismLauncher uses Breeze for the controls; it does not install PrismLauncher.
Select `Caelestia` if it is not already selected. PrismLauncher reads custom
themes at startup, so an already-running window needs one restart.

## Shared Qt theme

Enabling `qt` generates a shared Qt5/Qt6 palette and a Breeze colour scheme.
Breeze draws the base Qt 6 controls. Application-specific stylesheets are used
only for applications whose widget structure requires them; they do not own the
global Qt palette.

## XDG portals

`portal sync` copies the generated GTK theme and Qt palette into portal-specific
locations. GTK and GNOME portal services share the isolated GTK theme. The Qt
screencast picker uses Breeze directly; no portal-wide
stylesheet overrides its native tab panes, frames, or interaction states. The
Home Manager module also installs the service drop-ins needed to keep the
portal processes isolated from the global Qt and GTK settings.

The GTK and GNOME portals remain isolated under their own theme name. They do
not write to normal applications' GTK configuration; the GTK integration owns
the safe global colour-token layer.

## Home Manager and systemd

By default, the Home Manager module runs an aggregate sync during activation
and creates one long-running watcher when at least one enabled integration has
generated files to watch. `syncOnActivation` and `systemd.enable` control those
two actions independently. A 300 ms trailing-edge delay lets Caelestia finish
all generated files before they are copied. Buffered worker channels retain
the newest pending update without launching overlapping processes, and XCursor
is generated only after ten seconds without another change. Starting the
watcher does not rebuild XCursor by itself; activation's aggregate sync
initializes it when that hook is enabled.

For a manual install, run `scripts/install.sh`. It renders the units with the
right binary and XDG paths, then manages them through symlinks. Do not copy the
template files in `systemd/` by hand.

GTK applications that keep an old D-Bus-activated process can be launched
directly with a Home Manager desktop-entry override:

```nix
programs.caelestia-extras.gtk.directLaunch."org.gnome.Nautilus" = {
  name = "Files";
  exec = "nautilus --new-window %U";
  icon = "org.gnome.Nautilus";
};
```
