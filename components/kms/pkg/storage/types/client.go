package types

import (
	"context"
	"time"
)

type (
	Client interface {
		// Ping returns if the storage can be accessed.
		Ping(ctx context.Context) error

		// NewCredentials returns a scoped token to bucket `bucket` with expiration `expiration` and permissions `permissions`.
		NewCredentials(
			ctx context.Context,
			bucket string,
			permissions []Permission,
			expiration time.Time) (creds Credentials, err error)
	}
)
