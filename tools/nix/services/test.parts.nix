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
      servicesCfg = config.process-compose.test-services;
      ak = servicesCfg.services.authentik;

      setDataDir =
        procs:
        lib.mapAttrs (
          k: v:
          v
          // {
            dataDir = ".output/process-compose/data/${k}";
          }
        ) procs;
    in
    {

      # Note: This `test-services` package is exposed on `packages` outputs in the flake.
      process-compose."test-services" =
        # Process-compose NixOS module.
        {
          pkgs,
          pkgsStable,
          ...
        }:
        {
          imports = [
            inputs.services-flake.processComposeModules.default
            self.processComposeModules.keycloak
            self.processComposeModules.authentik
            self.processComposeModules.mailhog
          ];

          cli.options = {
            keep-project = true;
            unix-socket = "./.output/process-compose/pc.sock";
            log-file = ".output/process-compose/log.txt";
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
              log_location = ".output/process-compose/log/${name}.log";
              log_configuration = {
                disable_json = true;
                no_color = true;
              };
            };

          services.postgres = setDataDir ak.services.postgres;

          services.authentik = {
            enable = true;
            dataDir = ".output/process-compose/data";

            components = inputs'.authentik-nix.legacyPackages.authentikComponents;
            secretKey = "test";

            server.http.port = 9001;
            worker.http.port = 9002;

            blueprints = {
              modos = {
                path = ../../configs/authentik/modos-blueprint.yaml;
                import = true;
              };
              enrollment = {
                path = ../../configs/authentik/modos-enrollment-blueprint.yaml;
                import = true;
              };
            };

            path = ../../configs/authentik/modos-enrollment-blueprint.yaml;
          };

          services.keycloak = {
            enable = true;
            dataDir = ".output/process-compose/data";

            settings.http-port = 8081;

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

          services.mailhog = {
            smtp.port = 1026;
            enable = true;
          };
        };

      modos.services.config = {
        test-services = servicesCfg;
      };
    };
}
