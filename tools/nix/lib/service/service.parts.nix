{ ... }:
{
  flake.lib.service = {
    # Service modules for `services-flake`.
    modules = {
      garage = import ./garage.nix;
    };
  };
}
