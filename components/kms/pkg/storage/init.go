package storage

import (
	"context"

	"github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage/s3"
	st "github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage/types"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/errors"
)

// NewStorageS3 initializes an object storage (S3) client.
func NewStorageS3(ctx context.Context, conf *st.S3Connection) (st.Client, error) {
	store, err := s3.NewClient(ctx, conf)
	if err != nil {
		return nil, errors.AddContext(err, "Could not setup storage client.")
	}

	err = store.Ping(ctx)
	if err != nil {
		return nil, errors.AddContext(err, "Could not ping attachment S3.")
	}

	return store, nil
}
