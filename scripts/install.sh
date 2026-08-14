#!/bin/sh
set -eu

umask 022

usage() {
  printf '%s\n' \
    "Usage: scripts/install.sh [install|update] [--enable cursor,gtk,hyprtoolkit,qt,qbittorrent,portal]" \
    "" \
    "Builds the current checkout, installs managed files, and preserves user configuration." \
    "Use --enable after configuring the integrations you want to run." \
    "Set CAELESTIA_EXTRAS_SCHEME_FILE, CAELESTIA_EXTRAS_THEME_DIR, or" \
    "CAELESTIA_EXTRAS_PORTAL_THEME_NAME when the matching config values differ."
}

fail() {
  printf 'error: %s\n' "$*" >&2
  exit 1
}

need() {
  command -v "$1" >/dev/null 2>&1 || fail "required command not found: $1"
}

escape_sed() {
  printf '%s' "$1" | sed 's/[\\&|]/\\&/g'
}

link_managed() {
  managed=$1
  destination=$2

  mkdir -p "$(dirname "$destination")"
  if [ -L "$destination" ] && [ "$(readlink "$destination")" = "$managed" ]; then
    return
  fi
  if [ -e "$destination" ] || [ -L "$destination" ]; then
    printf 'keeping user file: %s\n' "$destination"
    return
  fi
  ln -s "$managed" "$destination"
}

stage_managed() {
  source=$1
  relative=$2

  install -Dm644 "$source" "$staged_managed/$relative"
}

render_staged() {
  source=$1
  relative=$2
  rendered="$work_dir/rendered"

  sed \
    -e "s|@BINARY@|$escaped_binary|g" \
    -e "s|@CONFIG_FILE@|$escaped_config_file|g" \
    -e "s|@CONFIG_HOME@|$escaped_config_home|g" \
    -e "s|@SCHEME_FILE@|$escaped_scheme_file|g" \
    -e "s|@THEME_DIR@|$escaped_theme_dir|g" \
    -e "s|@PORTAL_THEME_NAME@|$escaped_portal_theme_name|g" \
    "$source" > "$rendered"
  install -Dm644 "$rendered" "$staged_managed/$relative"
}

commit_managed() {
  relative=$1
  destination=$2

  managed="$managed_root/$relative"
  install -Dm644 "$staged_managed/$relative" "$managed"
  link_managed "$managed" "$destination"
}

cancel() {
  printf '%s\n' "cancelled before applying changes" >&2
  exit 130
}

: "${HOME:?HOME is required}"

mode=install
enable=
while [ "$#" -gt 0 ]; do
  case "$1" in
    install|update)
      mode=$1
      ;;
    --enable)
      [ "$#" -ge 2 ] || fail "--enable needs a comma-separated list"
      enable=$2
      shift
      ;;
    --help|-h)
      usage
      exit 0
      ;;
    *)
      usage >&2
      fail "unknown argument: $1"
      ;;
  esac
  shift
done

need go
need basename
need dirname
need install
need ln
need mkdir
need mktemp
need readlink
need sed

script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
repo_dir=$(CDPATH= cd -- "$script_dir/.." && pwd)
[ -f "$repo_dir/go.mod" ] || fail "run this script from a caelestia-extras checkout"
cd "$repo_dir"

config_home=${XDG_CONFIG_HOME:-"$HOME/.config"}
data_home=${XDG_DATA_HOME:-"$HOME/.local/share"}
state_home=${XDG_STATE_HOME:-"$HOME/.local/state"}
bin_home=${XDG_BIN_HOME:-"$HOME/.local/bin"}
config_file="$config_home/caelestia-extras/config.toml"
scheme_file=${CAELESTIA_EXTRAS_SCHEME_FILE:-"$state_home/caelestia/scheme.json"}
theme_dir=${CAELESTIA_EXTRAS_THEME_DIR:-"$state_home/caelestia/theme"}
portal_theme_name=${CAELESTIA_EXTRAS_PORTAL_THEME_NAME:-Caelestia-Portal}
binary="$bin_home/caelestia-extras"
managed_root="$config_home/caelestia-extras/managed"
work_dir=$(mktemp -d "${TMPDIR:-/tmp}/caelestia-extras.XXXXXX")
staged_root="$work_dir/staged"
staged_managed="$staged_root/managed"
trap 'rm -rf "$work_dir"' EXIT
trap cancel HUP INT TERM

case "$portal_theme_name" in
  ''|.|..|*/*)
    fail "CAELESTIA_EXTRAS_PORTAL_THEME_NAME must be a theme name, not a path"
    ;;
esac

escaped_binary=$(escape_sed "$binary")
escaped_config_file=$(escape_sed "$config_file")
escaped_config_home=$(escape_sed "$config_home")
escaped_scheme_file=$(escape_sed "$scheme_file")
escaped_theme_dir=$(escape_sed "$theme_dir")
escaped_portal_theme_name=$(escape_sed "$portal_theme_name")

printf '%s\n' "Building and staging files. Press Ctrl-C to cancel without changing your setup."
mkdir -p "$staged_root"
go build -trimpath -buildvcs=false -o "$staged_root/caelestia-extras" ./cmd/caelestia-extras

if [ ! -e "$config_file" ] && [ ! -L "$config_file" ]; then
  install -Dm600 "$repo_dir/assets/manual/config.toml" "$staged_root/config.toml"
  create_config=true
else
  create_config=false
fi

for template in gtk-portal.css gtk-global.css gtk.css pavucontrol-qt.qss prismlauncher.json prismlauncher.qss qt-caelestia.conf breeze-caelestia.colors qbittorrent.qss qbittorrent.json qt6ct-caelestia.conf qt6ct-portal.qss qt6ct-caelestia.qss; do
  stage_managed \
    "$repo_dir/assets/manual/templates/$template" \
    "templates/$template"
done

for version in gtk-3.0 gtk-4.0; do
  stage_managed \
    "$repo_dir/assets/manual/theme/base-dark.css" \
    "theme/$version/base-dark.css"
  stage_managed \
    "$repo_dir/assets/manual/theme/base-light.css" \
    "theme/$version/base-light.css"
done

render_staged \
  "$repo_dir/assets/manual/portal-qt/qt6ct.conf.in" \
  "portal-qt/qt6ct/qt6ct.conf"
for version in qt5ct qt6ct; do
  render_staged \
    "$repo_dir/assets/manual/qt/$version.conf.in" \
    "$version/$version.conf"
done
stage_managed \
  "$repo_dir/assets/manual/environment.d/10-caelestia-qt.conf" \
  "environment.d/10-caelestia-qt.conf"
for dropin in xdg-desktop-portal-gtk.service.d/10-caelestia-theme.conf xdg-desktop-portal-hyprland.service.d/10-caelestia-theme.conf; do
  render_staged \
    "$repo_dir/assets/manual/systemd/$dropin.in" \
    "systemd/$dropin"
done

for source in "$repo_dir"/systemd/*.in; do
  [ -f "$source" ] || continue
  unit=$(basename "$source" .in)
  render_staged \
    "$source" \
    "systemd/$unit"
done

printf '%s\n' "Applying staged files. Cancellation is disabled until this short step finishes."
trap '' HUP INT TERM
install -Dm755 "$staged_root/caelestia-extras" "$binary"
if [ "$create_config" = true ]; then
  install -Dm600 "$staged_root/config.toml" "$config_file"
  printf 'created config: %s\n' "$config_file"
else
  printf 'keeping user config: %s\n' "$config_file"
fi

for template in gtk-portal.css gtk-global.css gtk.css pavucontrol-qt.qss prismlauncher.json prismlauncher.qss qt-caelestia.conf breeze-caelestia.colors qbittorrent.qss qbittorrent.json qt6ct-caelestia.conf qt6ct-portal.qss qt6ct-caelestia.qss; do
  commit_managed \
    "templates/$template" \
    "$config_home/caelestia/templates/$template"
done
link_managed \
  "$theme_dir/prismlauncher.json" \
  "$data_home/PrismLauncher/themes/caelestia-breeze/theme.json"
link_managed \
  "$theme_dir/prismlauncher.qss" \
  "$data_home/PrismLauncher/themes/caelestia-breeze/themeStyle.css"
for version in qt5ct qt6ct; do
  commit_managed \
    "$version/$version.conf" \
    "$config_home/$version/$version.conf"
done
for version in gtk-3.0 gtk-4.0; do
  commit_managed \
    "theme/$version/base-dark.css" \
    "$data_home/themes/$portal_theme_name/$version/base-dark.css"
  commit_managed \
    "theme/$version/base-light.css" \
    "$data_home/themes/$portal_theme_name/$version/base-light.css"
done
commit_managed \
  "portal-qt/qt6ct/qt6ct.conf" \
  "$config_home/portal-qt/qt6ct/qt6ct.conf"
for dropin in xdg-desktop-portal-gtk.service.d/10-caelestia-theme.conf xdg-desktop-portal-hyprland.service.d/10-caelestia-theme.conf; do
  commit_managed \
    "systemd/$dropin" \
    "$config_home/systemd/user/$dropin"
done
for source in "$repo_dir"/systemd/*.in; do
  [ -f "$source" ] || continue
  unit=$(basename "$source" .in)
  commit_managed \
    "systemd/$unit" \
    "$config_home/systemd/user/$unit"
done

printf '%s complete: %s\n' "$mode" "$binary"

if [ -z "$enable" ]; then
  printf '%s\n' "No services enabled. Configure the integrations, then rerun with --enable cursor,gtk,hyprtoolkit,qt,qbittorrent,portal."
  exit 0
fi

need systemctl
"$binary" --config "$config_file" config validate
systemctl --user daemon-reload

old_ifs=$IFS
if [ "$enable" = "all" ]; then
  enable=cursor,gtk,hyprtoolkit,qt,qbittorrent,portal
fi
IFS=,
set -- $enable
IFS=$old_ifs
for integration in "$@"; do
  case "$integration" in
    cursor|gtk|hyprtoolkit|qt|qbittorrent|portal)
      path_unit="caelestia-extras-$integration.path"
      service_unit="caelestia-extras-$integration.service"
      ;;
    *)
      fail "unknown integration in --enable: $integration"
      ;;
  esac
  systemctl --user enable --now "$path_unit"
  systemctl --user start "$service_unit"
  if [ "$integration" = qt ]; then
    link_managed \
      "$config_home/caelestia-extras/managed/environment.d/10-caelestia-qt.conf" \
      "$config_home/environment.d/10-caelestia-qt.conf"
  fi
done

printf '%s\n' "Selected services are enabled. Log out and back in before testing portal changes."
