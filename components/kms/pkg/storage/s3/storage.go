package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"

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
