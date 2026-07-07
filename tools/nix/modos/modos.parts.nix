{
  lib,
  flake-parts-lib,
  ...
}:
let
  inherit (flake-parts-lib) mkPerSystemOption;

  mkLibOption =
    desc:
    lib.mkOption {
      type = lib.types.lazyAttrsOf lib.types.raw;
      description = desc;
      default = { };
    };

  libDefs = {
    build = mkLibOption "Build support functions and pinned packages.";
    common = mkLibOption "Common helper functions.";
    nixpkgs = mkLibOption "Nixpkgs package import utility.";
    service = mkLibOption "Service modules for 'service-flake'.";
    shell = mkLibOption "Shell utilities.";
    toolchain = mkLibOption "Nix shell modules for 'devenv'.";
  };
in
{
  # The flakes own library.
  options.flake.lib = libDefs;

  # Defining a `perSystem` scoped module option.
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
          # This includes all library functionality which is dependend
          # on `rootDir`, `pkgs`,`pkgsStable`.
          lib = libDefs // {
            image = mkLibOption "Image building functionality.";
            component = mkLibOption "Component scoped helper functions.";
            fileset = mkLibOption "Component scoped filesets utility.";
          };

          packages = lib.mkOption {
            type = lib.types.attrsOf lib.types.package;
            description = "All global exported packages.";
            default = { };
          };

          components = {
            packages = lib.mkOption {
              type = lib.types.attrsOf (lib.types.attrsOf lib.types.package);
              description = "Packages for each component: <compName> -> <pkgs-set>.";
              default = { };
            };

            packages-flat = lib.mkOption {
              type = lib.types.attrsOf lib.types.package;
              description = "Flat packages set for all component (read-only).";
              default = modos.lib.common.attrset.flattenDrvs modos.components.packages;
              readOnly = true;
            };
          };

          shells = lib.mkOption {
            type = lib.types.attrsOf lib.types.package;
            description = "All development Nix shells.";
            default = { };
          };

          debug = lib.mkOption {
            type = lib.types.raw;
            description = "Debugging Nix evaluations with free-form stuff just for `nix repl`.";
            internal = true;
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
        // modos.components.packages-flat # All component packages.
        // modos.packages; # All other packages.

        # Expose all shells.
        devShells = modos.shells;
      };
    }
  );
}
