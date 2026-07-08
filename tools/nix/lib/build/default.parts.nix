{ ... }:
{
  flake.lib.build = {
    # Function to create a build support function for building a Go module.
    createBuildGoModule = import ./go/build-module.nix;
    createBuildRustPackage = import ./rust/build-package.nix;
  };
}
