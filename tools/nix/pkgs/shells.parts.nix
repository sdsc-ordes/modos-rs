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
      pkgsStable,
      modos',
      ...
    }:
    let
      # Create a set of devenv modules.
      devenvs = modos.lib.toolchain.createDevenvModules {
        inherit pkgs;
        inherit pkgsStable;
        inherit modos';
      };

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
      ) devenvs;
    in
    {
      modos.shells = shells;
    };
}
