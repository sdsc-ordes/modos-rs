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
    {
      ...
    }:
    {
      # The modos library fully instantiated.
      modos.lib = self.lib.mkExtendedLib {
        inherit rootDir rootFileset;
      };
    };
}
