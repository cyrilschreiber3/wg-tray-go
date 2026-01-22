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
    mkShell {
      packages =
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
        ]
        ++ lib.optional stdenv.isLinux [
          libayatana-appindicator
        ];

      shellHook = ''

        echo -e "Welcome to the Go dev environment!\n"

        echo -e "$(${go}/bin/go version)\n"

      '';
    }
