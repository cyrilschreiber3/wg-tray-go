{
  config,
  lib,
  pkgs,
  ...
}:
with lib; let
  cfg = config.programs.wg-tray-go;
in {
  options.programs.wg-tray-go = {
    enable = mkOption {
      type = types.bool;
      default = false;
      description = "Enable wg-tray-go, a system tray application for managing WireGuard tunnels.";
    };

    package = mkOption {
      type = types.package;
      default = pkgs.wg-tray-go;
      defaultText = literalExpression "pkgs.bar-sika";
      description = "The wg-tray-go package to use.";
    };

    settings = mkOption {
      type = types.submodule {
        tunnelNames = types.listOf types.str;
        tunnelGroups = types.listOf (
          types.submodule {
            name = types.str;
            pickRandomly = mkOption {
              type = types.bool;
              default = false;
              description = "Whether to pick a random tunnel from this group when bringing the group up.";
            };
            tunnelNames = types.listOf types.str;
          }
        );
      };
      default = null;
      example = {
        tunnelNames = ["wg0" "wg1" "wg2"];
        tunnelGroups = [
          {
            name = "Work";
            pickRandomly = false;
            tunnelNames = ["wg3" "wg4"];
          }
          {
            name = "Personal";
            pickRandomly = true;
            tunnelNames = ["wg5"];
          }
        ];
      };
      description = "Configuration settings for wg-tray-go.";
    };
  };

  config = mkIf cfg.enable {
    xdg.configFiles."wg-tray-go/config.json".text =
      lib.mkIf cfg.settings
      != null (
        builtins.toJSON {
          tunnel_names = cfg.settings.tunnelNames;
          tunnel_groups =
            map (group: {
              name = group.name;
              pick_randomly = group.pickRandomly;
              tunnel_names = group.tunnelNames;
            })
            cfg.settings.tunnelGroups;
        }
      );
  };
}
