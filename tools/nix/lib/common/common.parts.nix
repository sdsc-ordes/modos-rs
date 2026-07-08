{ lib, ... }:
{
  flake.lib.common = {
    attrset = import ./attrset.nix { inherit lib; };
    yaml = import ./yaml.nix { inherit lib; };
  };
}
