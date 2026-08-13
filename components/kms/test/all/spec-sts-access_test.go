//go:build test && integration

package all

import (
	"testing"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwt"
	. "github.com/onsi/ginkgo/v2"
	mdJwt "github.com/sdsc-ordes/modos-rs/components/kms/internal/jwt/test"
	st "github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage/types"
	"github.com/sdsc-ordes/quitsh/pkg/log"
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
				mdJwt.WithSign(testCtx.Authentik.JWTPrivateKey),
				mdJwt.WithModifications(func(b *jwt.Builder) {
					b.Claim("bp", []st.BucketPermissions{})
				}),
			)

			// JWT validate -> no - because not needed.
			// JWT -> serialize to JSON
			// JSON -> parse with lib-common claim `bp`.

			c, err := testCtx.Storage.NewCredentials(
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
