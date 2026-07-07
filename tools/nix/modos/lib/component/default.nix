{
  lib,
  libCommon,
  rootDir,
}:
let
  configFileName = ".component.yaml";

  # All query functions.
  query = (import ./query.nix) { inherit lib libCommon configFileName; };

  # All components.
  compsDir = lib.path.append rootDir "./components";
  comps = query.getComponents { path = rootDir; };
in
{
  inherit configFileName;
  inherit query;

  # All components found in `rootDir`.
  # Type: [Component] (see `getComponents`).
  inherit comps;

  # The components root dir.
  inherit compsDir;

  # Root paths of components.
  getRootPathRel = compName: (lib.getAttr compName comps).pathRel;
  getRootPath = compName: (lib.getAttr compName comps).path;

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
    `compRoot`: The component source directory where `.component.yaml` resides.

    # Type
    ```
    readVersion :: Path -> String
    ```
    # Examples
    ```nix
    readVersion "./components/a"
    => "3.0.1"
    ```
  */
  readVersion = compName: comps.${compName}.config.version;

  /*
    Returns the build directory in the components directory.

    # Inputs
    `compRoot` : The component directory.

    # Type
    ```
    getBuildDir :: String -> String
    ```
    # Examples
    ```nix
    getBuildDir "./components/a"
    => ./components/a/.output/build
    ```
  */
  getBuildDir =
    compRoot:
    lib.concatStringsSep "/" [
      compRoot
      ".output"
      "build"
    ];
}
