{
  lib,
  modosLib,
  compName,
  buildType ? "release",
  environmentType ? "production",
  ...
}:
modosLib.build.buildGoModule {
  inherit buildType environmentType;
  inherit compName;

  pname = compName;
  version = modosLib.component.readVersion compName;

  src = modosLib.fileset.toSource [
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
