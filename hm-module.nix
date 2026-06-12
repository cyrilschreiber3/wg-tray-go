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
      defaultText = literalExpression "pkgs.wg-tray-go";
      description = "The wg-tray-go package to use.";
    };

    settings = mkOption {
      type = types.nullOr (types.submodule {
        options = {
          enableNotifications = mkOption {
            type = types.bool;
            default = true;
            description = "Enable notifications for wg-tray-go.";
          };
          tunnelNames = mkOption {
            type = types.listOf types.str;
            default = [];
            description = "List of tunnel names to show in the tray.";
          };
          tunnelGroups = mkOption {
            type = types.listOf (
              types.submodule {
                options = {
                  name = mkOption {
                    type = types.str;
                    description = "Name of the tunnel group.";
                  };
                  pickRandomly = mkOption {
                    type = types.bool;
                    default = false;
                    description = "Whether to pick a random tunnel from this group when bringing the group up.";
                  };
                  tunnelNames = mkOption {
                    type = types.listOf types.str;
                    description = "List of tunnel names in this group.";
                  };
                };
              }
            );
            default = [];
            description = "List of tunnel groups.";
          };
        };
      });
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
    home.packages = [cfg.package];

    xdg.configFile."wg-tray-go/config.json" = lib.mkIf (cfg.settings != null) {
      text = builtins.toJSON {
        enable_notifications = cfg.settings.enableNotifications;
        tunnel_names = cfg.settings.tunnelNames;
        tunnel_groups =
          map (group: {
            name = group.name;
            pick_randomly = group.pickRandomly;
            tunnel_names = group.tunnelNames;
          })
          cfg.settings.tunnelGroups;
      };
    };
  };
}
