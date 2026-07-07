{
  lib,
  makeWrapper,
  # Own arguments.
  modos,
  compName,
  buildType ? "release",
  environmentType ? "production",
  ...
}:
let
  target = "service";

  # Specify some more runtime dependencies if needed.
  # Dependencies should be added on the `service`.
  runtimeDeps = [ ];
in
modos.lib.build.buildRustPackage {
  inherit
    compName
    buildType
    environmentType
    target
    ;
  pname = compName;

  version = modos.lib.component.readVersion compName;

  src = modos.lib.fileset.toSource [ compName ];

  vendorHash = "sha256-ugPJchw4qIT05+aMFaQ2a4oO757/XSUkyUJwwnf6VQA=";

  doCheck = false;

  buildInputs = [ makeWrapper ];

  # Add runtime dependencies.
  postInstall = ''
    wrapProgram "$out/bin/${target}" \
      --prefix PATH : ${lib.makeBinPath runtimeDeps}
  '';

  meta = {
    description = compName;
    homepage = "https://github.com/sdsc-ordes/modos-rs";
    license = lib.licenses.apsl20;
    maintainers = [ "sdcs-ordes" ];
    mainProgram = compName;
  };
}
