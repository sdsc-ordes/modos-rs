{
  config,
  lib,
  ...
}:
let
  fs = lib.fileset;
  inherit (config.modos.lib) component;
  inherit (config.modos.lib) repo;

  # The root fileset.
  rootFileset = fs.gitTracked repo.rootDir;

  # All component filesets.
  # = { compName = FileSet; ...}
  compFileSets = lib.mapAttrs (compName: comp: fs.fromSource comp.path) component.comps;

  isFileset = a: (a._type or "") == "fileset";
in
{
  modos.lib.fileset = {
    getRootPathRel = component.getRootPathRel;
    getRootPath = component.getRootPath;

    # The filesets for each component.
    # Note: You cannot evaluate this directly with `nix repl`
    #       Use `toSource` and the inspect.
    components = compFileSets // {
      root = rootFileset;
    };

    /*
      Add the local files contained in the given filesets or fileset names (in `components.$filesetName`)
      to the store by using the `rootPath` as `root` in `toSource`.
      The `rootPath` in `toSource` represents the toplevel directory in the generated
      `/nix/store/...` path.

      See [fileset.toSource]().

      # Inputs

      `filesets`

      : The fileset names in `components` or pure `Fileset`s to put into the store.

      # Type

      ```
      toSource:: { filesets = [String or Fileset...], ... } -> SourceLike
      ```
      :::
    */
    toSource =
      {
        filesets,
        root ? repo.rootDir,
      }:
      let
        # Either accept a fileset name (look it up in `compFileSets` or it is a fileset itself.)
        sets = lib.map (arg: if (isFileset arg) then arg else compFileSets."${arg}") filesets;

        # Unionify all filesets together.
        sum = fs.unions sets;

        # Intersect with rootFileset to only include
        # what is in the root (no non-Git files).
        final = fs.intersection rootFileset sum;
      in
      fs.toSource {
        inherit root;
        fileset = final;
      };
  };
}
