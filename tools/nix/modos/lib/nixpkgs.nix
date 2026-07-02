{
  inputs,
  ...
}:
let
  config = {
    allowUnfree = true;
  };

  # Imports nixpkgs from `nixpkgs` (input)
  # where the localSystem depends on some env. variables which require
  # `--no-pure-eval`, if not enabled, it will by default report the
  # normal host platform.
  #
  # Options:
  # - If `USE_LIBC_MUSL` is set it will use the `musl` version..
  importPkgs =
    {
      nixpkgs,
      system,
      overlays ? [ ],
    }:
    let
      lib = nixpkgs.lib;
      p = import nixpkgs (
        lib.info "modos: Importing nixpkgs for system = '${system}'" {
          inherit system overlays config;
        }
      );
      pkgs = if builtins.getEnv "USE_LIBC_MUSL" == "true" then p.pkgsMusl else p;
    in
    pkgs;
in
{
  pkgs =
    {
      system,
      overlays ? [ ],
    }:
    importPkgs {
      nixpkgs = inputs.nixpkgs;
      inherit system overlays;
    };

  pkgsStable =
    {
      system,
      overlays ? [ ],
    }:
    importPkgs {
      nixpkgs = inputs.nixpkgs-stable;
      inherit system overlays;
    };
}
