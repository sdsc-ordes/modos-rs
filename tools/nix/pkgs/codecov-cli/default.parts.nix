{ ... }:
{
  perSystem =
    { pkgs, ... }:
    {
      packages = {
        codecov-cli = pkgs.callPackage ./. { };
      };
    };
}
