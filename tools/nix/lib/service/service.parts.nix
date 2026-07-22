{ ... }:
{
  # Service NixOS modules same as `services-flake` useable within an `import` statement in
  # `process-compose-flake`.
  # FIXME: Once
  # - https://github.com/juspay/services-flake/pull/700
  # - https://github.com/juspay/services-flake/pull/703
  # - https://github.com/juspay/services-flake/pull/702
  # are upstreamed one can remove these modules.
  flake.processComposeModules = {
    keycloak = import ./keycloak.nix;
    authentik = import ./authentik.nix;
    mailhog = import ./mailhog.nix;
    rustfs = import ./rustfs.nix;
  };
}
