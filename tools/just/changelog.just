set positional-arguments
set shell := ["bash", "-cue"]

default:
    just --list --unsorted -f "{{source_file()}}"

# Generate the changelog for unreleased commits.
# Only the first parent reachable traversal.
# See: https://github.com/git-chglog/git-chglog/issues/284
[group('aux')]
gen *args:
    #!/usr/bin/env bash
    set -eu

    start="${1:-HEAD}"
    end="${2:-}"
    config="${3:-tools/configs/git-cliff/config.toml}"
    file="${4:-CHANGELOG.md}"

    root_dir=$(git rev-parse --show-toplevel) || exit 1
    cd "$root_dir"

    export GIT_CLIFF_TAG_PATTERN="${GIT_CLIFF_TAG_PATTERN:-v\d+}"

    if [ -z "$end" ]; then
        end=$(git describe --abbrev=0 --tags)
        if [ "$end" = "" ]; then
            echo "No last tag found!" >&2
            echo "Create one!" >&2

            exit 1
        fi
    fi

    echo "Changelog in '$end..$start'" >&2

    if ! grep -q '<!-- next-content -->' "$file"; then
        echo "no '<!-- next-content -->' tag in '$file'"
        exit 1
    fi

    non_first_parent_commits=()
    readarray -t non_first_parent_commits < <(
        comm -23 <(git rev-list "$end..$start" | sort) \
                <(git rev-list --ancestry-path --first-parent "$end..$start" |sort)
    )

    skip_args=()
    if [ "${#non_first_parent_commits[@]}" -ne 0 ]; then
        # shellcheck disable=SC2086,SC2206
        skip_args=(--skip-commit)
        skip_args+=("${non_first_parent_commits[@]}")
    fi

    out=$(git-cliff --config "$config" \
        "$end..$start" \
        --strip header \
        "${skip_args[@]}" \
        )

    # Replace in file.
    echo "$out" | sed -i '/<!-- next-content -->/{
        r /dev/stdin
        d
    }' "$file"
