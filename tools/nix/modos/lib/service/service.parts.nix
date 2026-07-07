{ ... }:
{
  flake.lib.service = {
    # Service modules for `service-flake`.
    modules = {
      service.garage = import ./garage.nix;
    };
  };
}
