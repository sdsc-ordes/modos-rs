set positional-arguments
set shell := ["bash", "-cue"]
set dotenv-load := true
root_dir := `git rev-parse --show-toplevel`
flake_dir := root_dir / "tools/nix"
output_dir := root_dir / ".output"
build_dir := output_dir / "build"
pc_socket := output_dir / "process-compose" / "pc.sock"

mod nix "./tools/just/nix.just"
mod ci "./tools/just/ci.just"
mod changelog "./tools/just/changelog.just"
mod jwt "./tools/just/jwt.just"

# Default target if you do not specify a target.
default:
    just --list

# Enter the default Nix development shell and execute the command `"$@`.
alias dev := develop
[no-cd]
[group('general')]
develop *args:
    just nix::develop "default" "$@"

# Run quitsh (by compiling it directly and executing it from the root).
alias q := quitsh
[group('general')]
quitsh *args:
   quitsh-direct "$@"

# Run quitsh (by compiling it directly and executing it from the current directory).
[no-cd]
[group('general')]
quitsh-nocd *args:
    quitsh-direct "$@"

# Clean cleans the components output folders.
[group('general')]
clean-comps comppattern *args:
    just quitsh clean --components "{{comppattern}}" "$@"

# Cleans the whole repository and all untracked files (careful !)
[group('general')]
clean *args:
    #!/usr/bin/env bash
    set -eu
    if [ -d ".devenv/state/go" ]; then
        chmod -R +w .devenv/state/go
    fi
    git clean -dfX
    direnv reload

# Format the whole repository.
[group('general')]
format *args:
    just quitsh format "$@"

# Build the project.
[group('general')]
build *args: setup
    just quitsh build \
        --components "pdf-rendrer" \
        "$@"

# Lint the project.
[group('general')]
lint: setup
    just quitsh lint --components "*" --parallel --fix
    just quitsh nix fix-hash

# Test the project.
[group('general')]
test: setup
    just quitsh test --components "*" --parallel --fix

# Run the test services.
[group('services')]
services-start *args:
    nix run -L --show-trace "./tools/nix#test-services" -- --detached --detached-with-tui "$@"

# Stop the test services.
[group('services')]
services-stop *args:
    process-compose -u "{{pc_socket}}" down

# Attach to the running test services.
[group('services')]
services-attach *args:
    process-compose -u "{{pc_socket}}"  attach

[group('keycloak')]
export-realm:
    process-compose -u "{{pc_socket}}" process stop keycloak || true
    process-compose -u "{{pc_socket}}" process start keycloak-realm-export-all

# Update dependencies.
[group('aux')]
update-deps quitshBranch="main":
    #!/usr/bin/env bash
    set -eu

    # Update to quitsh flake.
    (cd "{{flake_dir}}" &&
        nix flake update quitsh &&
        git add .)

    # Update components.
    (cd ./tools/quitsh && just update-deps "main" "{{quitshBranch}}")

    # Update go.mod tidy in all modules.
    readarray -t gomods < <(find "{{root_dir}}/components" -name "go.mod" -and -not -ipath "*.devenv*")
    for gomod in "${gomods[@]}"; do
        echo "Update go mod in '$gomod'."
        (cd "$(dirname "$gomod")" && go mod tidy && git add go.mod go.sum)
    done

    just quitsh nix fix-hash

# Setup development files (default done in `.envrc`).
[private]
[group('aux')]
setup *args:
    just quitsh setup "$@"
