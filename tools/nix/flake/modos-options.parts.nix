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
    build = mkLibOption "Build functions.";
    common = mkLibOption "Common helper functions.";
    component = mkLibOption "Component helper functions.";
    nixpkgs = mkLibOption "Nixpkgs import utility.";
    service = mkLibOption "Service modules for 'service-flake'.";
    shell = mkLibOption "Nix dev. shell utilities.";
    toolchain = mkLibOption "Toolchain utilities for defining shells.";
  };

  libModosDefs = libDefs // {
    repo = mkLibOption "Repository definitions.";
    fileset = mkLibOption "Component scoped filesets utilities.";
  };

  packageSetType = lib.types.attrsOf (lib.types.either lib.types.package packageSetType);
in
{
  options = {
    # Define `lib` as flake output which is our
    # own standalone reusable library which is not
    # repository/system/pkgs dependent.
    flake.lib = libDefs;

    # Define the modos namespace with repository scoped functionality.
    # Include libs in `flake.lib`.
    modos.lib = libModosDefs;

    # Define the modos namespace as a flake output.
    flake.modos = lib.mkOption {
      type = lib.types.raw;
      readOnly = true;
    };

    # Defining a `perSystem` scoped module option `modos`.
    perSystem = mkPerSystemOption (
      {
        lib,
        ...
      }:
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
              description = "Flat packages set for all components (read-only).";
              readOnly = true;
            };
          };

          shells = lib.mkOption {
            type = lib.types.attrsOf lib.types.package;
            description = "All development Nix shells.";
            default = { };
          };
        };
      }
    );
  };
}
