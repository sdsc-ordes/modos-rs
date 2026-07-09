{
  config,
  lib,
  ...
}:
let
  inherit (config.modos.lib) repo component;

  # The components root dir.
  compsDir = lib.path.append repo.rootDir "./components";

  # All components found in `rootDir`.
  # Type: [Component] (see `getComponents`).
  comps = component.query.getComponents { path = repo.rootDir; };
in
{
  modos.lib.component = {
    inherit comps compsDir;

    # Root paths of components.
    getRootPathRel = compName: (lib.getAttr compName component.comps).pathRel;
    getRootPath = compName: (lib.getAttr compName component.comps).path;

    # Get the component in `comps` determined by a subpath `path` inside
    # a components directory.
    getCompFromPath =
      path:
      let
        matchPrefix = entry: lib.path.hasPrefix entry.value.path path;
        compFound = lib.findFirst matchPrefix null (lib.attrsets.attrsToList comps);
      in
      assert lib.assertMsg (compFound != null) "Component from '${path}' not found.";
      compFound.value;

    /*
      Get the version of a component.

      # Inputs
      `compName`: The component name (as read from `.component.yaml`).

      # Type
      ```
      readVersion :: String -> String
      ```
      # Examples
      ```nix
      readVersion "service-a"
      => "3.0.1"
      ```
    */
    readVersion = compName: comps.${compName}.config.version;
  };
}
