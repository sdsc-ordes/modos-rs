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
      process-compose."test" =
        # Process-compose NixOS module.
        { pkgs, pkgsStable, ... }:
        {
          imports = [
            inputs.services-flake.processComposeModules.default
            self.processComposeModules.keycloak
          ];

          cli.options = {
            keep-project = true;
            unix-socket = "./.output/process-compose/pc.sock";
          };

          settings = {
            log_level = "debug";
            log_location = ".output/process-compose/log.txt";
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

          services.keycloak.enable = true;
        };

      modos.services.config = {
        test = config.process-compose.test;
      };
    };
}
