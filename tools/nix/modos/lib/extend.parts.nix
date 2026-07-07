{
  config,
  lib,
  ...
}:
let
  # Make the full modos library,
  # by instantiating it with `rootDir` and `rootFileset`.
  # The output is a merge of `lib` with additional functionality
  # like `components` and `fileset` functionality.
  mkExtendedLib =
    {
      rootDir,
      rootFileset,
    }:
    let
      # The libraries only belonging to this flake,
      # since `rootDir` is baked.
      component = (import ./component) {
        inherit lib rootDir;
        libCommon = config.flake.lib.common;
      };

      fileset = (import ./fileset.nix) {
        inherit
          lib
          rootDir
          rootFileset
          ;
        libComponent = component;
      };
    in
    {
      # Redefine libs.
      inherit (config.flake.lib)
        build
        common
        nixpkgs
        service
        shell
        toolchain
        ;

      inherit component fileset;
    };

in
{
  flake.lib = {
    inherit mkExtendedLib;
  };
}
