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

	type D struct {
		perm      BucketPermission
		bucketExp string
		restExp   string
	}

	It("bucket name should work", func() {
		tests := []D{
			{
				perm:      BucketPermission{Path: "bucket-a/a/b/c"},
				bucketExp: "bucket-a",
				restExp:   "a/b/c",
			},
			{
				perm:      BucketPermission{Path: "/bucket-a///"},
				bucketExp: "bucket-a",
				restExp:   "/",
			},
			{
				perm:      BucketPermission{Path: "/bucket-a/a/b///a"},
				bucketExp: "bucket-a",
				restExp:   "a/b/a",
			},
		}

		for _, d := range tests {
			assert.Equal(t, d.bucketExp, d.perm.Bucket())
			bucket, rest := d.perm.PathSplit()
			assert.Equal(t, d.bucketExp, bucket)
			assert.Equal(t, d.restExp, rest)
		}
	})
})
