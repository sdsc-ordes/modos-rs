{ lib, ... }:
let

  # Create `devenv` modules where `pkgs` and `pkgsStable` are the normal packages from `nixpkgs`.
  createDevenvModules =
    {
      pkgs,
      pkgsStable,
      modos,
    }:
    let
      # The pinned packages.
      pkgsPinned = modos.lib.build.pinned;

      quitsh-direct-drv = pkgs.writeShellApplication {
        name = "quitsh-direct";
        text = ''
          #!/usr/bin/env sh
          root=$(git rev-parse --show-toplevel) || {
            echo "Could not determine repo. root dir." >&2
          }
          just -f "$root/tools/quitsh/justfile" run "$@"
        '';
        runtimeInputs = [
          pkgs.bash
          pkgs.git
          pkgs.just
          pkgsPinned.go
        ];
      };

      quitsh-direct-ci-drv = pkgs.writeShellScriptBin "quitsh-direct" ''
        exec ${lib.getExe modos.packages.quitsh} "$@"
      '';

      quitsh-direct = [
        {
          packages = [
            quitsh-direct-drv
          ];
        }
      ];

      quitsh-setup = [
        (
          { config, ... }:
          {
            enterShell = ''
              quitsh=$(command -v "quitsh-direct" 2>/dev/null || echo "quitsh")
              f="$DEVENV_DOTFILE/state/quitsh/setup-done"
              if [ ! -f "$f" ]; then
                mkdir -p $(dirname "$f") && touch "$f"
                "$quitsh" setup
              else
                ${config.quitsh.log.package}/bin/log info "Setup 'quitsh setup' already performed. ✅"
              fi
              unset quitsh
            '';
          }
        )
      ];

      addSetup = modules: modules ++ quitsh-setup;

      build-go = [
        {
          quitsh.toolchains = [ "build-go" ];
          quitsh.languages.go = {
            enable = true;
            # To make CGO and the debugger delve work.
            # https://nixos.wiki/wiki/Go#Using_cgo_on_NixOS
            enableHardeningWorkaround = true;
            package = pkgsPinned.go;
          };

          packages = [
            pkgs.git
          ];
        }
      ];

      build-rust = [
        {
          quitsh.toolchains = [ "build-rust" ];
          packages = [
            pkgs.cargo-watch
            # Coverage
            pkgs.cargo-llvm-cov
          ];

          languages.rust = {
            enable = true;
            toolchainPackage = pkgsPinned.rust.shell.toolchain;
          };

          languages.c.debugger = pkgs.lldb_18;
        }
      ];

      # Go development.
      dev-go = [
        {
          packages = [
            pkgs.golangci-lint
            pkgs.golangci-lint-langserver
            pkgs.typos-lsp
          ];
        }
      ];

      lint-go = [
        {
          quitsh.toolchains = [ "lint-go" ];

          packages = [
            pkgs.git
            pkgs.golangci-lint
            pkgsPinned.go
          ];
        }
      ];

      image-containerfile = [
        {
          quitsh.toolchains = [ "image-containerfile" ];
          packages = [
            pkgs.git
            pkgs.buildah
            pkgs.skopeo
          ];
        }
      ];

      image-nix = [
        {
          quitsh.toolchains = [ "image-nix" ];

          packages = with pkgs; [
            git
            skopeo
          ];
        }
      ];

      manifest-ytt = [
        {
          quitsh.toolchains = [ "manifest-ytt" ];

          packages = [
            pkgs.ytt
            pkgs.imgpkg
            pkgs.kbld
            pkgs.kubernetes-helm
            pkgs.kubeseal
            pkgs.vendir
            pkgs.sops
          ];
        }
      ];

      coverage-upload = [
        {
          quitsh.toolchains = [ "coverage-upload" ];

          packages = [
            pkgsPinned.codecov-cli
          ];
        }
      ];

      lint-trivy = [
        {
          quitsh.toolchains = [ "lint-trivy" ];
          packages = [
            pkgs.trivy
          ];
        }
      ];

      default =
        ci
        ++ build-go
        ++ dev-go
        ++ manifest-ytt
        ++ quitsh-direct
        ++ [
          (
            { lib, ... }:
            {
              quitsh.toolchains = [ "general" ];

              quitsh.config = lib.mkForce "tools/configs/quitsh/config.yaml";
              quitsh.configUser = "tools/configs/quitsh/config.user.yaml";

              dotenv.enable = true;

              packages = [
                pkgs.cachix

                # Essentials.
                pkgs.git
                pkgs.just
                pkgs.fd

                # Manifests
                # added by manifest-ytt module.

                # Web-Traffic
                pkgs.xh # WARNING: Use this instead of httpie adds PYTHONPATH
                pkgs.jwt-cli # Decode jwt tokens.

                # Inspect/upload images.
                pkgs.dive
                pkgs.skopeo

                # Process manager.
                pkgsPinned.process-compose

                # Changelog
                pkgs.git-cliff
              ];
            }
          )
        ];

      ci = [
        {
          quitsh.toolchains = [
            "ci"
            "git"
          ];
          quitsh.config = "tools/configs/quitsh/config-ci.yaml";

          packages = [
            quitsh-direct-ci-drv
            modos.packages.bootstrap
            modos.packages.quitsh
            pkgs.podman

            pkgs.openssh # SSH agent
          ];
        }
      ];
    in
    {
      # Main shells:
      default = addSetup default;
      ci = addSetup ci;

      # Toolchains:
      inherit
        # General CI ---------
        build-rust
        build-go
        lint-go
        lint-trivy

        image-nix
        image-containerfile

        manifest-ytt
        # --------------------

        # Auxiliary Tooling --
        coverage-upload
        # --------------------

        # General Development --
        dev-go
        # ----------------------
        ;

      # Quitsh:
      inherit
        quitsh-direct
        quitsh-setup
        ;
    };
in
{
  flake.lib.toolchain = {
    inherit createDevenvModules;
  };
}
