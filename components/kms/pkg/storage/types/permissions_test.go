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

	Describe("Sanitize path", func() {
		It("trims surrounding slashes", func() {
			got, err := Path("/bucket-a/sub/path/").Sanitize()
			require.NoError(t, err)
			assert.Equal(t, Path("bucket-a/sub/path"), got)
		})

		It("keeps an already-clean path unchanged", func() {
			got, err := Path("bucket-a/sub").Sanitize()
			require.NoError(t, err)
			assert.Equal(t, Path("bucket-a/sub"), got)
		})

		It("rejects a wildcard in the bucket (first) segment", func() {
			for _, p := range []string{"buck*et/sub", "bu?cket", "*/sub"} {
				_, err := Path(p).Sanitize()
				require.Error(t, err, "path %q must be rejected", p)
				assert.Contains(t, err.Error(), "not supported")
			}
		})

		It("allows a wildcard beyond the first segment (only the bucket is guarded)", func() {
			got, err := Path("bucket-a/su*b").Sanitize()
			require.NoError(t, err)
			assert.Equal(t, Path("bucket-a/su*b"), got)
		})
	})

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
			assert.Equal(t, d.bucketExp, d.perm.Path.Bucket())
			bucket, rest := d.perm.Path.Split()
			assert.Equal(t, d.bucketExp, bucket)
			assert.Equal(t, d.restExp, rest)
		}
	})
})
