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
      ...
    }:
    let
      comp = modos.lib.component.getCompFromPath ./.;

      args = {
        inherit modos;
        compName = comp.name;
      };

      plugin = pkgs.callPackage ./plugin args;
    in
    {
      modos.packages.component.${comp.name} = {
        inherit plugin;
      };
    };
}
