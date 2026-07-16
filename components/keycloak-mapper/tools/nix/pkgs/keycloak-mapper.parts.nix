{
  modos,
  ...
}:
let

in
{
  perSystem =
    {
      pkgs,
      modos',
      ...
    }:
    let
      comp = modos.lib.component.getCompFromPath ./.;

      args = {
        inherit modos;
        compName = comp.name;
        build = modos'.build;
        etcGroupAndPasswd = modos'.packages.image.etcGroupAndPasswd;
      };

      plugin = pkgs.callPackage ./plugin args;
    in
    {
      modos.packages.component.${comp.name} = {
        inherit plugin;
      };
    };
}
