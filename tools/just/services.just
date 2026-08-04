set positional-arguments
set shell := ["bash", "-cue"]
set dotenv-load := true
root_dir := `git rev-parse --show-toplevel`
output_dir := root_dir / ".output" / "run" / "process-compose"
pc_socket := output_dir / "pc-test-services.sock"

default:
    just --list --unsorted -f "{{source_file()}}"

# Run the test services without attaching to the TUI.
[group('services')]
start-notui *args:
    #!/usr/bin/env bash
    cd "{{root_dir}}"
    just quitsh process-compose start \
        --flake-dir "./tools/nix" "test-services" \
        --socket-path-file "{{pc_socket}}" "$@"

# Run the test services with attaching to the TUI.
[group('services')]
start *args:
    just -f "{{source_file()}}" start-notui -a "$@"

# Monitor services. For debugging.
[group('services')]
monitor *args: check-started
    just -f "{{source_file()}}" pc process monitor

# Stop the test services.
[group('services')]
stop *args:
    #!/usr/bin/env bash
    cd "{{root_dir}}"
    just -f "{{source_file()}}" check-started &&
        just -f "{{source_file()}}" pc down || true

# Push the test services to cachix. Run `cachix authtoken` first.
[group('cachix')]
cache:
    #!/usr/bin/env bash
    cd "{{root_dir}}"
    nix build --show-trace -L --print-out-paths "./tools/nix#test-services" |
        cachix push "$CACHIX_CACHE_NAME"

# Attach to the running test services.
[group('services')]
attach *args: check-started
    #!/usr/bin/env bash
    cd "{{root_dir}}"
    just -f "{{source_file()}}" pc attach

# Delete all data from the services.
delete-data:
    #!/usr/bin/env bash
    cd "{{root_dir}}"
    just -f "{{source_file()}}" stop || true

    if [ -d "{{output_dir}}" ]; then
        echo "Delete output."
        chmod -R +w "{{output_dir}}"
        rm -rf "{{output_dir}}"
    fi

[group('keycloak')]
keycloak-export: check-started
    #!/usr/bin/env bash
    cd "{{root_dir}}"
    if just -f "{{source_file()}}" pc process list | grep keycloak; then
        just -f "{{source_file()}}" pc process stop keycloak || true
    fi
    just -f "{{source_file()}}" pc process start keycloak-realm-export-all

[group('authentik')]
authentik-export: check-started
    #!/usr/bin/env bash
    cd "{{root_dir}}"
    just -f "{{source_file()}}" pc process start authentik-export

[group('oauth')]
device-code-flow *args:
    #!/usr/bin/env bash
    cd "{{root_dir}}"
    just quitsh oauth device-code-flow "$@"

[group('process-compose')]
pc *args:
    #!/usr/bin/env bash
    cd "{{root_dir}}"
    process-compose -u "$(cat {{pc_socket}})" "$@"

[private]
check-started:
    #!/usr/bin/env bash
    [ -f "{{pc_socket}}" ] || {
        echo "Process compose is not started.";
        exit 1
    }
