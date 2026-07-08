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
buildGoApplication rec {
  prettyPname = "WG Tray Go";
  pname = "wg-tray-go";
  version = "1.2.1";
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

  nativeBuildInputs = with pkgs;
    [
      makeWrapper
    ]
    ++ (lib.optionals stdenv.isLinux [
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

  postPatch = ''
    substituteInPlace wgutils/wg-utils.go \
      --replace-warn '"wg-quick"' '"${pkgs.wireguard-tools}/bin/wg-quick"' \
      --replace-warn '"wg"' '"${pkgs.wireguard-tools}/bin/wg"'
  '';

  postInstall = pkgs.lib.concatStringsSep "\n" [
    (pkgs.lib.optionalString pkgs.stdenv.isLinux ''
      mkdir -p $out/share/icons/hicolor/48x48/apps
      mkdir -p $out/share/icons/hicolor/scalable/apps
      cp ${./assets/icon.png} $out/share/icons/hicolor/48x48/apps/wg-tray-go.png
      cp ${./assets/icon.svg} $out/share/icons/hicolor/scalable/apps/wg-tray-go.svg
    '')
    (pkgs.lib.optionalString pkgs.stdenv.isDarwin ''
      mkdir -p "$out/Applications/${prettyPname}.app/Contents/MacOS"
      mkdir -p "$out/Applications/${prettyPname}.app/Contents/Resources"
      cp ${./assets/icon.icns} "$out/Applications/${prettyPname}.app/Contents/Resources/${pname}.icns"
      cp $out/bin/${pname} "$out/Applications/${prettyPname}.app/Contents/MacOS/${pname}"

      # Create a minimal Info.plist
      cat > "$out/Applications/${prettyPname}.app/Contents/Info.plist" <<EOF
      <?xml version="1.0" encoding="UTF-8"?>
      <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
      <plist version="1.0">
      <dict>
          <key>CFBundleName</key>
          <string>${prettyPname}</string>
          <key>CFBundleDisplayName</key>
          <string>${prettyPname}</string>
          <key>CFBundleIdentifier</key>
          <string>com.cyrilschreiber.${pname}</string>
          <key>CFBundleShortVersionString</key>
          <string>${version}</string>
          <key>CFBundleVersion</key>
          <string>${version}</string>
          <key>CFBundleExecutable</key>
          <string>${pname}</string>
          <key>CFBundleIconFile</key>
          <string>${pname}.icns</string>
          <!-- avoid having a blurry icon and text -->
          <key>NSHighResolutionCapable</key>
          <string>True</string>
          <!-- avoid showing the app on the Dock -->
          <key>LSUIElement</key>
          <string>1</string>
      </dict>
      </plist>
      EOF
    '')
  ];
}
