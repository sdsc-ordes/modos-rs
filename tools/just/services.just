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
    nix run -L --show-trace "./tools/nix#test-services" -- --detached --detached-with-tui "$@"

# Stop the test services.
[group('services')]
stop *args:
    process-compose -u "{{pc_socket}}" down

# Attach to the running test services.
[group('services')]
attach *args:
    process-compose -u "{{pc_socket}}"  attach

# Delete all data from the services.
delete-data:
    just -f "{{source_file()}}" stop
    rm -rf "{{output_dir}}/process-compose"

[group('keycloak')]
keycloak-export:
    process-compose -u "{{pc_socket}}" process stop keycloak || true
    process-compose -u "{{pc_socket}}" process start keycloak-realm-export-all
