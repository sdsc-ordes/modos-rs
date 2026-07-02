{
  self,
  lib,
  ...
}:
let
  fs = lib.fileset;
  rootDir = ../../../..;
  rootFileset = fs.gitTracked rootDir;
in
{
  perSystem =
    { pkgs, ... }:
    let
      buildSystem = pkgs.stdenv.buildPlatform.system;
    in
    {
      # The modos library with root fileset and components.
      modos.lib = self.lib.mkExtendedLib {
        inherit buildSystem rootDir rootFileset;
      };
    };
}
