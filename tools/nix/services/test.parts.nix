{
  self,
  lib,
  inputs,
  options,
  ...
}:
{
  imports = [
    inputs.process-compose-flake.flakeModule
  ];

  perSystem =
    {
      config,
      self',
      inputs',
      modos',
      ...
    }:
    let
      makeServices =
        {
          useNix ? false,
        }:
        import ./make-proc-comp-setup.nix {
          inherit
            inputs
            inputs'
            self
            self'
            modos'
            useNix
            ;
        };
    in
    {

      # This test-services are with modos components running from `just run`
      # instead of a Nix derivation.
      # Note: This `test-services-*` package is exposed on `packages` outputs in the flake.
      process-compose."test-services" = makeServices { useNix = false; };
      # This test-services are with all Nix derivations.
      process-compose."test-services-nix" = makeServices { useNix = true; };

      modos.services.config = {
        # Expose to further inspect.
        test-services = config.process-compose.test-services;
        test-services-opts = options;
      };
    };
}
