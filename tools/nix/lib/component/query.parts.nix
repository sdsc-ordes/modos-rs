{
  lib,
  config,
  ...
}:
let
  inherit (lib.attrsets) filterAttrs;
  inherit (builtins) readDir pathExists;

  inherit (config.flake.lib) component;
  inherit (config.flake.lib.common) yaml;
in
{
  flake.lib.component.query = {
    /*
      Gets all components in `path`.

      # Inputs
      `path` : Path to search for components.

      # Type
      ```
      getComponents :: { path, configName ? } -> { "<name>" = { name, path, pathRel, config }; ... }
      ```

      # Examples
      ```nix
      getComponents { path = ./components; }
      =>
      {
        # Keyed by `name` read from `.component.yaml`.
        service-a = {
          name = "service-a";
          path = ./components/service-a;         # Absolute component path.
          pathRel = "service-a";                 # Path relative to the search root.
          config = { name; version; language; }; # Parsed from `.component.yaml`.
        };
        ...
      }
      ```
    */
    getComponents =
      {
        path,
        configName ? component.configFileName,
      }:
      let
        # Get the component config in the dir.
        getCompConfig =
          dir:
          let
            config = dir + "/${configName}";
            exists = pathExists config;
          in
          {
            inherit config exists;
          };

        # Determine if output from `readDir` is a directory.
        isDir = basename: type: type == "directory";

        # List directory entries in the `dir`: [ entry, ... ]
        # where `entry := { path = "..."; basename = "..."}`
        getDirs =
          path:
          let
            res = filterAttrs isDir (readDir path);
          in
          lib.mapAttrsToList (basename: type: { path = lib.path.append path basename; }) res;

        # Visit a directory entry and determine if its a component.
        # otherwise recurse.
        visitDir =
          { path, ... }@entry:
          let
            res = getCompConfig path;
          in
          if res.exists then
            # Return the component.
            [ (entry // { inherit (res) config; }) ]
          else
            # Recurse to the directory.
            getComponents path;

        # Gets all components (directories with `.component.yaml` file)
        # recursively inside `dir`
        # = [ { basename = "..."; path = "...";}, ... ]
        getComponents =
          path:
          let
            dirs = getDirs path; # dirs = [ entry, ...]
          in
          # Concatenating will reduce all [ visitDir(entry) => [], ... ] to
          # one flat map.
          builtins.concatMap visitDir dirs;

        # Map that to proper components: { "comp-name" = {...}; ... }
        toComponentAttr =
          comps:
          builtins.listToAttrs (
            map (
              e:
              let
                conf = yaml.readSimple e.config [
                  "name"
                  "version"
                  "language"
                ];
              in
              {
                name = conf.name;
                # Update the component entry
                # with the read YAML file.
                value = e // {
                  name = conf.name;
                  config = conf;
                  pathRel = lib.path.removePrefix path e.path;
                };
              }
            ) comps
          );
      in
      toComponentAttr (getComponents path);
  };
}
