{
  description = "Development environment for Go projects";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  inputs.flake-utils.url = "github:numtide/flake-utils";
  inputs.gomod2nix.url = "github:nix-community/gomod2nix";
  inputs.gomod2nix.inputs.nixpkgs.follows = "nixpkgs";
  inputs.gomod2nix.inputs.flake-utils.follows = "flake-utils";

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    gomod2nix,
  }:
    {
      homeManagerModules = {
        wg-tray-go = import ./hm-module.nix;
        default = self.homeManagerModules.wg-tray-go;
      };

      overlays.default = final: prev: {
        wg-tray-go = final.callPackage ./. {
          inherit (gomod2nix.legacyPackages.${final.system}) buildGoApplication;
        };
      };
    }
    // (
      flake-utils.lib.eachDefaultSystem
      (system: let
        pkgs = nixpkgs.legacyPackages.${system};
      in {
        packages = {
          wg-tray-go = pkgs.callPackage ./. {
            inherit (gomod2nix.legacyPackages.${system}) buildGoApplication;
          };
          default = self.packages.wg-tray-go;
        };
        devShells.default = pkgs.callPackage ./shell.nix {
          inherit (gomod2nix.legacyPackages.${system}) mkGoEnv gomod2nix;
        };
      })
    );
}
