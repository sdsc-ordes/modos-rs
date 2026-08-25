//go:build test && integration

package all

import (
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("S3", func() {
	// var t testing.TB

	BeforeEach(func() {
		// t = GinkgoTB()
	})

	Describe("requesting STS token triplet", Label("sts-access"), func() {
		It("should give access to bucket-b", func() {
			// testCtx := NewTestContext(t)
			// defer testCtx.Close(t)
			//
			// // - Create a dummy JWT with Permissions.
			// token := testCtx.NewTokenKeycloak(t,
			// 	mdJwtT.WithModifications(func(b *jwt.Builder) {
			// 		b.Claim("bp",
			// 			st.BucketPermissions{
			// 				st.BucketPermission{
			// 					Path:        "bucket-a",
			// 					Permissions: []st.Permission{st.PermissionWrite}},
			// 			})
			// 	}),
			// )
			// b, err := json.Marshal(token)
			//
			// signedToken := mdJwtT.SignToken(t, token, testCtx.Keycloak.JWTPrivateKey)
			//
			// cl := mdJwt.NewClaims(&testCtx.Cfg.OIDC.ClaimBucketPermissions)
			// err := auth.ValidateJWT(testCtx.Ctx, testCtx.JWTVerifier, signedToken, nil, cl)
			// require.NoError(t, err)

		})
	})
})
