{
  lib,
  buildGo125Module,
  installShellFiles,
  testers,
  git,
  modosLib,
  self,
}:
let
  compName = "quitsh";
  libComponent = modosLib.component;
  libFileset = modosLib.fileset;

  version = libComponent.readVersion compName;
in
# NOTE: It can be that this derivation does not build anymore.
#       This is probably due to caching of the source directory.
#       Before debugging: Change the .component.yaml version and increment its patch number
#       to trigger another source hash.
buildGo125Module {
  inherit version;
  pname = compName;
  src = libFileset.toSource [
    compName
  ];

  modRoot = libFileset.getRootPathRel compName;

  # This hash defines the fixed-output derivation of the dependencies (FOD).
  # You can set the hash to "" and do:
  # ```
  # just nix::package quitsh
  # ```
  # to check if a new hash must be here:
  vendorHash = "sha256-HHM+XAGyJYgkG3YYmhFtG+eG6OvVX8LZPWBoPqM6Adw=";
  proxyVendor = true;

  nativeBuildInputs = [ installShellFiles ];
  nativeCheckInputs = [ git ];

  ldflags =
    let
      modulePath = "modos-rs/tools/quitsh";
    in
    [
      "-s"
      "-w"
      "-X ${modulePath}/pkg/build.buildVersion=${version}"
    ];

  postInstall = ''
    installShellCompletion --cmd quitsh \
      --bash <($out/bin/quitsh completion bash) \
      --fish <($out/bin/quitsh completion fish) \
      --zsh <($out/bin/quitsh completion zsh)
  '';

  passthru.tests.version = testers.testVersion {
    package = self;
    command = "quitsh --version";
    inherit version;
  };

  meta = {
    description = "Tool to build/test/lint/deploy components in a monorepo - quit using `sh`.";
    homepage = "https://github.com/sdsc-ordes/modos-rs";
    license = lib.licenses.apsl20;
    maintainers = [ "gabyx" ];
    mainProgram = "quitsh";
  };
}
