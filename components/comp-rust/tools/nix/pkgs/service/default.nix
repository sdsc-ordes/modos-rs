{
  lib,
  makeWrapper,
  # Own arguments.
  modos,
  build,
  compName,
  buildType ? "release",
  environmentType ? "production",
  ...
}:
let
  inherit (modos.lib) component fileset;
  target = "service";

  # Specify some more runtime dependencies if needed.
  # Dependencies should be added on the `service`.
  runtimeDeps = [ ];
in
build.buildRustPackage {
  inherit
    compName
    buildType
    environmentType
    target
    ;
  pname = compName;

  version = component.readVersion compName;

  src = fileset.toSource { filesets = [ compName ]; };

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
