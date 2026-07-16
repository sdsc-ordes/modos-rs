set positional-arguments
set shell := ["bash", "-cue"]
set dotenv-load := true
root_dir := `git rev-parse --show-toplevel`
output_dir := root_dir / ".output"

default:
    just --list --unsorted -f "{{source_file()}}"

# Test the device code flow.
device-flow *args:
    nu "{{root_dir}}/tools/nu/device-code-flow.nu" "$@"
