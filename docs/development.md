# Development

The project is a Go command with a Nix package, Home Manager module, and a
standalone installer. Changes to generated paths or defaults usually need to be
made in more than one of those surfaces.

## Local checks

```sh
go test ./...
go test -race ./...
go vet ./...
gofumpt -d .
nix flake check
tests/install.bash
```

`nix flake check` builds the Go package, checks that the CLI and Nix package
versions match, evaluates a real Home Manager configuration, and runs the
installer test in a Nix build. Running `tests/install.bash` separately is still
useful while editing the shell script because it is faster and prints each
simulated install/update pass.

The dev shell provides Go and `gofumpt`:

```sh
nix develop
```

Format Nix files with the flake formatter:

```sh
nix fmt
```

## Release checks

Before changing the version or tagging a release:

1. Run every command in the local-check block from a clean checkout.
2. Run `nix flake check --all-systems --no-build` to evaluate every declared
   platform. Build on each architecture that will receive a release artifact.
3. Check `git diff --check` and review the full diff, including `flake.lock`.
4. Update the version in both `cmd/caelestia-extras/main.go` and
   `nix/package.nix`; `caelestia-extras version` and the Nix package must agree.
5. Test one Home Manager activation and one manual update in a real graphical
   session. The automated installer test uses fake `go` and `systemctl`
   commands and cannot prove live portal, dconf, or compositor behavior.

The main pieces are:

- `cmd/caelestia-extras` — CLI and shell completion output
- `internal/config` — TOML loading, defaults, validation, and schema tests
- `internal/compositor` — compositor-specific actions behind a small backend interface
- `internal/cursor` — cursor generation and installation
- `internal/integration` — GTK, Hyprtoolkit, pavucontrol, and portal actions
- `internal/scheme` — Caelestia scheme parsing
- `nix/` — package and Home Manager module
- `systemd/` — standalone user-unit templates
- `scripts/install.sh` — manual install and update entry point

The configuration contract is shared by `internal/config`,
`config/caelestia-extras.schema.json`, `assets/manual/config.toml`, the Home
Manager TOML generator, and `docs/configuration.md`. The schema parity test
catches added or removed sections, but defaults and descriptions still require
a manual review across those files.
