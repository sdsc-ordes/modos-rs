{
  lib,
  inputs,
  ...
}:
{
  perSystem =
    {
      self',
      config,
      inputs',
      pkgs,
      ...
    }:
    let
      modos = config.modos;

      process-compose = inputs'.process-compose.legacyPackages.process-compose;
      codecov-cli = self'.packages.codecov-cli;

      rust = lib.genAttrs [ "shell" "release" ] (
        name:
        let
          pkgsRust = pkgs.extend (import inputs.rust-overlay);
          toolchainFile = modos.lib.fileset.rootDir + "/tools/configs/rust/rust-toolchain-${name}.toml";
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

        garage = pkgs.garage_2;
        rustfs = inputs'.rustfs.packages.default;
      };
    in
    {
      modos.lib.build = {
        pinned = pkgsPinned;

        # Generate the Go build-support functions.
        buildGoModule = pkgs.callPackage modos.lib.build.createBuildGoModule {
          inherit (modos.packages) quitsh;
          inherit (pkgsPinned) go;
          libComponent = modos.lib.component;
        };

        # Generate the Rust build-support functions.
        buildRustPackage = pkgs.callPackage modos.lib.build.createBuildRustPackage {
          rust-platform = pkgsPinned.rust.release.platform;
          libComponent = modos.lib.component;
        };
      };
    };
}
