{
  pkgs,
  lib,
  config,
  ...
}:

let
  cfg = config.services.kms;

  dataDir = cfg.dataDir + "/kms";
  name = "kms";
  inherit (lib) types mkEnableOption mkOption;

in
{
  options.services.kms = {
    enable = mkEnableOption "rustfs";

    dataDir = mkOption {
      type = types.str;
      default = "./data";
      description = ''
        Base directory where keycloak stores its data `<dataDir>/kms`.
      '';
    };

    package = lib.mkOption {
      type = types.either types.str types.package;
      description = ''
        Which package or script of KMS.
      '';
    };

    service = {
      configDir = mkOption {
        type = types.pathWith { };
        default = "./data/config";
        description = ''
          The config directory which contains the config to use for the service.
        '';
      };

      dataDir = mkOption {
        type = types.pathWith { };
        default = "./data";
        description = ''
          The data directory which contains the data for the service.
        '';
      };
    };
  };

  config.settings.processes = {
    ${name} = {
      command = pkgs.writeShellApplication {
        name = "kms";
        text =
          # Bash
          ''
            mkdir -p "${cfg.dataDir}"
            exec ${lib.getExe cfg.package} \
              --config-dir "${cfg.service.configDir}" \
              --data-dir "${cfg.service.dataDir}"

          '';
      };

      # readiness_probe = {
      #   http_get = {
      #     host = cfg.server.host;
      #     port = cfg.server.port;
      #     path = "/health";
      #   };
      #   initial_delay_seconds = 1;
      #   period_seconds = 2;
      #   timeout_seconds = 2;
      #   success_threshold = 1;
      #   failure_threshold = 10;
      # };
    };
  };
}
