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

  desktopItems = [
    (makeDesktopItem {
      name = "wg-tray-go";
      desktopName = "wg-tray-go";
      exec = "wg-tray-go %U";
      icon = "wg-tray-go";
      comment = "System tray application for managing WireGuard tunnels";
      categories = ["Network" "Utility"];
      startupWMClass = "wg-tray-go";
    })
  ];

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
# TODO: create .app bundle for macOS

