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
    { config, ... }:
    {
      process-compose."test-services" =
        # Process-compose NixOS module.
        { pkgs, pkgsStable, ... }:
        {
          imports = [
            inputs.services-flake.processComposeModules.default
            self.processComposeModules.keycloak
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
              namespace = lib.mkDefault "services-test";
              availability.restart = lib.mkDefault "on_failure";
              availability.max_restarts = lib.mkDefault 3;
              log_location = ".output/process-compose/log/${name}.log";
            };

          services.keycloak = {
            enable = true;
            dataDir = ".output/process-compose/data";

            settings.http-port = 8081;

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
        test = config.process-compose.test-services;
      };
    };
}
