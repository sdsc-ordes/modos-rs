{
  lib,
  # Add. packages.
  buf,
  # Own arguments.
  modos,
  build,
  compName,
  buildType ? "release",
  environmentType ? "production",
  ...
}:
let
  inherit (modos.lib) fileset component;
in
build.buildGoModule {
  inherit buildType environmentType;
  inherit compName;

  pname = compName;
  version = component.readVersion compName;

  src = fileset.toSource {
    filesets = [
      compName
    ];
  };

  preConfigure = ''
    # Executing `buf` in `go generate` needs somehow `HOME` to be set.
    # Nix sandbox has not home.
    export HOME=$(mktemp -d)
  '';

  nativeBuildInputs = [
    buf
  ];

  target = "service";
  vendorHash = "sha256-69EGfYXwNJXdTZZvp0jHQSIPEF0RrDWZLNMnKHdglko=";

  doCheck = true;

  meta = {
    description = compName;
    homepage = "https://github.com/sdcs-ordes/modos-rs";
    license = lib.licenses.apsl20;
    maintainers = [ "sdcs-ordes" ];
    mainProgram = compName;
  };
}
