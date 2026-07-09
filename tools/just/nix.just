set positional-arguments
set shell := ["bash", "-cue"]

default:
    just --list --unsorted -f "{{source_file()}}"

# Enter the Nix development shell `$1` and execute the command `${@:2}`.
[no-cd]
[group('nix')]
develop *args:
    #!/usr/bin/env bash
    set -eu
    root_dir=$(git rev-parse --show-toplevel) || exit 1
    flake_dir="$root_dir/tools/nix"

    shell="$1"; shift 1;
    nix_args=()
    args=("$@") && [ "${#args[@]}" != 0 ] ||
        args=(env SHELL="$SHELL" "$SHELL")

    nix develop \
        --no-pure-eval \
        --accept-flake-config \
        "$flake_dir#$shell" \
        --command "${args[@]}"

# Start an Nix interpreter with loaded flake.
[no-cd]
[group('nix')]
repl:
    #!/usr/bin/env bash
    set -eu

    root_dir=$(git rev-parse --show-toplevel) || exit 1
    flake_dir="$root_dir/tools/nix"

    nix repl \
        --no-pure-eval \
        --accept-flake-config \
        "$flake_dir"

# Build a package in the `packages` output of `flake.nix`
alias build := package
[no-cd]
[group('nix')]
package attrname *args:
    #!/usr/bin/env bash
    set -eu
    attrname="$1"
    nix_args=("${@:2}")
    root_dir=$(git rev-parse --show-toplevel) || exit 1
    flake_dir="$root_dir/tools/nix"
    out_dir="$root_dir/.output/package"

    mkdir -p "$out_dir"
    outlink="$out_dir/$attrname"

    nix build \
        --accept-flake-config \
        --no-pure-eval \
        -L "$flake_dir#$attrname" \
        "${nix_args[@]}" \
        --out-link "$outlink" --json

# Show all packages.
[no-cd]
[group('nix')]
show:
    #!/usr/bin/env bash
    set -eu
    root_dir=$(git rev-parse --show-toplevel) || exit 1
    flake_dir="$root_dir/tools/nix"

    system=$(nix eval --impure --raw --expr builtins.currentSystem)
    nix eval --json "$flake_dir#packages.$system" | jq

# Fix all Nix hashes in repository.
[no-cd]
[group('nix')]
fix:
    just quitsh nix fix-hash
