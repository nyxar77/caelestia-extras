# Integrations

## Cursor

`cursor sync` reads Caelestia's active scheme, recolours the Bibata SVGs,
builds a Hyprcursor theme, and applies it through Hyprland. With
`xcursor_fallback = true`, it also refreshes the XCursor version.

This is the only integration that needs Hyprland directly.

## GTK

`gtk sync` sets the GNOME colour preference and selects the configured light or
dark GTK theme. It uses `dconf`.

## Hyprtoolkit

`hyprtoolkit sync` copies the generated Caelestia Hyprtoolkit config to the
active Hyprtoolkit config path.

## pavucontrol

`pavucontrol` launches the configured `pavucontrol-qt` command. If Caelestia
generated a stylesheet, it is passed to the application.

## XDG portals

`portal sync` copies the generated GTK and Qt files into the portal-specific
locations. The Home Manager module also installs the service drop-ins needed
to keep the portal processes isolated from the global Qt and GTK settings.

## Home Manager and systemd

The Home Manager module creates a service for each enabled sync and a path
watcher for the files that Caelestia updates. It also runs the initial
Hyprtoolkit and portal sync during activation.

For a manual install, the matching units are in `systemd/`. Copy them to
`~/.config/systemd/user/`, make sure `caelestia-extras` is in the service
`PATH`, and enable the `.path` units you need.

GTK applications that keep an old D-Bus-activated process can be launched
directly with a Home Manager desktop-entry override:

```nix
programs.caelestia-extras.gtk.directLaunch."org.gnome.Nautilus" = {
  name = "Files";
  exec = "nautilus --new-window %U";
  icon = "org.gnome.Nautilus";
};
```
