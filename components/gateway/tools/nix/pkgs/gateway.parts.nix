{
  ...
}:
{
  perSystem =
    { config, pkgs, ... }:
    let
      modos = config.modos;
      comp = modos.lib.component.getCompFromPath ./.;

      args = {
        inherit modos;
        compName = comp.name;
      };

      # The service.
      service = pkgs.callPackage ./service args;

      # The service development build.
      service-dev = pkgs.callPackage ./service (
        args
        // {
          buildType = "debug";
          environmentType = "development";
        }
      );

      # The service image.
      service-image = pkgs.callPackage ./service-image (args // { inherit service; });
    in
    {
      modos.components.packages.${comp.name} = {
        inherit
          service
          service-dev
          service-image
          ;

      };
    };
}
