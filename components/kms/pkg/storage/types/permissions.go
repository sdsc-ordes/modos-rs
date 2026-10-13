package types

import (
	"path"
	"slices"
	"strings"

	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/errors"
)

const PermissionRead Permission = "read"
const PermissionWrite Permission = "write"

type (
	Permission  = string
	Permissions []Permission
	Path        string

	BucketPermission struct {
		// The bucket permission resource path
		// (without '*' or '?' due to safetey)
		// E.g `bucket-a/bla/a/b/c
		Path Path

		// The permissions for this bucket.
		Permissions Permissions
	}

	BucketPermissions = []BucketPermission
)

// Contains checks if permissions contains permission `p`.
func (ps Permissions) Contains(p Permission) bool {
	return slices.Contains(ps, p)
}

// Bucket returns the bucket name.
func (p Path) Bucket() string {
	s := strings.SplitN(
		strings.TrimLeft((string)(p), "/"),
		"/", 2) //nolint:mnd

	return s[0]
}

// Split returns the bucket name and rest of the path.
func (p Path) Split() (string, string) {
	s := strings.SplitN(
		strings.TrimLeft((string)(p), "/"),
		"/", 2) //nolint:mnd

	if len(s) == 1 {
		return s[0], ""
	}

	return s[0], path.Clean(s[1])
}

// Sanitize sanitzes the path.
func (p Path) Sanitize() (Path, error) {
	bucket, rest := p.Split()

	if bucket == "" {
		return "", errors.New("Bucket name '%v' results in an empty resource path -> Ignore.", p)
	} else if strings.ContainsAny(bucket, "*?") {
		return "", errors.New("Bucket name '%v' contains '?' or '*' which is not supported yet.", p)
	}

	return Path(path.Join(bucket, rest)), nil
}
