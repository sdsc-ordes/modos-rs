{
  lib,
  inputs,
  ...
}:
{
  perSystem =
    {
      system,
      self',
      config,
      inputs',
      pkgs,
      ...
    }:
    let
      cn = config.modos;

      codecov-cli = self'.packages.codecov-cli;
      # let
      #   pkgs = import inputs.nixpkgs-codecov {
      #     inherit system;
      #     config = {
      #       allowUnfree = true;
      #     };
      #   };
      # in
      # pkgs.codecov-cli;

      rust = lib.genAttrs [ "shell" "release" ] (
        name:
        let
          pkgsRust = pkgs.extend (import inputs.rust-overlay);
          toolchainFile = cn.lib.fileset.rootDir + "/tools/configs/rust/rust-toolchain-${name}.toml";
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
        inherit (pkgs) process-compose;

        inherit codecov-cli;
        inherit rust;

        minio = inputs'.nixpkgs-minio.legacyPackages.minio;
        garage = pkgs.garage_2;
        rustfs = inputs'.rustfs.packages.default;
      };
    in
    {
      modos = {
        build = {
          # All pinned packages.
          pinned = pkgsPinned;

          # Generate the Go build-support functions.
          buildGoModule = pkgs.callPackage cn.lib.build.createBuildGoModule {
            inherit (cn.pkgs) quitsh;
            inherit (pkgsPinned) go;
            libComponent = cn.lib.component;
          };

          buildRustPackage = pkgs.callPackage cn.lib.build.createBuildRustPackage {
            inherit (cn.pkgs) quitsh;
            rust-platform = pkgsPinned.rust.release.platform;
            libComponent = cn.lib.component;
          };
        };
      };
    };
}
