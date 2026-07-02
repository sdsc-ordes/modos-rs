set positional-arguments
set shell := ["bash", "-cue"]

default:
    just --list --unsorted -f "{{source_file()}}"

# Get an access token on the online keycloak instance with device flow.
access-token-device-flow \
    client_id="dac-cli" \
    realm="fhdportal" \
    endpoint="https://login.biomedit.ch/realms" \
    scopes='openid profile roles dac:read dac:write policy:read policy:write submission:read submission:write membership:read membership:write dac:write dataset:read datarequest:read datarequest:write':
    #!/usr/bin/env bash

    if [ -t 2 ]; then
        TTY_OUT="/dev/fd/2"   # stderr is a terminal
    elif [ -e /dev/tty ]; then
        TTY_OUT="/dev/tty"   # fallback: controlling terminal
    fi

    endpoint="{{endpoint}}/{{realm}}/protocol/openid-connect"
    device_resp=$(curl --header 'Content-Type: application/x-www-form-urlencoded' \
        -d 'client_id={{client_id}}' \
        -d "scope=$scopes" \
        --url "$endpoint/auth/device" 2>"$TTY_OUT") || {
        echo "Could not get device flow response." &>"$TTY_OUT"
        exit 1
    }

    echo -e "{{GREEN}}Device flow response: \n$device_resp{{NORMAL}}" &>"$TTY_OUT"
    device_code=$(echo "$device_resp" | jq -r ".device_code")
    uri=$(echo "$device_resp" | jq -r ".verification_uri_complete")

    echo "{{BLUE}}{{BOLD}}Open the following url '${uri}' in the browser{{NORMAL}}" &>"$TTY_OUT"
    if [ "{{os()}}" = "macos" ]; then
        open "$uri" &>"$TTY_OUT" || {
            echo "Could not open url." &>"$TTY_OUT"
            exit 1
        }
    else
        xdg-open "$uri" &>"$TTY_OUT" || {
            echo "Could not open url." &>"$TTY_OUT"
            exit 1
        }
    fi

    echo "Press enter when you logged in in the browser." &>"$TTY_OUT"
    read -r test

    access_token=$(curl \
        --header 'Content-Type: application/x-www-form-urlencoded' \
        -d 'grant_type=urn:ietf:params:oauth:grant-type:device_code' \
        -d "device_code=$device_code" \
        -d 'client_id={{client_id}}' \
        --url "$endpoint/token" 2>"$TTY_OUT" | jq -r ".access_token") || {
            echo "Could not get access token" &>"$TTY_OUT"
            exit 1
        }

    echo -e "{{RED}}Access Token:\n$access_token{{NORMAL}}" &>"$TTY_OUT"
