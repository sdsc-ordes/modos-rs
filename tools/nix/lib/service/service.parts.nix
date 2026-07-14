{ config, ... }:
{
  # Service NixOS modules same as `services-flake` useable within an `import` statement in
  # `process-compose-flake`.
  flake.processComposeModules = {
    # garage = import ./garage.nix;

    keycloak = import ./keycloak.nix;
  };
}
