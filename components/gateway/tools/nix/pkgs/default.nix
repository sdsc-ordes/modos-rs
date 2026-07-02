args:
let
  # debug/release builds -> debug with debug symbols for the debugger (also we misuse this, to
  # enable certain things)

  compName = "gateway";

  service = args.pkgs.callPackage ./service (args // { inherit compName; });
  service-dev = args.pkgs.callPackage ./service (
    args
    // {
      inherit compName;
      buildType = "debug";
      environmentType = "development";
    }
  );

  service-image = args.pkgs.callPackage ./service-image (args // { inherit compName service; });
in
{
  inherit service service-dev service-image;
}
