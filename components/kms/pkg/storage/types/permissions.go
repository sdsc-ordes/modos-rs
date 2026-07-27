package types

type Permission string

const PermissionRead Permission = "read"
const PermissionWrite Permission = "write"

type BucketPermission struct {
	// The bucket permission resource path (without '*' or '?' due to safetey)
	Path string
	// The permissions for this bucket.
	Permissions []Permission
}

type BucketPermissions = []BucketPermission
