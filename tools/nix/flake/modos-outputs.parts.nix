{
  config,
  ...
}:
let
  inherit (config.modos.lib) common;
in
{
  # Inject standalone library into `modos.lib`.
  modos.lib = config.flake.lib;

  # Expose `modos` global namespace on the flake.
  flake.modos = config.modos;

  # Defining a `perSystem` scoped module option `modos`.
  perSystem =
    {
      config,
      ...
    }:
    let
      modos = config.modos;
    in
    {
      modos = {
        # Compute the component flat packages.
        packages.component-flat = common.attrset.flattenDrvs modos.packages.component;
      };

      # Expose the perSystem config on the
      legacyPackages.modos = modos;
      # Expose all packages.
      packages = modos.packages.global // modos.packages.component-flat;
      # Expose all shells.
      devShells = modos.shells;
    };
}
