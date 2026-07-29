package all

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("S3", func() {
	var t testing.TB
	BeforeEach(func() { t = GinkgoTB() })

	Describe("requestin STS token triplet", func() {
		It("should give access to bucket-a", func() {
			t.Log("Do something.")
		})
	})
})
