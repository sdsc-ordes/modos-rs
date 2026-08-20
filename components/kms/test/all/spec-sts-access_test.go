//go:build test && integration

package all

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwt"
	. "github.com/onsi/ginkgo/v2"
	mdJwtT "github.com/sdsc-ordes/modos-rs/components/kms/internal/jwt/test"
	st "github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage/types"
	"github.com/sdsc-ordes/quitsh/pkg/log"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/auth"
)

var _ = Describe("S3", func() {
	var t testing.TB

	BeforeEach(func() {
		t = GinkgoTB()
	})

	Describe("requesting STS token triplet", Label("sts-access"), func() {
		It("should give access to bucket-b", func() {
			testCtx := NewTestContext(t)
			defer testCtx.Close(t)

			// - Create a dummy JWT with Permissions.
			token := testCtx.NewTokenKeycloak(t,
				mdJwtT.WithSign(testCtx.Authentik.JWTPrivateKey),
				mdJwtT.WithModifications(func(b *jwt.Builder) {
					b.Claim("bp", st.BucketPermissions{
						st.BucketPermission{
							Path:        "bucket-a",
							Permissions: []st.Permission{st.PermissionWrite}},
					})
				}),
			)

			tokenS, err := json.Marshal(token)

			cl, err := auth.ValidateJWT[mdJwt.Claims](testCtx.Ctx, testCtx.JWTVerifier, tokenS, nil)

			// provider, err := oidc.NewProvider(ctx, issuer) // fetches .well-known/openid-configuration
			// // endpoints for oauth2.Config:
			// endpoint := provider.Endpoint()
			// // verifier (handles JWKS fetch internally):
			// verifier := provider.Verifier(&oidc.Config{ClientID: clientID})
			// idToken, err := verifier.Verify(ctx, rawIDToken)
			//
			// // grab extra discovery fields not exposed as methods:
			// var extra struct {
			// 	ScopesSupported []string `json:"scopes_supported"`
			// }
			// _ = provider.Claims(&extra)
			// Which to pick

			// JWT validate -> no - because not needed.
			// JWT -> serialize to JSON
			// JSON -> parse with lib-common claim `bp`.

			_, err := testCtx.Storage.NewCredentials(
				testCtx.Ctx,
				[]st.BucketPermission{},
				1*time.Hour,
			)
			if err != nil {
				log.ErrorE(err, "Credentials could not be created.")
			}

			// - Create a STS Credential
			// - Write file to bucket and test it works.
		})

		It("should not give access to bucket-a", func() {

		})
	})
})
