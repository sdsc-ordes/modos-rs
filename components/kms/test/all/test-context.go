//go:build test && integration

package all

import (
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"path"
	"testing"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/sdsc-ordes/modos-rs/components/kms/internal/config"
	mdJwt "github.com/sdsc-ordes/modos-rs/components/kms/internal/jwt/test"
	"github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage"
	st "github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage/types"
	"github.com/sdsc-ordes/modos-rs/components/kms/test/common"
	"github.com/stretchr/testify/require"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/auth"
	cmc "gitlab.com/data-custodian/custodian/components/lib-common/pkg/config"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/log"
	clog "gitlab.com/data-custodian/custodian/components/lib-common/pkg/log/context"
)

type (
	TestContext struct {
		Ctx context.Context

		Log     log.Logger
		RootDir string

		Keycloak  OAuth
		Authentik OAuth

		Storage     st.Client
		JWTVerifier *auth.JWTVerifier
	}

	OAuth struct {
		JWTPrivateKey jwk.Key
	}

	TestContextOption func(*testContextOpts)
	testContextOpts   struct{}
)

// NewTestContext sets up test context.
func NewTestContext(t testing.TB, opts ...TestContextOption) (testCtx *TestContext) {
	log.Info("Construct new test context.")
	var o testContextOpts
	o.Apply(opts...)

	ctx := clog.ContextFrom(context.Background(), log.GetLoggerOrig().Named("kms"))

	rootDir := common.GetTestRootDir()
	dataDir := path.Join(rootDir, "data")
	configDir := path.Join(dataDir, "config")

	conf, err := cmc.LoadConfigs[config.Config](configDir)
	require.NoError(t, err, "Could not load config files in '%s'.", configDir)
	conf.WithDataDir(dataDir)

	client, err := storage.NewStorageS3(ctx, &conf.Storage.Connection)
	log.PanicEf(err, "Could not create S3 storage.")

	jwtVerifier, err := createJWTVerifier(ctx, &conf.OIDC)
	log.PanicEf(err, "Could not create JWT verifier.")

	return &TestContext{
		Ctx:         ctx,
		RootDir:     rootDir,
		Storage:     client,
		Keycloak:    OAuth{JWTPrivateKey: getKeycloakPrivateKey(t)},
		Authentik:   OAuth{JWTPrivateKey: getAuthentikPrivateKey(t)},
		JWTVerifier: jwtVerifier,
	}
}

func (c *TestContext) Close(t testing.TB) {
}

func (c *testContextOpts) Apply(options ...TestContextOption) {
	for _, f := range options {
		f(c)
	}
}

// NewTokenKeycloak returns a signed token with default values.
func (c *TestContext) NewTokenKeycloak(t testing.TB, options ...mdJwt.Option) jwt.Token {
	options = append(
		options,
		mdJwt.WithModifications(func(b *jwt.Builder) {
			b.Audience([]string{mdJwt.DefaultIssuerKeycloak})
		}))

	return mdJwt.NewToken(t, options...)
}

// NewTokenAuthentik returns a signed token with default values.
func (c *TestContext) NewTokenAuthentik(t testing.TB, options ...mdJwt.Option) jwt.Token {
	options = append(options, mdJwt.WithModifications(func(b *jwt.Builder) {
		b.Audience([]string{mdJwt.DefaultIssuerAuthentik})
	}))

	return mdJwt.NewToken(t, options...)
}

func getAuthentikPrivateKey(t testing.TB) jwk.Key {
	// Note: same as in `modos-oauth-blueprint.yaml`
	pemKey := `-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgE2rYNDdsy50q8Z+c
FJRWAy7ZXtn8EOYpYhrLH3P425+hRANCAAS2qulHIyTz2hobJ3eBxzhd4+iH2Kxj
nyK+Tb15OjCIhIFTqknYaMyQcsu+1btcpvR9E8KlMsTv7awVoA5+9+7z
-----END PRIVATE KEY-----`

	block, _ := pem.Decode([]byte(pemKey))
	require.NotNil(t, block, "Failed to decode PEM block.")

	rawKey, e := x509.ParsePKCS8PrivateKey(block.Bytes)
	require.NoError(t, e, "failed to parse key as PKCS#8.")
	privKey, ok := rawKey.(*ecdsa.PrivateKey)
	require.True(t, ok, "Failed to convert key to RSA private key.")

	key, err := jwk.Import(privKey)
	require.NoError(t, err, "Could not import key.")
	err = key.Set("alg", "ES256")
	require.NoError(t, err, "Could not set `alg` field.")
	jwk.AssignKeyID(key)
	require.NoError(t, err, "Could not set `kid` field.")

	return key
}

func getKeycloakPrivateKey(t testing.TB) jwk.Key {
	// Note: same as in `modos-realm.json`
	b64Key := "MC4CAQAwBQYDK2VwBCIEIGaiqDys9Gpq+mGrDeFus9q9WTmrD9x8rJW1Bvdynbix"

	derBytes, err := base64.StdEncoding.DecodeString(b64Key)
	require.NoError(t, err, "Failed to decode key.")

	rawKey, e := x509.ParsePKCS8PrivateKey(derBytes)
	require.NoError(t, e, "failed to parse key as PKCS#8.")
	privKey, ok := rawKey.(ed25519.PrivateKey)
	require.True(t, ok, "Failed to convert key to ed25519 private key.")

	key, err := jwk.Import(privKey)
	require.NoError(t, err, "Could not import key.")
	err = key.Set("alg", "EdDSA")
	require.NoError(t, err, "Could not set `alg` field.")
	jwk.AssignKeyID(key)
	require.NoError(t, err, "Could not set `kid` field.")

	return key
}

func createJWTVerifier(ctx context.Context, oidcCfg *config.OIDC) (*auth.JWTVerifier, error) {
	const timeout = 30 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return auth.NewJWTVerifier(
		ctx,
		oidcCfg.Issuer,
		oidcCfg.ClientID,
		auth.WithTrustedAlgorithms(oidcCfg.TrustedAlgorithms...),
		auth.WithTrustedAudiences(oidcCfg.TrustedAudiences...),
	)
}
