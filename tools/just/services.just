set positional-arguments
set shell := ["bash", "-cue"]
set dotenv-load := true
root_dir := `git rev-parse --show-toplevel`
output_dir := root_dir / ".output"
pc_socket := output_dir / "process-compose" / "pc.sock"

default:
    just --list --unsorted -f "{{source_file()}}"

# Run the test services.
[group('services')]
start *args:
    #!/usr/bin/env bash
    cd "{{root_dir}}"
    nix run -L --show-trace "./tools/nix#test-services" -- "$@"

# Stop the test services.
[group('services')]
stop *args:
    #!/usr/bin/env bash
    cd "{{root_dir}}"
    just -f "{{source_file()}}" pc down

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
    echo "Delete output."
    chmod -R +w "{{output_dir}}/process-compose"
    rm -rf "{{output_dir}}/process-compose"

[group('keycloak')]
keycloak-export: check-started
    #!/usr/bin/env bash
    cd "{{root_dir}}"
    if process-compose process list | grep keycloak; then
        process-compose process stop keycloak || true
    fi
    just -f "{{source_file()}}" pc process start keycloak-realm-export-all

[group('authentik')]
authentik-export: check-started
    #!/usr/bin/env bash
    cd "{{root_dir}}"
    just -f "{{source_file()}}" pc process start authentik-export

[group('process-compose')]
pc *args:
    #!/usr/bin/env bash
    cd "{{root_dir}}"
    process-compose -u "{{pc_socket}}" "$@"

[private]
check-started:
    #!/usr/bin/env bash
    [ -f "{{pc_socket}}" ] || {
        echo "Process compose is not started.";
        exit 1
    }
