{ ... }:
{
  perSystem =
    {
      config,
      pkgs,
      ...
    }:
    let
      quitsh = pkgs.callPackage ./. {
        self = quitsh;
        inherit (config) modos;
      };
    in
    {
      modos.packages = { inherit quitsh; };
    };
}
