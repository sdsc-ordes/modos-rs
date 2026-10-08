//go:build test && unittest

package s3

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage/types"
)

var _ = Describe("policy", func() {
	var t require.TestingT
	BeforeEach(func() { t = GinkgoT() })

	Describe("toAction", func() {
		ctx := context.Background()

		It("maps read to object-get and bucket-list", func() {
			assert.Equal(t,
				[]string{s3GetObject, s3ListBucket},
				toAction(ctx, []types.Permission{types.PermissionRead}))
		})

		It("maps write to object-put", func() {
			assert.Equal(t,
				[]string{s3PutObject},
				toAction(ctx, []types.Permission{types.PermissionWrite}))
		})

		It("combines read and write actions in order", func() {
			assert.Equal(t,
				[]string{s3GetObject, s3ListBucket, s3PutObject},
				toAction(ctx, []types.Permission{types.PermissionRead, types.PermissionWrite}))
		})

		It("deduplicates repeated permissions", func() {
			assert.Equal(t,
				[]string{s3GetObject, s3ListBucket},
				toAction(ctx, []types.Permission{types.PermissionRead, types.PermissionRead}))
		})

		It("skips unknown permissions", func() {
			assert.Empty(t,
				toAction(ctx, []types.Permission{types.Permission("bogus")}))
		})

		It("returns nothing for no permissions", func() {
			assert.Empty(t, toAction(ctx, nil))
		})
	})

	Describe("NewScopedPolicy", func() {
		ctx := context.Background()

		It("uses the current IAM policy language version", func() {
			doc, err := NewScopedPolicy(ctx, types.BucketPermissions{})
			require.NoError(t, err)
			assert.Equal(t, "2012-10-17", doc.Version)
			assert.Empty(t, doc.Statement)
		})

		It("builds one allow statement scoped to the sanitized subpath", func() {
			doc, err := NewScopedPolicy(ctx,
				types.BucketPermissions{
					{
						Path:        "/bucket-a/sub/",
						Permissions: []types.Permission{types.PermissionRead}},
				})
			require.NoError(t, err)
			require.Len(t, doc.Statement, 1)

			st := doc.Statement[0]
			assert.Equal(t, "Allow", st.Effect)
			assert.Equal(t, []string{s3GetObject, s3ListBucket}, st.Action)
			assert.Equal(t, []string{"arn:aws:s3:::bucket-a/sub/*"}, st.Resource)
		})

		It("emits one statement per bucket permission", func() {
			doc, err := NewScopedPolicy(ctx, types.BucketPermissions{
				{Path: "bucket-a/read", Permissions: []types.Permission{types.PermissionRead}},
				{Path: "bucket-b/write", Permissions: []types.Permission{types.PermissionWrite}},
			})
			require.NoError(t, err)
			require.Len(t, doc.Statement, 2)

			assert.Equal(t, []string{"arn:aws:s3:::bucket-a/read/*"}, doc.Statement[0].Resource)
			assert.Equal(t, []string{s3PutObject}, doc.Statement[1].Action)
			assert.Equal(t, []string{"arn:aws:s3:::bucket-b/write/*"}, doc.Statement[1].Resource)
		})

		It("skips permissions whose bucket path is invalid", func() {
			doc, err := NewScopedPolicy(ctx, types.BucketPermissions{
				{Path: "bad*bucket/sub", Permissions: []types.Permission{types.PermissionRead}},
				{Path: "bucket-b/ok", Permissions: []types.Permission{types.PermissionWrite}},
			})
			require.NoError(t, err)
			require.Len(t, doc.Statement, 1)
			assert.Equal(t, []string{"arn:aws:s3:::bucket-b/ok/*"}, doc.Statement[0].Resource)
		})
	})
})
