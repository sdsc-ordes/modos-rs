{ ... }:
{

  config = {
    perSystem =
      {
        config,
        pkgs,
        ...
      }:
      let
        modos = config.modos;

        # Add the build support stuff to the library.
        modosLib = modos.lib // {
          inherit (modos) build;
        };

        # Load all components packages by loading
        # `<compPath>/tools/nix/pkgs/default.nix` which must define a
        # an AttrSet of derivations.
        # The result is `{ compName : { pkg-a: derivation-a, pkgs-b: ...}, ...}`
        comp-packages = modosLib.component.loadPackages {
          inherit pkgs modosLib;
        };
        comp-pkgs-flattened = modosLib.common.attrset.flattenDrvs comp-packages;
      in
      {
        # All modos components packages.
        modos = {
          components = {
            packages = comp-packages;
            packages-flat = comp-pkgs-flattened;
          };
        };
      };
  };
}
