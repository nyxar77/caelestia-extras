{
  description = "Optional live integrations for Caelestia";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = {self, nixpkgs}: let
    systems = ["x86_64-linux" "aarch64-linux"];
    forAllSystems = nixpkgs.lib.genAttrs systems;
  in {
    packages = forAllSystems (system: {
      default = nixpkgs.legacyPackages.${system}.callPackage ./nix/package.nix {};
    });
    devShells = forAllSystems (system: {
      default = nixpkgs.legacyPackages.${system}.mkShell {
        packages = with nixpkgs.legacyPackages.${system}; [
          go
          gofumpt
        ];
      };
    });
    checks = forAllSystems (system: {
      default = self.packages.${system}.default;
    });
    homeModules.default = import ./nix/home-manager.nix;
  };
}
