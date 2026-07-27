package s3

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage/types"

	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/errors"
	clog "gitlab.com/data-custodian/custodian/components/lib-common/pkg/log/context"
)

// Ping implements [types.Client].
func (c *clientS3) Ping(ctx context.Context) (err error) {
	clog.Infof(ctx, "Pinging buckets.")

	ctxT, cancel := context.WithTimeout(ctx, defaultPingTimeout)
	defer cancel()

	_, e := c.client.ListBuckets(ctxT, &s3.ListBucketsInput{})

	if e != nil {
		return errors.AddContext(e,
			"could not ping buckets (timeout: '%v')",
			defaultPingTimeout)
	}

	return
}

// Credentials implements [types.Client].
func (c *clientS3) NewCredentials(
	ctx context.Context,
	bucket string,
	permissions []types.Permission,
	expiration time.Time,
) (creds types.Credentials, err error) {
	// https://docs.rustfs.com/administration/iam/sts#sts-temporary-credentials
	return types.Credentials{}, nil
}
