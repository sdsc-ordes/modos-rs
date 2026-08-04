package types

const PermissionRead Permission = "read"
const PermissionWrite Permission = "write"

type (
	Permission string

	BucketPermission struct {
		// The bucket permission resource path (without '*' or '?' due to safetey)
		Path string
		// The permissions for this bucket.
		Permissions []Permission
	}

	BucketPermissions = []BucketPermission
)
