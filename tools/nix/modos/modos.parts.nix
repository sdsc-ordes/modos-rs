{
  flake-parts-lib,
  ...
}:
let
  inherit (flake-parts-lib) mkPerSystemOption;
in
{
  # Defining a `perSystem` scoped NixOS module-system option.
  options.perSystem = mkPerSystemOption (
    {
      config,
      lib,
      ...
    }:
    let
      modos = config.modos;
    in
    {
      # Define a new `modos` namespace options available with
      # `config.modos` inside `perSystem = {config,...}`.
      # The options are collected in different flake-parts modules.
      # and exported below into the flake.
      options = {
        modos = {
          lib = lib.mkOption {
            type = lib.types.lazyAttrsOf lib.types.raw;
            description = "The library functions.";
            default = { };
          };

          build = lib.mkOption {
            type = lib.types.lazyAttrsOf lib.types.raw;
            description = "The build functionality.";
            default = { };
          };

          pkgs = lib.mkOption {
            type = lib.types.attrsOf lib.types.package;
            description = "All packages.";
            default = { };
          };

          components = {
            packages = lib.mkOption {
              type = lib.types.raw;
              description = "Packages per component.";
              default = { };
            };
            packages-flat = lib.mkOption {
              type = lib.types.attrsOf lib.types.package;
              description = "Flat packages set per component.";
              default = [ ];
            };
          };

          shells = lib.mkOption {
            type = lib.types.attrsOf lib.types.package;
            description = "All development Nix shells.";
            default = { };
          };
        };
      };

      config = {
        # Expose the namespace in `legacyPackages`
        # just for debugging.
        legacyPackages.modos = modos;

        # Expose all packages.
        packages = {
          inherit (config.modos) quitsh;
        }
        // modos.components.packages-flat
        // modos.pkgs;

        # Expose all shells.
        devShells = modos.shells;
      };
    }
  );
}
