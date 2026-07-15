{ config, ... }:
{
  # Service NixOS modules same as `services-flake` useable within an `import` statement in
  # `process-compose-flake`.
  flake.processComposeModules = {
    keycloak = import ./keycloak.nix;
    mailhog = import ./mailhog.nix;
  };
}
