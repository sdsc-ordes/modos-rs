{
  lib,
  # Own arguments.
  modos,
  build,
  compName,
  buildType ? "release",
  environmentType ? "production",
  ...
}:
let
  inherit (modos.lib) component;
in
build.buildGoModule {
  inherit buildType environmentType;
  inherit compName;

  pname = compName;
  version = component.readVersion compName;

  src = modos.fileset.toSource [
    compName
    "quitsh" # just for tests.
  ];

  target = "service";
  vendorHash = "sha256-KTq/NAjeKJxoR7UTY3KxGGoSVcuUWhzJ/IedVcgOqjk=";

  doCheck = false;

  meta = {
    description = compName;
    homepage = "https://gitlab.com/data-custodian/dac-portal";
    license = lib.licenses.apsl20;
    maintainers = [ "sdcs-ordes" ];
    mainProgram = compName;
  };
}
