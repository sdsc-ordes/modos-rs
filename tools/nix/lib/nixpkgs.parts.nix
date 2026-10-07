{
  lib,
  inputs,
  ...
}:
let
  config = {
    allowUnfree = true;
  };

  mkMultiverse =
    { system }:
    inputs.multiverse.lib.mkMultiverse {
      inherit system config;
      # No overlays for now needed.
    };
in
{

  flake.lib.nixpkgs = {
    # Shorthand to access the std library in the `nix repl`.
    lib = inputs.nixpkgs.lib;

    inherit mkMultiverse;

    importPkgs =
      {
        system,
      }:
      let
        mvs = mkMultiverse { inherit system; };

        # Use a 7 days behind a pinned unstable for security reasons.
        # To update run: `nix flake update multiverse` and `direnv reload`
        # and check the assert below and take the unstable tip commit.
        pkgsUnstableCooldown = mvs.daysBehind "34ab99075ac4f7e40cf037eef32cb1c360bb85e9" 7;

        pkgsUnstable =
          assert lib.assertMsg (pkgsUnstableCooldown.multiverse.rev == inputs.nixpkgs.rev) ''
            Input 'nixpkgs' must be
            aligned with cooldown 7 days behind "${pkgsUnstableCooldown.multiverse.rev}".
            NixOS unstable tip: '${mvs.tip.multiverse.rev}'.
          '';
          lib.trace "NixOS unstable tip: '${mvs.tip.multiverse.rev}'." pkgsUnstableCooldown;
      in
      pkgsUnstable;

    importPkgsStable =
      {
        system,
      }:
      let
        mvs = mkMultiverse { inherit system; };
      in
      mvs.at "26.05"; # Branch: nixos-26.05
  };
}
