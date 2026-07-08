{
  lib,
  rust-platform,
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
  # The SRI hash of the vendored dependencies.
  # If vendor hash is `nulL`, then no dependencies are fetched and
  # the build relies on the vendor folder within the source.
  vendorHash,

  # Meta information for `mkDerivation.
  meta,

  # The output binary.
  target ? "service",

  # Build flags for `quitsh`.
  buildType ? "release",
  environmentType ? "production",
  ...
}@args: # NOTE: `args` doesnt capture default arguments.
let
  compDirRel = libComponent.getRootPathRel compName;
  # The name of the derivation.
  name = "${args.pname}-${args.version}";

  forwardArgs = lib.removeAttrs args [
    "vendorHash"
    "compName"
  ];
in
rust-platform.buildRustPackage (
  forwardArgs
  // {
    inherit
      version
      pname
      name
      meta
      ;

    src = "${src}/${compDirRel}";

    cargoHash = vendorHash;
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
