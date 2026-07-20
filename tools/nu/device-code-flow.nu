#!/usr/bin/env nu
# Tests the OAuth 2.0 Device Authorization Grant (a.k.a. "device code flow")
# against a local Keycloak *or* Authentik instance.
#
#   Keycloak : http://localhost:8081  realm `modos`, client `modos-cli`
#   Authentik: http://localhost:9001  client `modos-cli` (no realm concept)
#
# Usage:
#   nu device-code-flow.nu                        # keycloak (default)
#   nu device-code-flow.nu --provider authentik
#   nu device-code-flow.nu --provider authentik --scope "openid permissions"
#   nu device-code-flow.nu --host http://localhost:8081 --realm modos --client modos-cli
#
# Notes:
#   * `--realm` only applies to Keycloak.
#   * PKCE is always used for Keycloak (the `modos-cli` client requires S256).
#     For Authentik it is off by default; pass `--pkce` to force it.
#   * To obtain the `bucket-permissions` claim from Authentik, request it
#     explicitly: `--scope "openid permissions"` (Authentik scopes are opt-in).

def main [
    --provider: string = "keycloak" # identity provider: keycloak | authentik
    --host: string = "" # base URL; defaults per provider if empty
    --realm: string = "modos" # Keycloak realm (ignored for authentik)
    --client: string = "modos-cli"
    --scope: string = "permissions"
    --pkce # force PKCE (already implied for keycloak)
] {
    let host = if ($host | is-empty) {
        match $provider {
            keycloak => "http://localhost:8081"
            authentik => "http://localhost:9001"
            _ => {
                print $"(ansi red)Unknown provider '($provider)'. Use 'keycloak' or 'authentik'.(ansi reset)"
                exit 1
            }
        }
    } else {
        $host
    }

    let ep = match $provider {
        "keycloak" => ({
            device: $"($host)/realms/($realm)/protocol/openid-connect/auth/device"
            token: $"($host)/realms/($realm)/protocol/openid-connect/token"
        })
        "authentik" => ({
            device: $"($host)/application/o/device/",
            token: $"($host)/application/o/token/"
        })
        _ => {
            print $"(ansi red)Unknown provider '($provider)'. Use 'keycloak' or 'authentik'.(ansi reset)"
            exit 1
        }
    }
    let device_endpoint = $ep.device
    let token_endpoint = $ep.token

    # Keycloak's `modos-cli` mandates PKCE; for authentik it is optional.
    let use_pkce = ($provider == "keycloak") or $pkce

    print $"(ansi cyan)Provider     :(ansi reset) ($provider)"
    print $"(ansi cyan)Host         :(ansi reset) ($host)"
    if $provider == "keycloak" {
        print $"(ansi cyan)Realm        :(ansi reset) ($realm)"
    }
    print $"(ansi cyan)Client       :(ansi reset) ($client)"
    print $"(ansi cyan)Scopes       :(ansi reset) ($scope)"
    print $"(ansi cyan)PKCE         :(ansi reset) ($use_pkce)"
    print ""

    # --- Step 1: request a device code + user code --------------------------
    print $"(ansi cyan)→ POST ($device_endpoint)(ansi reset)"

    # PKCE: generate a random verifier and derive its SHA-256, base64url-encoded
    # challenge. The verifier itself is presented later at the token-polling step.
    let verifier = $"(random uuid)(random uuid)"
    let challenge = (
        $verifier
        | hash sha256 --binary
        | encode base64
        | str replace --all "+" "-"
        | str replace --all "/" "_"
        | str replace --all "=" ""
    )

    mut init_body = {client_id: $client, scope: $scope}
    if $use_pkce {
        $init_body = ($init_body | merge {
            code_challenge: $challenge
            code_challenge_method: "S256"
        })
    }

    let init = (
        http post $device_endpoint --full --allow-errors --content-type "application/x-www-form-urlencoded" $init_body
    )

    if $init.status != 200 {
        print $"(ansi red)Device authorization request failed \(HTTP ($init.status)\):(ansi reset)"
        print ($init.body | to json)
        print ""
        print $"(ansi yellow)Hint: ensure the client is public and the Device Authorization Grant is enabled.(ansi reset)"
        if $provider == "authentik" {
            print $"(ansi yellow)For authentik, the default brand must have `flow_device_code` set \(see modos-blueprint.yaml\).(ansi reset)"
        }
        exit 1
    }

    let dev = $init.body
    let interval = $dev.interval? | default 5
    let expires_in = $dev.expires_in? | default 600
    let verify_complete = $dev.verification_uri_complete? | default $dev.verification_uri

    # --- Step 2: show the user what to do -----------------------------------
    print ""
    print $"(ansi green_bold)Open this URL in your browser and log in:(ansi reset)"
    print $"  (ansi green)($verify_complete)(ansi reset)"
    print ""
    print $"Or visit ($dev.verification_uri) and enter code: (ansi yellow_bold)($dev.user_code)(ansi reset)"
    print ""

    # --- Step 3: poll the token endpoint ------------------------------------
    mut poll_body = {grant_type: "urn:ietf:params:oauth:grant-type:device_code", device_code: $dev.device_code, client_id: $client}
    if $use_pkce {
        $poll_body = ($poll_body | merge { code_verifier: $verifier })
    }

    mut wait = $interval
    let start = (date now)
    print $"(ansi cyan)Polling for token every ($wait)s \(code expires in ($expires_in)s\)...(ansi reset)"

    loop {
        if ((date now) - $start) > ($expires_in * 1sec) {
            print $"(ansi red)✗ Device code expired before authorization completed.(ansi reset)"
            exit 1
        }

        let resp = (
            http post $token_endpoint
              --full --allow-errors
              --content-type "application/x-www-form-urlencoded"
              $poll_body
        )

        if $resp.status == 200 {
            let tok = $resp.body
            print ""
            print $"(ansi green_bold)✓ Authorization successful!(ansi reset)"
            print ""
            print $"  token_type    : ($tok.token_type)"
            print $"  expires_in    : ($tok.expires_in)"
            print $"  refresh_token : (if 'refresh_token' in $tok { 'yes' } else { 'no' })"
            print $"  access_token  : (($tok.access_token) | str substring 0..48)..."
            print ""

            # Best-effort decode of the JWT access-token payload (nushell
            # version dependent — skipped silently if unsupported).
            try {
                let claims = (
                    $tok.access_token
                    | split row "."
                    | get 1
                    | decode base64 --url --nopad
                    | decode
                    | from json
                )
                print $"(ansi cyan)Access-token claims:(ansi reset)"
                print ($claims | table --expand)
            } catch {
                print $"(ansi dark_gray)\(JWT payload not decoded on this nushell version\)(ansi reset)"
            }

            exit 0
        }

        let err = $resp.body.error? | default "unknown_error"
        match $err {
            "authorization_pending" => {
                print $"  … waiting for user to authorize \(HTTP ($resp.status)\)"
            }
            "slow_down" => {
                $wait = $wait + 5
                print $"  … slow_down — increasing poll interval to ($wait)s"
            }
            _ => {
                print $"(ansi red)✗ Token request failed: ($err)(ansi reset)"
                print ($resp.body | to json)
                exit 1
            }
        }

        sleep ($wait * 1sec)
    }
}
