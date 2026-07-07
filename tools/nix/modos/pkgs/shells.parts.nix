# Define different shells.
{
  lib,
  inputs,
  ...
}:
let
in
{
  perSystem =
    {
      config,
      pkgs,
      pkgsStable,
      ...
    }:
    let
      modos = config.modos;

      # Add modos package-namespace to the `nixpkgs` package set.
      # For  `devenv` module functions
      pkgsEx = pkgs // {
        modos = modos.packages;
      };

      # Create a set of devenv modules.
      devenvs = modos.lib.toolchain.createDevenvModules {
        inherit pkgs pkgsStable modos;
      };

      # Define all shells over the set of `devenvs` modules.
      shells = lib.attrsets.mapAttrs (
        name: modules:
        modos.lib.shell.mkShell {
          pkgs = pkgsEx;
          inherit
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
