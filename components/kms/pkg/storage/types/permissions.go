package types

import (
	"slices"
	"strings"
)

const PermissionRead Permission = "read"
const PermissionWrite Permission = "write"

type (
	Permission  = string
	Permissions []Permission

	BucketPermission struct {
		// The bucket permission resource path
		// (without '*' or '?' due to safetey)
		// E.g `bucket-a/bla/a/b/c
		Path string

		// The permissions for this bucket.
		Permissions Permissions
	}

	BucketPermissions = []BucketPermission
)

// Contains checks if permissions contains permission `p`.
func (ps Permissions) Contains(p Permission) bool {
	return slices.Contains(ps, p)
}

// Bucket returns the bucket string.
func (p *BucketPermission) Bucket() string {
	s := strings.SplitN(
		strings.TrimLeft(p.Path, "/"),
		"/", 2) //nolint:mnd

	return s[0]
}
