//go:build test && unittest

package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/onsi/ginkgo/v2"
)

func TestTypes(t *testing.T) {
	RunSpecs(t, "storage types")
}

var _ = Describe("permissions", func() {
	var t require.TestingT
	BeforeEach(func() { t = GinkgoT() })

	It("bucket name should work", func() {
		b := BucketPermission{Path: "bucket-a/a/b/c"}
		assert.Equal(t, "bucket-a", b.Bucket())

		b = BucketPermission{Path: "/bucket-a///"}
		assert.Equal(t, "bucket-a", b.Bucket())
	})
})
