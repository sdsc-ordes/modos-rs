{
  config,
  lib,
  flake-parts-lib,
  ...
}:
let
  configG = config;
  inherit (flake-parts-lib) mkPerSystemOption;

  mkLibOption =
    desc:
    lib.mkOption {
      type = lib.types.lazyAttrsOf lib.types.raw;
      description = desc;
      default = { };
    };

  libDefs = {
    build = mkLibOption "Build functions.";
    common = mkLibOption "Common helper functions.";
    component = mkLibOption "Component helper functions.";
    nixpkgs = mkLibOption "Nixpkgs import utility.";
    service = mkLibOption "Service modules for 'service-flake'.";
    shell = mkLibOption "Nix dev. shell utilities.";
    toolchain = mkLibOption "Toolchain utilities for defining shells.";
  };

  libModosDefs = libDefs // {
    repo = mkLibOption "Repository definintions.";
    component = mkLibOption "Component scoped helper functions.";
    fileset = mkLibOption "Component scoped filesets utilities.";
  };

  packageSetType = lib.types.attrsOf (lib.types.either lib.types.package packageSetType);
in
{
  options = {
    # Define the flakes own standalone reusable library
    # which is not repository/system/pkgs depedend.
    flake.lib = libDefs;

    # Define the modos namespace with repository scoped functionality.
    # Include libs in `flake.lib`.
    modos.lib = libModosDefs;

    # Defining a `perSystem` scoped module option `modos`.
    perSystem = mkPerSystemOption (
      {
        config,
        lib,
        ...
      }:
      let
        cfg = config.modos;
      in
      {
        # Define a new `modos` namespace options available with
        # `config.modos` inside `perSystem = {config,...}`.
        # The options are collected in different flake-parts modules.
        # and exported below into the flake.
        options.modos = {
          build = mkLibOption "Build functionality and pinned packages.";

          # This includes all library functionality which is dependent
          # on `rootDir`, `pkgs`,`pkgsStable`.
          packages = {
            global = lib.mkOption {
              type = lib.types.attrsOf lib.types.package;
              description = "All global exported packages.";
              default = { };
            };

            pinned = lib.mkOption {
              type = packageSetType;
              description = "Pinned nested packages sets.";
              default = { };
            };

            image = lib.mkOption {
              type = lib.types.attrsOf lib.types.package;
              description = "Packages for images.";
              default = { };
            };

            component = lib.mkOption {
              type = lib.types.attrsOf (lib.types.attrsOf lib.types.package);
              description = "Packages for each component: <compName> -> <pkgs-set>.";
              default = { };
            };

            component-flat = lib.mkOption {
              type = lib.types.attrsOf lib.types.package;
              description = "Flat packages set for all component (read-only).";
              default = configG.modos.lib.common.attrset.flattenDrvs cfg.packages.component;
              readOnly = true;
            };
          };

          shells = lib.mkOption {
            type = lib.types.attrsOf lib.types.package;
            description = "All development Nix shells.";
            default = { };
          };
        };

        config = {
          legacyPackages.modos = cfg;

          # Expose all packages.
          packages = cfg.packages.global // cfg.packages.component-flat;

          # Expose all shells.
          devShells = cfg.shells;
        };
      }
    );

    # Define the modos namespace as a flake output.
    flake.modos = lib.mkOption {
      type = lib.types.raw;
      readOnly = true;
    };
  };

  config = {
    # Inject standalone library into `modos.lib`.
    modos.lib = configG.flake.lib;

    # Expose `modos` global namespace on the flake.
    flake.modos = configG.modos;

    # Make `modos` function parameter available on flake-parts modules.
    # and make `modos'` function parameter available also `perSystem`.
    _module.args.modos = configG.modos;
    perSystem =
      {
        config,
        ...
      }:
      {
        _module.args.modos' = config.modos;
      };
  };
}
