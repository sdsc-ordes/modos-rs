{ modos, ... }:
{
  perSystem =
    {
      pkgs,
      ...
    }:
    let
      quitsh = pkgs.callPackage ./. {
        self = quitsh;
        inherit modos;
      };
    in
    {
      modos.packages.global = { inherit quitsh; };
    };
}
