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

  provisionEnable = cfg.buckets != [ ] || cfg.iam.import.path != null || cfg.provisionScript != null;
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

    region = mkOption {
      type = types.str;
      default = "us-east-1";
      description = "The service region reported to clients.";
    };

    buckets = lib.mkOption {
      type = types.listOf types.str;
      default = [ ];
      description = "Buckets to create on startup.";
      example = [
        "uploads"
        "assets"
      ];
    };

    iam = {
      import = {
        path = mkOption {
          type = types.nullOr (
            types.either
              (types.pathWith {
                inStore = false;
                absolute = false;
              })
              # A nix store path.
              (types.pathWith { inStore = true; })
          );
          default = null;
          description = ''
            Path to the folder from RustFS IAM export (unzipped) to restore via the admin `import-iam` endpoint on startup.
            Produce it via the console IAM export tab.
            Import is get-or-create, so it is safe to re-apply on an already-populated data dir.
          '';
        };
      };

      export = {
        enable = mkEnableOption "export of IAM settings on a process '${name}-iam-export'.";

        path = mkOption {
          type = types.pathWith {
            inStore = false;
            absolute = false;
          };
          default = "${config.dataDir}/export/iam-settings";
          description = ''
            Path to the folder where to unzip the RustFS IAM export when the
            manual process '${name}-iam-export runs'.
          '';
        };
      };
    };

    provisionScript = mkOption {
      type = types.nullOr types.package;
      default = null;
      description = ''
        Extra provision script with custom provisioning steps.
      '';
    };
  };

  config.settings.processes = {
    ${name} = {
      environment = {
        RUST_LOG = cfg.logLevel;
        RUSTFS_ADDRESS = "${cfg.server.host}:${lib.toString cfg.server.port}";
        RUSTFS_CONSOLE_ENABLE = lib.boolToString cfg.console.enable;
        RUSTFS_CONSOLE_ADDRESS = "${cfg.server.host}:${lib.toString cfg.console.port}";

        RUSTFS_ACCESS_KEY = cfg.accessKey;
        RUSTFS_SECRET_KEY = cfg.secretKey;

        RUSTFS_DATA_DIR = "${dataDir}/data";

        RUSTFS_REGION = cfg.region;
      }
      // cfg.extraEnvironment;

      command = pkgs.writeShellApplication {
        name = "rustfs";
        text =
          # Bash
          ''
            mkdir -p "$RUSTFS_DATA_DIR"
            exec ${cfg.package}/bin/rustfs server "$RUSTFS_DATA_DIR"
          '';
      };

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
  // lib.optionalAttrs provisionEnable {
    "${name}-provision" = {
      command = pkgs.writeShellApplication {
        name = "rustfs-provision";
        runtimeInputs = [
          pkgs.curl
          pkgs.awscli2
          pkgs.zip
        ];
        text =
          # Bash
          ''
            # shellcheck disable=SC2034
            endpoint="${cfg.server.host}:${lib.toString cfg.server.port}"

            # Scratch dir (for the IAM zip); nothing is written to $HOME.
            tmp="$(mktemp -d)"
            trap 'rm -rf "$tmp"' EXIT

            export AWS_ACCESS_KEY_ID="${cfg.accessKey}"
            export AWS_SECRET_ACCESS_KEY="${cfg.secretKey}"
            export AWS_DEFAULT_REGION="${cfg.region}"
          ''
          + lib.concatStringsSep "\n" (
            lib.map (
              b:
              # Bash
              ''
                echo "Provision: Ensuring bucket '${b}'."
                aws --endpoint-url "http://$endpoint" s3 mb "s3://${b}" 2>/dev/null
                echo "Provision: Bucket '${b}' created."
              '') cfg.buckets
          )
          + (lib.optionalString (cfg.iam.import.path != null) ''
            echo "Provision: Importing IAM from zipping '${cfg.iam.import.path}'"
            zip -rq "$tmp/iam.zip" "${cfg.iam.import.path}"

            curl -fsS -X PUT \
              --aws-sigv4 "aws:amz:${cfg.region}:s3" \
              -u "${cfg.accessKey}:${cfg.secretKey}" \
              --data-binary "@$tmp/iam.zip" \
              -H "Content-Type: application/zip" \
              "http://$endpoint/rustfs/admin/v3/import-iam"

            echo "Provision: IAM import done."
          '')
          + (lib.optionalString (cfg.provisionScript != null) "${lib.getExe cfg.provisionScript}")
          + ''
            echo "Provision: Done."
          '';
      };

      depends_on.${name}.condition = "process_healthy";
      availability.restart = "no";
    };
  }
  // lib.optionalAttrs cfg.iam.export.enable {
    "${name}-iam-export" = {
      command = pkgs.writeShellApplication {
        name = "${name}-iam-export";
        runtimeInputs = [
          pkgs.curl
          pkgs.unzip
        ];
        text =
          # Bash
          ''
            endpoint="${cfg.server.host}:${lib.toString cfg.server.port}"

            tmp="$(mktemp -d)"
            trap 'rm -rf "$tmp"' EXIT

            # IAM export — SigV4-signed GET from the admin export endpoint.
            echo "Export: Downloading IAM settings into '${cfg.iam.export.path}'."
            curl -fsS -X GET \
              --aws-sigv4 "aws:amz:${cfg.region}:s3" \
              -u "${cfg.accessKey}:${cfg.secretKey}" \
              -H "Accept: application/zip" \
              -o "$tmp/iam.zip" \
              "http://$endpoint/rustfs/admin/v3/export-iam"

            echo "Unzipping into '${cfg.iam.export.path}'."
            mkdir -p "${cfg.iam.export.path}"
            unzip -oq "$tmp/iam.zip" -d "${cfg.iam.export.path}"

            echo "Export: IAM export done."
          '';
      };
      disabled = true;
    };
  };
}
