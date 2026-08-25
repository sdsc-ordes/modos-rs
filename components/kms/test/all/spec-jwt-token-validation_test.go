//go:build test && integration

package all

import (
	"encoding/json"
	"testing"

	"github.com/lestrrat-go/jwx/v3/jwt"
	. "github.com/onsi/ginkgo/v2"
	mdJwt "github.com/sdsc-ordes/modos-rs/components/kms/internal/jwt"
	mdJwtT "github.com/sdsc-ordes/modos-rs/components/kms/internal/jwt/test"
	st "github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage/types"
	"github.com/stretchr/testify/require"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/auth"
)

// CreateToken creates a token with bucket permissions claims.
func CreateToken(
	t testing.TB,
	testCtx *TestContext,
	permissions st.BucketPermissions,
) (jwt.Token, string) {
	token := testCtx.NewToken(t,
		mdJwtT.WithBucketPermissionsClaim(
			&testCtx.Cfg.OIDC.ClaimBucketPermissions,
			permissions,
		),
	)

	signedToken := mdJwtT.SignToken(t, token, testCtx.OIDC.JWTPrivateKey)

	return token, signedToken
}

// Run with
//
// ```bash
//
//	just quitsh exec-target \
//		--log-level debug
//		-K "test.showTestLog: true" \
//		-K 'test.testArgs: [ "-ginkgo.label-filter=jwt" ]'
//		"kms::test-integration"
//
// ```
var _ = Describe("A JWT", func() {
	var t testing.TB

	BeforeEach(func() {
		t = GinkgoTB()
	})

	Describe("with normal bucket permissions", Label("jwt", "bucket-permissions"), func() {
		It("should validate", func() {
			testCtx := NewTestContext(t)
			defer testCtx.Close(t)

			token, signedToken := CreateToken(
				t, testCtx, st.BucketPermissions{
					st.BucketPermission{
						Path:        "bucket-a",
						Permissions: []st.Permission{st.PermissionWrite},
					},
				})

			b, err := json.Marshal(token)
			require.NoError(t, err)

			cl := mdJwt.NewClaims(&testCtx.Cfg.OIDC.ClaimBucketPermissions)
			err = auth.ValidateJWT(
				testCtx.Ctx, testCtx.JWTVerifier, signedToken, nil, cl,
			)
			require.NoError(t, err, "JWT: '%v'", string(b))
		})
	})
})
