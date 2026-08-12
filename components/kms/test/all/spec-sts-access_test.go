//go:build test && integration

package all

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("S3", func() {
	var t testing.TB
	BeforeEach(func() { t = GinkgoTB() })

	Describe("requesting STS token triplet", Label("sts-access"), func() {
		It("should give access to bucket-a", func() {
			tCtx := NewTestContext(t)
			defer tCtx.Close(t)

			t.Log("Do something.", tCtx.Keycloak)
		})
	})
})
