package types

import (
	"context"
	"time"
)

type (
	Client interface {
		// Ping returns if the storage can be accessed.
		Ping(ctx context.Context) error

		// NewCredentials returns a scoped token to bucket `bucket`
		// with duration `duration` and permissions `permissions`.
		NewCredentials(
			ctx context.Context,
			bucket string,
			permissions []Permission,
			duration time.Duration) (creds *Credentials, err error)
	}
)
