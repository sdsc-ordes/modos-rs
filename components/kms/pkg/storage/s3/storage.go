package s3

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage/types"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/secret"

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
	duration time.Duration,
) (creds *types.Credentials, err error) {
	// For RustFS: https://docs.rustfs.com/administration/iam/sts#sts-temporary-credentials
	// NOTE: - Service Accounts (an 'access-key' bound to a user) cannot call `AssumeRole`
	// therefore we use a dedicated user account.

	in := sts.AssumeRoleInput{
		RoleArn:         aws.String("arn:aws:iam::rustfs:role/scoped"),
		RoleSessionName: aws.String("bucket-access"),
		DurationSeconds: aws.Int32(int32(duration.Seconds())),
		// Policy: json.Marshal(),
	}

	res, err := c.sts.AssumeRole(ctx, &in)
	if err != nil {
		return nil, errors.AddContext(
			err,
			"Could not call 'AssumeRole' in getting new credentials.",
		)
	}

	return &types.Credentials{
		AccessKeyID:     secret.RedactedString(*res.Credentials.AccessKeyId),
		SessionToken:    secret.RedactedString(*res.Credentials.SessionToken),
		SecretAccessKey: secret.RedactedString(*res.Credentials.SecretAccessKey),
	}, nil
}
