package s3

import (
	"context"
	"encoding/json"
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
	permissions types.BucketPermissions,
	duration time.Duration,
) (creds *types.Credentials, err error) {
	clog.Info(ctx,
		"Create credential for permissions.",
		"permissions", permissions, "duration", duration)

	// For RustFS: https://docs.rustfs.com/administration/iam/sts#sts-temporary-credentials
	// NOTE:
	// - Service Accounts (an 'access-key' bound to a user) cannot call `AssumeRole`
	// therefore we use a dedicated user account.
	policy, err := NewScopedPolicy(ctx, permissions)
	if err != nil {
		return nil, errors.AddContext(err, "Could not create scoped policy document.")
	}

	policyJSON, err := json.MarshalIndent(policy, "", "  ")
	if err != nil {
		return nil, errors.AddContext(err, "Could marshal scoped policy document.")
	}
	clog.Info(ctx, "Create credentials with policy.", "role", policy)

	in := sts.AssumeRoleInput{
		RoleArn:         aws.String("arn:aws:iam::rustfs:role/scoped"),
		RoleSessionName: aws.String("bucket-access"),
		DurationSeconds: aws.Int32(int32(duration.Seconds())),
		Policy:          aws.String(string(policyJSON)),
	}

	res, err := c.sts.AssumeRole(ctx, &in)
	if err != nil {
		return nil, errors.AddContext(
			err,
			"Could not call 'AssumeRole' in getting new credentials.",
		)
	}
	clog.Info(ctx, "Create credentials successful.")

	return &types.Credentials{
		AccessKeyID:     secret.RedactedString(*res.Credentials.AccessKeyId),
		SessionToken:    secret.RedactedString(*res.Credentials.SessionToken),
		SecretAccessKey: secret.RedactedString(*res.Credentials.SecretAccessKey),
	}, nil
}
