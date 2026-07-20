{
  lib,
  maven,
  # Our arguments.
  modos,
  compName,
  ...
}:
let
  inherit (modos.lib) component fileset;
  version = component.readVersion compName;
  pname = compName;
in
maven.buildMavenPackage {
  inherit pname version;

  src = fileset.toSource {
    filesets = [ compName ];
    root = component.comps.${compName}.path;
  };

  postPatch = ''
    substituteInPlace pom.xml \
      --replace-fail '<version>1.0.0</version>' '<version>${version}</version>'
  '';

  mvnHash = "sha256-7NyPWO0rskwdtZUplWPsxprdH4RN19hgNr1gIfkXYYA=";

  mvnParameters = "-DskipTests";

  installPhase = ''
    runHook preInstall
    install -Dm644 -t "$out/${pname}-${version}.jar" target/keycloak-mapper-${version}.jar
    runHook postInstall
  '';

  meta = with lib; {
    description = "Keycloak OIDC protocol mapper emitting per-bucket permissions from group attributes.";
    homepage = "https://github.com/sdsc-ordes/modos-rs";
    license = licenses.asl20;
    platforms = platforms.all;
    maintainers = [ "sdsc-ordes" ];
  };
}
