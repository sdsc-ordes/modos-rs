//go:build test && integration

package all

import (
	"encoding/json"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	mdJwt "github.com/sdsc-ordes/modos-rs/components/kms/internal/jwt"
	mdJwtT "github.com/sdsc-ordes/modos-rs/components/kms/internal/jwt/test"
	st "github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage/types"
	"github.com/stretchr/testify/require"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/auth"
)

var _ = Describe("JWT", func() {
	var t testing.TB

	BeforeEach(func() {
		t = GinkgoTB()
	})

	Describe("with normal bucket permissions", Label("jwt", "bucket-permissions"), func() {
		It("should validate", func() {
			testCtx := NewTestContext(t)
			defer testCtx.Close(t)

			token := testCtx.NewToken(t,
				mdJwtT.WithBucketPermissionsClaim(
					&testCtx.Cfg.OIDC.ClaimBucketPermissions,
					st.BucketPermissions{
						st.BucketPermission{
							Path:        "bucket-a",
							Permissions: []st.Permission{st.PermissionWrite},
						},
					}))

			signedToken := mdJwtT.SignToken(t, token, testCtx.OIDC.JWTPrivateKey)

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
