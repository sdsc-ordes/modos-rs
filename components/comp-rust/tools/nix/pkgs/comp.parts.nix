{
  modos,
  ...
}:
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
      modos.packages.component.${comp.name} = {
        inherit
          service
          service-dev
          service-image
          ;
      };
    };
}
