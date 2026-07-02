{
  lib,
  makeWrapper,
  cowsay,
  modosLib,
  compName,
  buildType ? "release",
  environmentType ? "production",
  ...
}:
let
  target = "service";

  # TODO: remove this.
  runtimeDeps = cowsay;
in
modosLib.build.buildRustPackage {
  inherit
    compName
    buildType
    environmentType
    target
    ;
  pname = compName;

  version = modosLib.component.readVersion compName;

  src = modosLib.fileset.toSource [ compName ];

  vendorHash = "sha256-ugPJchw4qIT05+aMFaQ2a4oO757/XSUkyUJwwnf6VQA=";

  doCheck = false;

  buildInputs = [ makeWrapper ];

  # Add runtime dependencies.
  postInstall = ''
    wrapProgram "$out/bin/${target}" \
      --prefix PATH : ${lib.makeBinPath [ runtimeDeps ]}
  '';

  meta = {
    description = compName;
    homepage = "https://github.com/sdsc-ordes/modos-rs";
    license = lib.licenses.apsl20;
    maintainers = [ "sdcs-ordes" ];
    mainProgram = compName;
  };
}
