# Development

The project is a small Go program with a Nix flake.

```sh
go test ./...
nix flake check
tests/install.bash
```

The dev shell provides Go and `gofumpt`:

```sh
nix develop
```

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
