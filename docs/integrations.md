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
Caelestia's generated `gtk.css` to GTK 3 and GTK 4. That stylesheet defines the
public GTK/libadwaita colour tokens only: no global `button`, `entry`, `window`,
or geometry selectors are used.

This is intentionally global at the palette boundary. App-specific styling is
reserved for an application that cannot consume the toolkit tokens correctly.
Already-running GTK processes may need to be reopened after a palette change;
new processes read the generated stylesheet directly.

## Hyprtoolkit

`hyprtoolkit sync` copies the generated Caelestia Hyprtoolkit config to the
active Hyprtoolkit config path.

## pavucontrol

`pavucontrol` launches the configured `pavucontrol-qt` command. If Caelestia
generated a stylesheet, it is passed to the application.

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

`portal sync` copies the generated GTK and Qt files into the portal-specific
locations. The Home Manager module also installs the service drop-ins needed
to keep the portal processes isolated from the global Qt and GTK settings.

The GTK portal remains isolated under its own theme name. It does not write to
normal applications' GTK configuration; the GTK integration owns the safe
global colour-token layer.

## Home Manager and systemd

The Home Manager module creates one long-running watcher and runs an aggregate
sync during activation. A 300 ms trailing-edge delay lets Caelestia finish all
generated files before they are copied. Buffered worker channels retain the
newest pending update without launching overlapping processes, and XCursor is
generated only after ten seconds without another change. Starting the watcher
does not rebuild XCursor by itself; activation's aggregate sync initializes it.

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
