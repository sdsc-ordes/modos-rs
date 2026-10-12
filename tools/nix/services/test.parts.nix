{
  self,
  lib,
  inputs,
  ...
}:
{
  imports = [
    inputs.process-compose-flake.flakeModule
  ];

  perSystem =
    {
      config,
      inputs',
      modos',
      ...
    }:
    let
      setDataDir =
        procs:
        lib.mapAttrs (
          k: v:
          v
          // {
            dataDir = ".output/run/process-compose/data/${k}";
          }
        ) procs;

      makeServices =
        {
          nix ? false,
        }:
        # Process-compose NixOS module.
        {
          config,
          pkgs,
          pkgsStable,
          ...
        }:
        let
          servicesCfg = config;
          ak = servicesCfg.services.authentik;
        in
        {
          imports = [
            inputs.services-flake.processComposeModules.default
            self.processComposeModules.keycloak
            self.processComposeModules.authentik
            self.processComposeModules.rustfs
            # self.processComposeModules.mailhog

            self.processComposeModules.kms
          ];

          cli.options = {
            unix-socket = "./.output/run/process-compose/pc-test-services.sock";
            log-file = ".output/run/process-compose/log.txt";
            no-server = false;
          };

          settings = {
            log_level = "debug";
            ordered_shutdown = true;
            log_configuration = {
              disable_json = true;
              no_color = true;
              no_metadata = true;
            };
          };

          defaults.processSettings =
            { name, ... }:
            {
              availability.restart = lib.mkDefault "no";
              availability.max_restarts = lib.mkDefault 0;
              log_location = ".output/run/process-compose/log/${name}.log";
              log_configuration = {
                disable_json = true;
                no_color = true;
              };
            };

          services.postgres = setDataDir ak.services.postgres;

          services.authentik = {
            enable = true;
            dataDir = ".output/run/process-compose/data";

            components = inputs'.authentik-nix.legacyPackages.authentikComponents;
            secretKey = "test";

            initialAdminEmail = "admin@example.com";
            initialAdminPassword = "admin";

            server.http.port = 9001;
            worker.http.port = 9002;

            blueprints = {
              export = {
                enable = true;
                path = "tools/configs/authentik/export.yaml";
              };

              imports = [
                ../../configs/authentik/modos-brand-blueprint.yaml
                ../../configs/authentik/modos-groups-users-blueprint.yaml
                ../../configs/authentik/modos-enrollment-blueprint.yaml
                ../../configs/authentik/modos-oauth-blueprint.yaml
              ];
            };
          };

          services.keycloak = {
            enable = true;
            dataDir = ".output/run/process-compose/data";

            settings.http-port = 8081;

            # User is `admin`.
            initialAdminPassword = "admin";

            plugins = [
              modos'.packages.component.keycloak-mapper.plugin
            ];

            realms = {
              modos = {
                path = "./tools/configs/keycloak/modos-realm.json";
                import = true;
                export = true;
              };
            };
          };

          services.mailhog."mailhog" = {
            smtp.port = 1026;
            enable = true;
          };

          services.rustfs = {
            enable = true;
            dataDir = ".output/run/process-compose/data";

            package = inputs'.rustfs-flake.packages.default;
            server.port = 9010;

            # Web Interface
            console = {
              enable = true;
              port = 9011;
            };

            accessKey = "rustfsadmin";
            secretKey = "rustfsadmin";

            logLevel = "debug";

            buckets = [
              "bucket-a"
              "bucket-b"
            ];

            iam.import.path = "tools/configs/rustfs/iam-export";
            iam.export = {
              enable = true;
              path = "tools/configs/rustfs/iam-export";
            };
          };

          services.kms = {
            enable = true;

            package =
              if nix then
                self.packages.kms
              else
                (pkgs.writeShellScriptBin "kms-service" ''
                  just nix::develop "default" just run
                '');
          };

          settings.processes.kms = {
            working_dir = "./components/kms";
            # watch = {
            #   debounce = "300ms";
            #   paths = [
            #     {
            #       path = "./components/kms";
            #     }
            #   ];
            # };
            depends_on = {
              keycloak.condition = "process_healthy";
              authentik.condition = "process_healthy";
              rustfs.condition = "process_healthy";
            };
          };
        };
    in
    {

      # Note: This `test-services-*` package is exposed on `packages` outputs in the flake.
      process-compose."test-services" = makeServices { nix = false; };
      process-compose."test-services-nix" = makeServices { nix = false; };

      modos.services.config = {
        test-services = config.process-compose.test-services;
      };
    };
}
