{
  pkgs,
  lib,
  config,
  ...
}:

let
  cfg = config.services.rustfs;

  dataDir = cfg.dataDir + "/rustfs";
  name = "rustfs";
  inherit (lib) types mkEnableOption mkOption;
in
{
  options.services.rustfs = {
    enable = mkEnableOption "rustfs";

    dataDir = mkOption {
      type = types.str;
      default = "./data";
      description = ''
        Base directory where keycloak stores its data `<dataDir>/rustfs`.
      '';
    };

    package = lib.mkOption {
      type = types.package;
      description = ''
        Which package of RustFS to use such as
        'inputs.rustfs.packages.''${pkgs.stdenv.hostPlatform.system}.default'
      '';
    };

    server = {
      host = lib.mkOption {
        type = types.nullOr types.str;
        default = "127.0.0.1";
        description = ''
          The IP interface to bind to.
          `null` means "all interfaces".
        '';
      };

      port = lib.mkOption {
        type = types.port;
        default = 9000;
        description = "The TCP port for the S3 API.";
      };
    };

    console = {
      port = lib.mkOption {
        type = types.port;
        default = 9001;
        description = "The TCP port for the web console.";
      };

      enable = lib.mkOption {
        type = types.bool;
        default = true;
        description = "Enable the console.";
      };
    };

    accessKey = lib.mkOption {
      type = types.str;
      default = "rustfsadmin";
      description = "Access key for authentication (5 to 20 characters).";
    };

    secretKey = lib.mkOption {
      type = types.str;
      default = "rustfsadmin";
      description = "Secret key for authentication (8 to 40 characters).";
    };

    logLevel = lib.mkOption {
      type = lib.types.str;
      default = "info";
      description = "Log level (error, warn, info, debug, trace).";
    };

    extraEnvironment = lib.mkOption {
      type = types.attrsOf types.str;
      default = { };
      description = ''
        Additional environment variables to pass to RustFS.
        See the RustFS documentation for available options
        (e.g. `RUSTFS_CORS_ALLOWED_ORIGINS`, `RUSTFS_TLS_PATH`).
      '';
      example = {
        RUSTFS_OBS_LOGGER_LEVEL = "debug";
        RUSTFS_OBJECT_CACHE_ENABLE = "true";
      };
    };
  };

  config.settings.processes.${name} = {
    environment = {
      RUST_LOG = cfg.logLevel;
      RUSTFS_ADDRESS = "${cfg.server.host}:${lib.toString cfg.server.port}";

      RUSTFS_CONSOLE_ENABLE = lib.boolToString cfg.console.enable;
      RUSTFS_CONSOLE_ADDRESS = "${cfg.server.host}:${lib.toString cfg.console.port}";

      RUSTFS_ACCESS_KEY = cfg.accessKey;
      RUSTFS_SECRET_KEY = cfg.secretKey;

      RUSTFS_DATA_DIR = dataDir;
    }
    // cfg.extraEnvironment;

    command =
      # Bash
      ''
        mkdir -p "$RUSTFS_DATA_DIR"
        exec ${cfg.package}/bin/rustfs server "$RUSTFS_DATA_DIR"
      '';

    readiness_probe = {
      http_get = {
        host = cfg.server.host;
        port = cfg.server.port;
        path = "/health";
      };
      initial_delay_seconds = 1;
      period_seconds = 2;
      timeout_seconds = 2;
      success_threshold = 1;
      failure_threshold = 10;
    };
  };
}
