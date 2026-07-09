{
  self,
  ...
}:
{
  perSystem =
    {
      system,
      ...
    }:
    let
      pkgs = self.lib.nixpkgs.importPkgs { inherit system; };
      pkgsStable = self.lib.nixpkgs.importPkgsStable { inherit system; };
    in
    {
      # Define two arguments `pkgs` and `pkgsStable` available on all flake-parts modules.
      _module.args.pkgs = pkgs;
      _module.args.pkgsStable = pkgsStable;

      legacyPackages.unstable = pkgs;
      legacyPackages.stable = pkgsStable;
    };
}
