{ ... }:
{
  perSystem =
    {
      config,
      pkgs,
      ...
    }:
    let
      modosLib = config.modos.lib;

      # The tool which acts as a helper to build stuff.
      quitsh = pkgs.callPackage ../../../quitsh/tools/nix/package {
        self = quitsh;
        inherit modosLib;
      };
    in
    {
      modos.pkgs = { inherit quitsh; };
    };
}
