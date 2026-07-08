{ lib, ... }:
{
  flake.lib.component = {
    # The file which defines a component.
    configFileName = ".component.yaml";

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
  };
}
