#!/usr/bin/env bash
set -euo pipefail

repo_dir=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
test_dir=$(mktemp -d "${TMPDIR:-/tmp}/caelestia-extras-install-test.XXXXXX")
trap 'rm -rf "$test_dir"' EXIT HUP INT TERM

fake_bin="$test_dir/fake-bin"
mkdir -p "$fake_bin"
mkdir -p "$test_dir/tmp"
cat > "$fake_bin/go" <<'EOF'
#!/bin/sh
set -eu

output=
while [ "$#" -gt 0 ]; do
  case "$1" in
    -o)
      output=$2
      shift
      ;;
  esac
  shift
done

[ -n "$output" ]
if [ "${TEST_INTERRUPT_GO:-}" = 1 ]; then
  kill -INT "$PPID"
fi
printf '%s\n' '#!/bin/sh' 'exit 0' > "$output"
chmod +x "$output"
EOF
chmod +x "$fake_bin/go"

cat > "$fake_bin/systemctl" <<'EOF'
#!/bin/sh
set -eu

printf '%s\n' "$*" >> "$TEST_SYSTEMCTL_LOG"
EOF
chmod +x "$fake_bin/systemctl"

export HOME="$test_dir/home"
export XDG_CONFIG_HOME="$test_dir/config"
export XDG_DATA_HOME="$test_dir/data"
export XDG_STATE_HOME="$test_dir/state"
export XDG_BIN_HOME="$test_dir/bin"
export TMPDIR="$test_dir/tmp"
export TEST_SYSTEMCTL_LOG="$test_dir/systemctl.log"
export PATH="$fake_bin:$PATH"

run_install() {
  "$repo_dir/scripts/install.sh" "$@"
}

run_install install

config_file="$test_dir/config/caelestia-extras/config.toml"
test -x "$test_dir/bin/caelestia-extras"
test -f "$config_file"
test -L "$test_dir/config/caelestia/templates/gtk-portal.css"
grep -qx '  "qssFilePath": "",' "$test_dir/config/caelestia/templates/prismlauncher.json"
test -L "$test_dir/config/caelestia/templates/breeze-caelestia.colors"
test -L "$test_dir/config/systemd/user/caelestia-extras-gtk.path"
test -L "$test_dir/data/themes/Caelestia-Portal/gtk-4.0/base-dark.css"
test -L "$test_dir/data/PrismLauncher/themes/caelestia-breeze/theme.json"

if TEST_INTERRUPT_GO=1 \
  HOME="$test_dir/cancel/home" \
  XDG_CONFIG_HOME="$test_dir/cancel/config" \
  XDG_DATA_HOME="$test_dir/cancel/data" \
  XDG_STATE_HOME="$test_dir/cancel/state" \
  XDG_BIN_HOME="$test_dir/cancel/bin" \
  TMPDIR="$test_dir/tmp" \
  "$repo_dir/scripts/install.sh" >/dev/null 2>&1; then
  printf '%s\n' "installer did not stop on SIGINT" >&2
  exit 1
fi
test ! -e "$test_dir/cancel/bin/caelestia-extras"
test ! -e "$test_dir/cancel/config/caelestia-extras/config.toml"

printf '%s\n' '# user-owned' > "$config_file"
run_install update
grep -qx '# user-owned' "$config_file"

portal_template="$test_dir/config/caelestia/templates/gtk-portal.css"
rm "$portal_template"
printf '%s\n' '/* user-owned */' > "$portal_template"
run_install update
grep -qx '/\* user-owned \*/' "$portal_template"

run_install update --enable gtk
grep -qx -- '--user daemon-reload' "$test_dir/systemctl.log"
grep -qx -- '--user enable --now caelestia-extras-gtk.path' "$test_dir/systemctl.log"
grep -qx -- '--user start caelestia-extras-gtk.service' "$test_dir/systemctl.log"
grep -Fqx "ExecStart=\"$test_dir/bin/caelestia-extras\" --config \"$config_file\" gtk sync" "$test_dir/config/caelestia-extras/managed/systemd/caelestia-extras-gtk.service"

run_install update --enable qt
grep -qx -- '--user enable --now caelestia-extras-qt.path' "$test_dir/systemctl.log"
grep -qx -- '--user start caelestia-extras-qt.service' "$test_dir/systemctl.log"
test -L "$test_dir/config/environment.d/10-caelestia-qt.conf"
grep -Fqx "PathChanged=$test_dir/state/caelestia/theme/breeze-caelestia.colors" "$test_dir/config/caelestia-extras/managed/systemd/caelestia-extras-qt.path"

export CAELESTIA_EXTRAS_SCHEME_FILE="$test_dir/custom/scheme.json"
export CAELESTIA_EXTRAS_THEME_DIR="$test_dir/custom/theme"
export CAELESTIA_EXTRAS_PORTAL_THEME_NAME="Manual-Portal"
run_install update
grep -Fqx "PathChanged=$test_dir/custom/scheme.json" "$test_dir/config/caelestia-extras/managed/systemd/caelestia-extras-gtk.path"
grep -Fqx "PathChanged=$test_dir/custom/theme/gtk-portal.css" "$test_dir/config/caelestia-extras/managed/systemd/caelestia-extras-portal.path"
grep -Fqx "Environment=GTK_THEME=Manual-Portal" "$test_dir/config/caelestia-extras/managed/systemd/xdg-desktop-portal-gtk.service.d/10-caelestia-theme.conf"
test -L "$test_dir/data/themes/Manual-Portal/gtk-3.0/base-dark.css"
