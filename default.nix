{
  pkgs ? (
    let
      inherit (builtins) fetchTree fromJSON readFile;
      inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix;
    in
      import (fetchTree nixpkgs.locked) {
        overlays = [
          (import "${fetchTree gomod2nix.locked}/overlay.nix")
        ];
      }
  ),
  buildGoApplication ? pkgs.buildGoApplication,
}:
buildGoApplication {
  pname = "wg-tray-go";
  version = "0.1.0";
  pwd = ./.;
  src = ./.;
  modules = ./gomod2nix.toml;

  # subPackages = ["models" "wgutils" "ui"];

  nativeBuildInputs = with pkgs; (lib.optional pkgs.stdenv.isLinux [
    gcc
    pkg-config
  ]);

  buildInputs = with pkgs;
    [
      wireguard-tools
    ]
    ++ (
      lib.optionals pkgs.stdenv.isLinux [
        libayatana-appindicator
        gtk3
      ]
    );
}
