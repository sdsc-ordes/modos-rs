package types

import (
	"context"
)

type (
	Client interface {
		// Ping returns if the storage can be accessed.
		Ping(ctx context.Context) error
	}
)
