{
  lib,
  rustPlatform,
  libComponent,
  ...
}:
{
  # The component name we are building.
  compName,
  # The package name.
  pname,
  # The package version.
  version,
  # The source (component directory) of this Rust build.
  src,

  lockFileRel ? "Cargo.lock",

  # Meta information for `mkDerivation`.
  meta,

  # The output binary.
  target ? "service",

  # Build flags for `quitsh`.
  buildType ? "release",
  environmentType ? "production",
  ...
}@args: # NOTE: `args` doesn't capture default arguments.
let
  compDirRel = libComponent.getRootPathRel compName;
  # The name of the derivation.
  name = "${args.pname}-${args.version}";

  forwardArgs = lib.removeAttrs args [
    "vendorHash"
    "compName"
  ];

  compSrc = "${src}/${compDirRel}";
  lockFile = "${compSrc}/${lockFileRel}";
in
rustPlatform.buildRustPackage (
  forwardArgs
  // {
    inherit
      version
      pname
      name
      meta
      ;

    src = compSrc;

    cargoDeps = rustPlatform.importCargoLock {
      inherit lockFile;
    };

    postPatch = ''
      install -m 644 ${lockFile} Cargo.lock
    '';

    buildType = if buildType == "debug" then "debug" else "release";
    buildFeatures = [ environmentType ];

    cargoBuildFlags = [
      "--bin"
      target
    ];

    postInstall =
      (args.postInstall or "")
      +
      # bash
      ''
        execPath="$out/bin/${target}"
        if [ ! -f "$execPath" ]; then
          echo "No binary in '$(pwd)/$execPath'." >&2
          echo "You must define an output '${target}' in 'Cargo.toml'." >&2
          exit 1
        fi

        # Symlink the process name to the built target.
        if [ ! -f "$out/bin/${pname}" ]; then
          ln -sf "$execPath" "$out/bin/${pname}"
        fi
      '';
  }
)
