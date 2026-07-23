{
  lib,
  inputs,
  modos,
  ...
}:
{
  perSystem =
    {
      self',
      inputs',
      modos',
      pkgs,
      ...
    }:
    let
      inherit (modos.lib)
        repo
        component
        build
        ;

      process-compose = inputs'.process-compose.legacyPackages.process-compose;
      codecov-cli = self'.packages.codecov-cli;

      rust = lib.genAttrs [ "shell" "release" ] (
        name:
        let
          pkgsRust = pkgs.extend (import inputs.rust-overlay);
          toolchainFile = repo.rootDir + "/tools/configs/rust/rust-toolchain-${name}.toml";
          toolchain = pkgsRust.pkgsBuildHost.rust-bin.fromRustupToolchainFile toolchainFile;
          platform = pkgs.makeRustPlatform {
            cargo = toolchain;
            rustc = toolchain;
          };
        in
        {
          inherit toolchain platform;
        }
      );

      pkgsPinned = {
        go = pkgs.go_1_26;
        inherit process-compose;

        inherit codecov-cli;
        inherit rust;
      };
    in
    {
      modos.packages.pinned = pkgsPinned;

      modos.build = {
        # Generate the Go build-support functions.
        buildGoModule = pkgs.callPackage build.createBuildGoModule {
          inherit (modos'.packages.global) quitsh;
          inherit (pkgsPinned) go;
          libComponent = component;
        };

        # Generate the Rust build-support functions.
        buildRustPackage = pkgs.callPackage build.createBuildRustPackage {
          rust-platform = pkgsPinned.rust.release.platform;
          libComponent = component;
        };
      };
    };
}
