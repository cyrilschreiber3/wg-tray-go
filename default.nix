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
  makeDesktopItem ? pkgs.makeDesktopItem,
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

  nativeBuildInputs = with pkgs; (lib.optionals stdenv.isLinux [
    gcc
    pkg-config
    copyDesktopItems
  ]);

  buildInputs = with pkgs;
    [
      wireguard-tools
    ]
    ++ (
      lib.optionals stdenv.isLinux [
        libayatana-appindicator
        gtk3
      ]
    );

  postInstall = pkgs.lib.optionalString pkgs.stdenv.isLinux ''
    mkdir -p $out/share/icons/hicolor/48x48/apps
    mkdir -p $out/share/icons/hicolor/scalable/apps
    cp ${./assets/icon.png} $out/share/icons/hicolor/48x48/apps/wg-tray-go.png
    cp ${./assets/icon.svg} $out/share/icons/hicolor/scalable/apps/wg-tray-go.svg
  '';
}
# TODO: create .app bundle for macOS

