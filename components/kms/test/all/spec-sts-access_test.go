//go:build test && integration

package all

import (
	"log"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	mdJwt "github.com/sdsc-ordes/modos-rs/components/kms/internal/jwt"
	st "github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage/types"
	"github.com/stretchr/testify/require"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/auth"
)

// Run with
//
// ```bash
//
//	just quitsh exec-target \
//		--log-level debug
//		-K "test.showTestLog: true" \
//		-K 'test.testArgs: [ "-ginkgo.label-filter=storage-access" ]'
//		"kms::test-integration"
//
// ```
var _ = Describe("S3", func() {
	var t testing.TB

	BeforeEach(func() {
		t = GinkgoTB()
	})

	Describe("requesting credentials (STS)", Label("storage-access"), func() {
		It("should give access to bucket-b", func() {
			testCtx := NewTestContext(t)
			defer testCtx.Close(t)

			_, signedToken := CreateToken(
				t, testCtx, st.BucketPermissions{
					st.BucketPermission{
						Path:        "bucket-a",
						Permissions: []st.Permission{st.PermissionWrite},
					},
				})

			cl := mdJwt.NewClaims(&testCtx.Cfg.OIDC.ClaimBucketPermissions)
			err := auth.ValidateJWT(testCtx.Ctx, testCtx.JWTVerifier, signedToken, nil, cl)
			require.NoError(t, err)

			cred, err := getStorageCredential(t, testCtx, cl)
			require.NoError(t, err)

			err = testStorageAccess(t, testCtx, cred, true)
			require.NoError(t, err)
		})
	})
})

func getStorageCredential(
	_ testing.TB,
	_ *TestContext,
	_ *mdJwt.Claims,
) (st.Credentials, error) {
	log.Panicf("Not implemented")

	return nil, nil
}

func testStorageAccess(
	_ testing.TB,
	_ *TestContext,
	_ st.Credentials,
	_ bool, /* test write flag */
) error {
	log.Panicf("Not implemented")

	return nil
}
