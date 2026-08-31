{
  description = "Optional live integrations for Caelestia";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  inputs.home-manager = {
    url = "github:nix-community/home-manager";
    inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs =
    {
      self,
      home-manager,
      nixpkgs,
    }:
    let
      systems = [
        "x86_64-linux"
        "aarch64-linux"
      ];
      forAllSystems = nixpkgs.lib.genAttrs systems;
    in
    {
      packages = forAllSystems (system: {
        default = nixpkgs.legacyPackages.${system}.callPackage ./nix/package.nix { };
      });
      devShells = forAllSystems (system: {
        default = nixpkgs.legacyPackages.${system}.mkShell {
          packages = with nixpkgs.legacyPackages.${system}; [
            go
            gofumpt
          ];
        };
      });
      checks = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
          mkHome =
            extras:
            home-manager.lib.homeManagerConfiguration {
              inherit pkgs;
              modules = [
                self.homeManagerModules.default
                {
                  home = {
                    username = "caelestia-extras-test";
                    homeDirectory = "/home/caelestia-extras-test";
                    stateVersion = "25.11";
                  };
                  programs.caelestia-extras = {
                    enable = true;
                    syncOnActivation = false;
                  }
                  // extras;
                }
              ];
            };
          homeManagerTest = mkHome {
            gtk.enable = true;
            systemd = {
              target = "caelestia-test-session.target";
              environment = [ "CAELESTIA_EXTRAS_TEST=1" ];
            };
          };
          homeManagerNoServiceTest = mkHome {
            gtk.enable = true;
            systemd.enable = false;
          };
        in
        {
          default = self.packages.${system}.default;
          home-manager-module =
            assert
              homeManagerTest.config.systemd.user.services.caelestia-extras-watch.Unit.After == [
                "caelestia-test-session.target"
              ];
            assert pkgs.lib.elem "CAELESTIA_EXTRAS_TEST=1"
              homeManagerTest.config.systemd.user.services.caelestia-extras-watch.Service.Environment;
            assert
              !builtins.hasAttr "caelestia-extras-watch" homeManagerNoServiceTest.config.systemd.user.services;
            homeManagerTest.activationPackage;
          installer =
            pkgs.runCommand "caelestia-extras-installer-test"
              {
                nativeBuildInputs = with pkgs; [
                  bash
                  coreutils
                  gnugrep
                  gnused
                ];
                src = ./.;
              }
              ''
                cp -r "$src" source
                chmod -R u+w source
                patchShebangs source/scripts/install.sh source/tests/install.bash
                substituteInPlace source/tests/install.bash \
                  --replace-fail '#!/bin/sh' '#!${pkgs.runtimeShell}'
                cd source
                tests/install.bash
                touch "$out"
              '';
          version =
            pkgs.runCommand "caelestia-extras-version-test"
              {
                nativeBuildInputs = [ self.packages.${system}.default ];
                expectedVersion = self.packages.${system}.default.version;
              }
              ''
                test "$(caelestia-extras version)" = "$expectedVersion"
                touch "$out"
              '';
        }
      );
      formatter = forAllSystems (system: nixpkgs.legacyPackages.${system}.nixfmt);
      homeManagerModules.default = import ./nix/home-manager.nix;
      homeModules.default = self.homeManagerModules.default;
    };
}
