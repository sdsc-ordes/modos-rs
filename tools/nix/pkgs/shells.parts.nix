# Define different shells.
{
  lib,
  inputs,
  modos,
  ...
}:
{
  perSystem =
    {
      pkgs,
      modos',
      ...
    }:
    let
      # Define all shells over the set of `devenvs` modules.
      shells = lib.attrsets.mapAttrs (
        name: modules:
        modos.lib.shell.mkShell {
          inherit
            pkgs
            modules
            inputs
            ;
          inherit (pkgs) system;
        }
      ) modos'.devenvModules;
    in
    {
      modos.shells = shells;
    };
}
