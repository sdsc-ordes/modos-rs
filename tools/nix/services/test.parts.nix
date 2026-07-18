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
          };

          defaults.processSettings =
            { name, ... }:
            {
              availability.restart = lib.mkDefault "on_failure";
              availability.max_restarts = lib.mkDefault 3;
              log_location = ".output/process-compose/log/${name}.log";
            };

          services.postgres = setDataDir ak.services.postgres;

          services.authentik = {
            enable = true;
            dataDir = ".output/process-compose/data";

            components = inputs'.authentik-nix.packages;
            secretKey = "test";
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
        test = servicesCfg;
      };
    };
}
