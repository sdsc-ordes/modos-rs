package types

import (
	"context"
	"time"
)

type (
	Client interface {
		// Ping returns if the storage can be accessed.
		Ping(ctx context.Context) error

		// NewCredentials returns a scoped token for bucket permissions `permissions`
		// with duration `dur`.
		NewCredentials(
			ctx context.Context,
			perms BucketPermissions,
			duration time.Duration,
		) (creds Credentials, err error)
	}
)
