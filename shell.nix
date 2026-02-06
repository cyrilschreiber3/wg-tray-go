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
  mkGoEnv ? pkgs.mkGoEnv,
  gomod2nix ? pkgs.gomod2nix,
}: let
  goEnv = mkGoEnv {pwd = ./.;};
in
  with pkgs;
    mkShell rec {
      packages =
        lib.flatten
        [
          [
            # General dependencies
            git

            # gomod2nix prerequisites
            goEnv
            gomod2nix

            # Go development
            delve
            go
            golangci-lint
            golangci-lint-langserver
            gomodifytags
            gopls
            gotests
            impl

            # Project specific dependencies
            wireguard-tools
            gtk3
            gcc
            pkg-config

            (writeScriptBin "create-icons" ''
              #!/usr/bin/env nix-shell
              #!nix-shell -i bash -p inkscape

              set -e

              pushd "$PROJECT_ROOT/assets" >/dev/null

              # Generate PNG
              inkscape "./icon.svg" --export-type=png --export-filename="./icon.png" --export-width=48 --export-height=48 2>/dev/null
              inkscape "./icon.svg" --export-type=png --export-filename="./icon_hd.png" --export-width=1000 --export-height=1000 2>/dev/null

              ${lib.optionalString stdenv.isDarwin ''
                # Generate ICNS
                mkdir -p icon.iconset
                ${
                  lib.concatStringsSep "\n" (
                    map
                    (size: "inkscape \"./icon.svg\" --export-type=png --export-filename=\"icon.iconset/icon_${size}x${size}.png\" --export-width=${size} --export-height=${size} 2>/dev/null")
                    ["16" "32" "64" "128" "256" "512" "1024"]
                  )
                }

                iconutil -c icns icon.iconset -o icon.icns

                rm -rf icon.iconset
              ''}

              popd >/dev/null
            '')
          ]
          (lib.optional stdenv.isLinux [
            libayatana-appindicator
          ])
          (lib.optional stdenv.isDarwin [
            # No additional dependencies for macOS
          ])
        ];

      shellHook = ''

        export LD_LIBRARY_PATH=${lib.makeLibraryPath [stdenv.cc.cc]}
        export PROJECT_ROOT=$(git rev-parse --show-toplevel)

        echo -e "Welcome to the Go dev environment!\n"

        echo -e "$(${go}/bin/go version)\n"

      '';
    }
