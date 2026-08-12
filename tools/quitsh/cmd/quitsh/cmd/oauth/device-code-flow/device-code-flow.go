// Package devicecodeflow implements small OAuth helper commands.
//
// The `device-code-flow` command tests the OAuth 2.0 Device Authorization
// Grant (a.k.a. "device code flow") against a local Keycloak *or* Authentik
// instance.
package devicecodeflow

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sdsc-ordes/quitsh/pkg/cli"
	"github.com/sdsc-ordes/quitsh/pkg/errors"
	"github.com/sdsc-ordes/quitsh/pkg/log"

	"github.com/spf13/cobra"
)

const (
	providerKeycloak  = "keycloak"
	providerAuthentik = "authentik"

	defaultPollInterval = 1   // seconds, per RFC 8628 if the server omits it.
	defaultExpiresIn    = 600 // seconds, fallback lifetime for the device code.
	pkceVerifierBytes   = 32  // yields a 43-char base64url verifier (spec minimum).
	jwtParts            = 2   // header.payload(.signature) — we only need payload.

	slowDownIncrement = 5 * time.Second // per RFC 8628 on `slow_down`.
)

// Settings holds the flags for the device-code-flow command.
type Settings struct {
	Provider string // identity provider: keycloak | authentik
	Host     string // base URL; defaults per provider if empty
	Realm    string // Keycloak realm (ignored for authentik)
	Client   string
	Scope    string
	PKCE     bool // force PKCE (already implied for keycloak)
}

// deviceResponse is the JSON returned by the device authorization endpoint.
//
//nolint:tagliatelle // OAuth 2.0
type deviceResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	Interval                int    `json:"interval"`
	ExpiresIn               int    `json:"expires_in"`
}

// tokenResponse is the JSON returned by the token endpoint on success.
//
//nolint:tagliatelle // OAuth 2.0
type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// tokenError is the JSON returned by the token endpoint on error.
type tokenError struct {
	Error string `json:"error"`
}

// AddCmd registers the `device-code-flow` command onto the given parent.
//
// Note: the context is taken from `cl.Ctx()` (the CLI's signal context) rather
// than `cmd.Context()`. quitsh runs cobra via plain `Execute()`, so
// `cmd.Context()` is an uncancellable `context.Background()` and would not abort
// on CTRL+C. `cl.Ctx()` is the `signal.NotifyContext` set up by
// `WithSignalContext`.
func AddCmd(cl cli.ICLI, parent *cobra.Command) {
	sett := Settings{
		Provider: providerKeycloak,
		Host:     "",
		Realm:    "modos",
		Client:   "modos-cli",
		Scope:    "openid permissions",
		PKCE:     false,
	}

	cmd := &cobra.Command{
		Use:   "device-code-flow",
		Short: "Test the OAuth 2.0 Device Authorization Grant against Keycloak or Authentik.",
		Long: "Test the OAuth 2.0 Device Authorization Grant (device code flow) " +
			"against a local Keycloak or Authentik instance.",
		RunE: func(_ *cobra.Command, _ []string) error {
			return run(cl.Ctx(), &sett)
		},
	}

	cmd.Flags().StringVar(&sett.Provider, "provider", sett.Provider,
		"Identity provider: keycloak | authentik.")
	cmd.Flags().StringVar(&sett.Host, "host", sett.Host,
		"Base URL; defaults per provider if empty.")
	cmd.Flags().StringVar(&sett.Realm, "realm", sett.Realm,
		"Keycloak realm (ignored for authentik).")
	cmd.Flags().StringVar(&sett.Client, "client", sett.Client,
		"OAuth client id.")
	cmd.Flags().StringVar(&sett.Scope, "scope", sett.Scope,
		"Space separated scopes to request.")
	cmd.Flags().BoolVar(&sett.PKCE, "pkce", sett.PKCE,
		"Force PKCE (already implied for keycloak).")

	parent.AddCommand(cmd)
}

// endpoints resolves the device and token endpoints for the configured provider.
func endpoints(sett *Settings) (host, device, token string, pkceEnabled bool, err error) {
	host = sett.Host

	switch sett.Provider {
	case providerKeycloak:
		if host == "" {
			host = "http://localhost:8081"
		}
		device = host + "/realms/" + sett.Realm + "/protocol/openid-connect/auth/device"
		token = host + "/realms/" + sett.Realm + "/protocol/openid-connect/token"
		pkceEnabled = false
	case providerAuthentik:
		if host == "" {
			host = "http://localhost:9001"
		}
		device = host + "/application/o/device/"
		token = host + "/application/o/token/"
		pkceEnabled = false
	default:
		err = errors.New(
			"unknown provider '%s'; use 'keycloak' or 'authentik'", sett.Provider)

		return
	}

	return host, device, token, pkceEnabled, nil
}

// run executes the full device code flow.
func run(ctx context.Context, sett *Settings) error {
	host, deviceEndpoint, tokenEndpoint, usePKCE, err := endpoints(sett)
	if err != nil {
		return err
	}

	log.Infof("Provider: %s", sett.Provider)
	log.Infof("Host: %s", host)
	if sett.Provider == providerKeycloak {
		log.Infof("Realm: %s", sett.Realm)
	}
	log.Infof("Client: %s", sett.Client)
	log.Infof("Scopes: %s", sett.Scope)
	log.Infof("PKCE: %v", usePKCE)

	var pkceChallenge, pkceVerifier string
	if usePKCE {
		pkceVerifier, pkceChallenge, err = newPKCE()
		if err != nil {
			return err
		}
	}

	dev, err := requestDeviceCode(ctx, sett, deviceEndpoint, pkceChallenge)
	if err != nil {
		return err
	}

	interval := dev.Interval
	if interval == 0 {
		interval = defaultPollInterval
	}
	expiresIn := dev.ExpiresIn
	if expiresIn == 0 {
		expiresIn = defaultExpiresIn
	}
	verifyComplete := dev.VerificationURIComplete
	if verifyComplete == "" {
		verifyComplete = dev.VerificationURI
	}

	log.Info("---------------------------------\n\n")
	log.Info("Open this URL in your browser and log in:")
	log.Infof("  %s", verifyComplete)
	log.Infof("Or visit %s and enter code: %s\n\n", dev.VerificationURI, dev.UserCode)
	log.Info("-----------------------------")

	return poll(ctx, tokenEndpoint, dev.DeviceCode, sett.Client, pkceVerifier, interval, expiresIn)
}

// requestDeviceCode requests a device code and user code.
func requestDeviceCode(
	ctx context.Context,
	sett *Settings,
	deviceEndpoint string,
	pkceChallenge string,
) (*deviceResponse, error) {
	initBody := url.Values{}
	initBody.Set("client_id", sett.Client)
	initBody.Set("scope", sett.Scope)
	if pkceChallenge != "" {
		initBody.Set("code_challenge", pkceChallenge)
		initBody.Set("code_challenge_method", "S256")
	}

	log.Infof("→ POST %s", deviceEndpoint)
	log.Infof("Init body: '%v'", initBody)

	status, body, err := postForm(ctx, deviceEndpoint, initBody)
	if err != nil {
		return nil, errors.AddContext(err, "device authorization request failed")
	}

	if status != http.StatusOK {
		log.Errorf("Device authorization request failed (HTTP %d):", status)
		log.Error(string(body))

		log.Warn("Hint: ensure the client is public and the Device Authorization Grant is enabled.")
		if sett.Provider == providerAuthentik {
			log.Warn(
				"For authentik, the default brand must have `flow_device_code` set (see modos-blueprint.yaml).",
			)
		}

		return nil, errors.New("device authorization request failed (HTTP %d)", status)
	}

	var dev deviceResponse
	if e := json.Unmarshal(body, &dev); e != nil {
		return nil, errors.AddContext(e, "could not decode device authorization response")
	}

	return &dev, nil
}

// poll repeatedly queries the token endpoint until authorization completes,
// the code expires, or the context is cancelled.
func poll(
	ctx context.Context,
	tokenEndpoint string,
	deviceCode string,
	clientID string,
	verifier string,
	interval int,
	expiresIn int,
) error {
	pollBody := url.Values{}
	pollBody.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
	pollBody.Set("device_code", deviceCode)
	pollBody.Set("client_id", clientID)
	if verifier != "" {
		pollBody.Set("code_verifier", verifier)
	}

	wait := time.Duration(interval) * time.Second
	deadline := time.Now().Add(time.Duration(expiresIn) * time.Second)

	log.Infof("Polling '%s' for token every %s (code expires in %ds)...\n",
		tokenEndpoint, wait, expiresIn)
	log.Infof("Body: '%v'", pollBody)

	for {
		if time.Now().After(deadline) {
			return errors.New("device code expired before authorization completed")
		}

		status, body, err := postForm(ctx, tokenEndpoint, pollBody)
		if err != nil {
			return errors.AddContext(err, "token request failed")
		}

		if status == http.StatusOK {
			return reportSuccess(body)
		}

		var e tokenError
		_ = json.Unmarshal(body, &e)
		errCode := e.Error
		if errCode == "" {
			errCode = "unknown_error"
		}

		switch errCode {
		case "authorization_pending":
			log.Infof("  … waiting for user to authorize (HTTP %d)", status)
		case "slow_down":
			wait += slowDownIncrement
			log.Infof("  … slow_down — increasing poll interval to %s", wait)
		default:
			log.Errorf("✗ Token request failed: %s", errCode)
			log.Error(string(body))

			return errors.New("token request failed: %s", errCode)
		}

		select {
		case <-ctx.Done():
			log.Warn("Aborted the polling.")

			return ctx.Err()
		case <-time.After(wait):
		}
	}
}

// reportSuccess prints the obtained token and best-effort decodes its claims.
func reportSuccess(body []byte) error {
	var tok tokenResponse
	if err := json.Unmarshal(body, &tok); err != nil {
		return errors.AddContext(err, "could not decode token response")
	}

	hasRefresh := "no"
	if tok.RefreshToken != "" {
		hasRefresh = "yes"
	}

	log.Info("✅ Authorization successful!")
	log.Infof("  token_type    : %s", tok.TokenType)
	log.Infof("  expires_in    : %d", tok.ExpiresIn)
	log.Infof("  refresh_token : %s", hasRefresh)
	log.Infof("  access_token  : %s", tok.AccessToken)

	// Best-effort decode of the JWT access-token payload.
	claims, err := decodeJWTClaims(tok.AccessToken)
	if err != nil {
		return errors.AddContext(err, "JWT payload not decoded")
	}
	log.Infof("Access-token claims:\n%s", claims)

	return nil
}

// newPKCE generates a random verifier and its base64url-encoded S256 challenge.
// The verifier itself is presented later at the token-polling step.
func newPKCE() (verifier, challenge string, err error) {
	buf := make([]byte, pkceVerifierBytes)
	if _, e := rand.Read(buf); e != nil {
		return "", "", errors.AddContext(e, "could not generate PKCE verifier")
	}
	verifier = base64.RawURLEncoding.EncodeToString(buf)

	sum := sha256.Sum256([]byte(verifier))
	challenge = base64.RawURLEncoding.EncodeToString(sum[:])

	return verifier, challenge, nil
}

// decodeJWTClaims decodes and pretty-prints the payload segment of a JWT.
func decodeJWTClaims(accessToken string) (string, error) {
	parts := strings.Split(accessToken, ".")
	if len(parts) < jwtParts {
		return "", errors.New("access token is not a JWT")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", errors.AddContext(err, "could not base64url-decode JWT payload")
	}

	var claims any
	if e := json.Unmarshal(payload, &claims); e != nil {
		return "", errors.AddContext(e, "could not parse JWT payload as JSON")
	}

	pretty, err := json.MarshalIndent(claims, "", "  ")
	if err != nil {
		return "", errors.AddContext(err, "could not re-encode JWT claims")
	}

	return string(pretty), nil
}

// postForm sends a form-encoded POST and returns the status code and raw body.
// Non-2xx responses are returned without error (like nushell's `--allow-errors`).
func postForm(
	ctx context.Context,
	endpoint string,
	form url.Values,
) (status int, body []byte, err error) {
	req, err := http.NewRequestWithContext(
		ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return 0, nil, errors.AddContext(err, "could not build request to '%s'", endpoint)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, nil, errors.AddContext(err, "request to '%s' failed", endpoint)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, errors.AddContext(err,
			"could not read response body from '%s'", endpoint)
	}

	return resp.StatusCode, body, nil
}
