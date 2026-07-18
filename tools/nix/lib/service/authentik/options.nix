{
  lib,
  pkgs,
  ...
}:

let
  inherit (lib)
    mkOption
    types
    ;
  name = "authentik";

  settingsFormat = pkgs.formats.yaml { };

  # A relative user-provided path, or a Nix store path (same pattern as keycloak realms).
  blueprintPath = types.nullOr (
    types.either (types.pathWith {
      inStore = false;
      absolute = false;
    }) (types.pathWith { inStore = true; })
  );

  hostAndPort = name: port: {
    host = mkOption {
      type = types.str;
      default = "0.0.0.0";
      description = "Host of the ${name}.";
    };
    port = mkOption {
      type = types.number;
      default = port;
      description = "Port of the ${name}.";
    };
  };
in
{
  options.services.authentik = {
    enable = mkOption {
      type = types.bool;
      default = false;
      example = true;
      description = ''
        Whether to enable the Keycloak identity and access management
        server.
      '';
    };
    dataDir = mkOption {
      type = types.str;
      default = "./data";
      description = ''
        Base directory where keycloak stores its data `<dataDir>/keycloak`.
      '';
    };
    # Authentik is not in nixpkgs, so the package set has to be provided by the user
    # from their own `authentik-nix` flake input. This keeps services-flake input-free.
    components = mkOption {
      type = types.attrsOf types.package;
      example = lib.literalExpression "inputs.authentik-nix.packages.\${system}";
      description = ''
        The Authentik component packages, typically
        `inputs.authentik-nix.packages.''${system}`.

        Must provide at least the following attributes:
        - `gopkgs`            – provides `bin/server`
        - `rust`              – provides `bin/authentik` (the `worker` subcommand)
        - `migrate`           – provides `bin/migrate.py`
        - `staticWorkdirDeps` – the working-directory dependencies
                                (`authentik/`, `blueprints/`, `templates/`, static assets)
        - `manage`            – the management CLI (optional, for blueprint tooling)
      '';
    };

    secretKey = mkOption {
      type = types.nullOr types.str;
      example = "insecure-dev-secret";
      description = ''
        The Authentik secret key, exported as `AUTHENTIK_SECRET_KEY`.

        This is written into the process environment and therefore ends up in the
        Nix store; only use it for local development. For anything else set
        {option}`environmentFile` instead and put `AUTHENTIK_SECRET_KEY` there.
      '';
    };

    initialAdminEmail = mkOption {
      type = types.str;
      default = "admin@example.com";
      description = ''
        E-mail for the bootstrap `akadmin` user, exported as
        `AUTHENTIK_BOOTSTRAP_EMAIL`.
      '';
    };

    initialAdminPassword = mkOption {
      type = types.str;
      default = "admin";
      description = ''
        Initial password for the bootstrap `akadmin` user, exported as
        `AUTHENTIK_BOOTSTRAP_PASSWORD`.
      '';
    };

    environmentFile = mkOption {
      type = types.nullOr (
        types.pathWith {
          inStore = false;
          absolute = false;
        }
      );
      default = null;
      example = "authentik.env";
      description = ''
        Path to an environment file with additional
        `AUTHENTIK_*` variables, e.g. `AUTHENTIK_SECRET_KEY` and
        `AUTHENTIK_POSTGRESQL__PASSWORD`.
        Values here override {option}`settings`.
      '';
    };

    services = {
      postgres = mkOption {
        type = types.attrsOf types.raw;
        readOnly = true;
        description = ''
          The config to easily define the needed postgres process.
        '';
      };
    };

    logLevel = mkOption {
      type = types.str;
      default = "info";
      example = "debug";
      description = "Authentik log level.";
    };

    server = {
      http = hostAndPort "server endpoint" 9000;
      https = hostAndPort "server endpoint" 9443;
      metrics = hostAndPort "server metrics endpoint." 9300;
    };

    worker = {
      http = hostAndPort "worker endpoint" 9001;
      metrics = hostAndPort "worker metrics endpoint." 9302;
    };

    email = hostAndPort "${name}'s email connection" 25;

    postgres = (hostAndPort "${name}'s postgres DB" 5432) // {
      name = mkOption {
        type = types.str;
        default = "authentik";
        description = "PostgreSQL database name (`postgresql.name`).";
      };

      user = mkOption {
        type = types.str;
        default = "authentik";
        description = "PostgreSQL user (`postgresql.user`).";
      };

      password = mkOption {
        type = types.str;
        default = "authentik";
        description = ''
          PostgreSQL password (`postgresql.password`). Written to the config
          file in the Nix store; for non-dev use, override it via
          {option}`environmentFile` (`AUTHENTIK_POSTGRESQL__PASSWORD`).
        '';
      };
    };

    blueprints = mkOption {
      default = { };
      type = types.attrsOf (
        types.submodule {
          options = {
            path = mkOption {
              type = blueprintPath;
              default = null;
              example = "./blueprints/my-blueprint.yaml";
              description = ''
                Path (relative to the `process-compose` working dir, or a Nix store
                path) of a blueprint YAML file to make available for import.
              '';
            };

            import = mkOption {
              type = types.bool;
              default = true;
              description = "Whether to make this blueprint available for import.";
            };
          };
        }
      );

      example = lib.literalExpression ''
        {
          my-app = {
            path = ./blueprints/my-app.yaml;
          };
        }
      '';

      description = ''
        Blueprints to import on start up.
        Enabled blueprints are copied into the blueprints directoryr and
        auto-applied by the Authentik worker.
      '';
    };

    settings = mkOption {
      description = "YAML option for authentic which are merged with '<dataDir>/authentic/lib/default.yml'.";
      type = types.submodule {
        freeformType = settingsFormat.type;
        options = { };
      };
    };
  };
}
