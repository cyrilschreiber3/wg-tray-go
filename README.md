# wg-tray-go

wg-tray-go is a WireGuard tunnel manager for your system tray. It uses native platform components and is designed to work on Linux and macOS.

## Features

- Manage WireGuard tunnels from the tray
- Organize tunnels into groups
- Connect to a random or all tunnels in a group
- Native look and feel on Linux and macOS
- Home Manager module for Nix-based setups

## Installation

### Nix / NixOS

If you use Nix or Home Manager, you can install wg-tray-go from the flake.

Example:

```nix
{
  inputs.wg-tray-go.url = "github:cyrilschreiber3/wg-tray-go";

  outputs = { self, nixpkgs, wg-tray-go, ... }:
  {
    homeConfigurations."your-user" = home-manager.lib.homeManagerConfiguration {
      pkgs = import nixpkgs { system = "x86_64-linux"; };
      modules = [
        wg-tray-go.homeManagerModules.default
      ];
    };
  };
}
```

### Linux

You have two options on standard Linux systems:

1. Download the executable tarball

Download the archive that contains only the binary, extract it, and run it directly:

```bash
tar -xzf wg-tray-go_<version>_linux_amd64.tar.gz
./wg-tray-go
```

2. Download the package tarball

Download the `-package` tarball if you want the executable plus desktop integration files, then install them into the appropriate locations.

Example layout:

- executable -> ~/.local/bin/wg-tray-go or /usr/local/bin/wg-tray-go
- desktop file -> ~/.local/share/applications/wg-tray-go.desktop
- icon files -> ~/.local/share/icons/hicolor/...

Example install:

```bash
mkdir -p ~/.local/bin ~/.local/share/applications ~/.local/share/icons
cp bin/wg-tray-go ~/.local/bin/
cp share/applications/wg-tray-go.desktop ~/.local/share/applications/
cp -r share/icons ~/.local/share/
```

### MacOS (arm / Apple Silicon)

You also have two options on macOS:

1. Download the executable tarball

Download the archive with the executable and run it directly:

```bash
tar -xzf wg-tray-go_<version>_darwin_arm64.tar.gz
./wg-tray-go
```

2. Download the DMG

Download the DMG image, open it, and drag the app into Applications like a normal MacOS app.

### Build from source

Because wg-tray-go is written in Go, it is easy to build locally.

On Linux, you may need a few dependencies first. Here is an example for Debian/Ubuntu:

```bash
sudo apt-get install libayatana-appindicator3-dev libgtk-3-dev pkg-config
```

Example build:

```bash
go build -o wg-tray-go ./...
```

This can be useful for platforms not covered by CI builds.

## Configuration

wg-tray-go supports a JSON configuration file.
The purpose of this config file is to avoid using root access on `ls` to read in `/etc`; however, its correctness is not checked and it is your responsibility to keep it up to date. If no config is provided, the interfaces are dynamically looked up in `/etc/wireguard`.

If you want to avoid being prompted for your root password each time you run `wg-tray-go`, use a config file and add the following lines to your sudoers file (`sudo visudo` to edit), replacing `<username>` with your user name:

```text
<username> ALL=(ALL) NOPASSWD: /usr/bin/wg
<username> ALL=(ALL) NOPASSWD: /usr/bin/wg-quick
```

### JSON configuration

The configuration file is located at:

```text
~/.config/wg-tray-go/config.json
```

Configuration options:

- `enable_notifications`: boolean, whether to enable notifications when tunnels are connected or disconnected.
- `tunnel_names`: array of strings, names of the WireGuard tunnels to display at top level.
- `tunnel_groups`: array of objects, each representing a group of tunnels. Each object has:
  - `name`: string, name of the group.
  - `pick_randomly`: boolean, whether to pick a random tunnel instead of connecting to all tunnels in the group when using the group shortcut.
  - `tunnel_names`: array of strings, names of the WireGuard tunnels in this group.

Note:

- The tunnel names must match the names of the WireGuard tunnels as defined in your system's WireGuard configuration.
- A tunnel name can be used in multiple groups.

Example:

```json
{
  "enable_notifications": true,
  "tunnel_names": ["wg0", "wg1", "wg2"],
  "tunnel_groups": [
    {
      "name": "Work",
      "pick_randomly": false,
      "tunnel_names": ["wg3", "wg4"]
    },
    {
      "name": "Personal",
      "pick_randomly": true,
      "tunnel_names": ["wg5"]
    }
  ]
}
```

### Nix / Home Manager

If you use Home Manager, you can configure wg-tray-go declaratively:

```nix
{
    programs.wg-tray-go = {
        enable = true;
        settings = {
            enableNotifications = true;
            tunnelNames = [ "wg0" "wg1" "wg2" ];
            tunnelGroups = [
                {
                    name = "Work";
                    pickRandomly = false;
                    tunnelNames = [ "wg3" "wg4" ];
                }
                {
                    name = "Personal";
                    pickRandomly = true;
                    tunnelNames = [ "wg5" ];
                }
            ];
        };
    };
}
```

## Notes

Linux packaging may vary depending on your distribution and desktop environment.
The app uses native components, so behavior and appearance may differ slightly between platforms.
