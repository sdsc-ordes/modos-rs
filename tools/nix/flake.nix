{
  description = "modos-rs";

  nixConfig = {
    extra-trusted-substituters = [
      "https://nix-community.cachix.org"
      "https://devenv.cachix.org"
      "ssh://nix-ssh@nix-cache.swissmodos.ch"
      "https://modos-rs.cachix.org"
    ];
    extra-trusted-public-keys = [
      "modos-rs.cachix.org-1:qznktesR1I4KmWol3CwKfi0vM0BaH1+rSVzLYwfenG0="
      "nix-community.cachix.org-1:mB9FSh9qf2dCimDSUo8Zy7bkq5CX+/rkCWyvRCYg3Fs="
      "devenv.cachix.org-1:w1cLUi8dv3hnoSPGAuibQv+f9TZLr6cv/Hm9XgU50cw="
      "nix-cache.swissmodos.ch.1:rPQnp1nJav3UluO5MeomJTEPeqffeIu7Y41xpecBqMA="
    ];

    allow-import-from-derivation = "true";
  };

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";

    # Pinning some packages:
    rustfs.url = "github:rustfs/rustfs/5d737eaeb7fcab5d40c655ba60a494e93dd98922";
    # ref: commit points to nixpkgs-unstable
    process-compose.url = "github:nixos/nixpkgs/b5aa0fbd538984f6e3d201be0005b4463d8b09f8";

    # The devenv module to create good development shells.
    # The `nixpkgs-devenv` must aligned with the pinned version.
    devenv = {
      url = "github:cachix/devenv?ref=v2.1.2";
      inputs.nixpkgs.follows = "nixpkgs-devenv";
    };
    # This is the rolling nixpkgs with what devenv was tested.
    nixpkgs-devenv = {
      url = "github:cachix/devenv-nixpkgs?ref=ec3063523dcd911aeadb50faa589f237cdab5853";
    };

    # Format the repo with nix-treefmt.
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    # Quitsh functionality.
    quitsh = {
      url = "github:sdsc-ordes/quitsh?ref=main&dir=tools/nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    # The Rust overlay to include the latest toolchain.
    rust-overlay = {
      url = "github:oxalica/rust-overlay";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    # Importing flake-parts modules recursively.
    import-tree = {
      url = "github:vic/import-tree";
    };

    # Using `nix-systems` flake specification.
    systems = {
      url = "path:./flake/systems.nix";
      flake = false;
    };

    # Structuring the flake outputs with a NixOS modules.
    flake-parts = {
      url = "github:hercules-ci/flake-parts";
    };
  };
  # We use flake-parts to assemble all flake outputs.
  # This gives nicer modularity. All `.parts` files are
  # `flake-parts` module files.
  outputs =
    inputs:
    let
      lib = inputs.nixpkgs.lib;
    in
    inputs.flake-parts.lib.mkFlake { inherit inputs; } (
      lib.pipe inputs.import-tree [
        # NOTE: Uncomment the below to inspect what modules are loaded.
        (i: i.map (x: lib.info "modos: Importing: '${x}'" x))
        (i: i.filter (lib.hasInfix ".parts."))
        (
          i:
          i [
            ./.
            ../quitsh
            ../../components
          ]
        )
      ]
    );
}
