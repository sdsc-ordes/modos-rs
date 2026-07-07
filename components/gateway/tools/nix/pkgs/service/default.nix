{
  lib,
  # Own arguments.
  modos,
  compName,
  buildType ? "release",
  environmentType ? "production",
  ...
}:
modos.lib.build.buildGoModule {
  inherit buildType environmentType;
  inherit compName;

  pname = compName;
  version = modos.lib.component.readVersion compName;

  src = modos.lib.fileset.toSource [
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
