package s3

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage/types"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/secret"

	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/errors"
	clog "gitlab.com/data-custodian/custodian/components/lib-common/pkg/log/context"
)

// Ping implements [types.Client].
func (c *clientS3) Ping(ctx context.Context) error {
	clog.Infof(ctx, "Pinging buckets.")

	ctxT, cancel := context.WithTimeout(ctx, defaultPingTimeout)
	defer cancel()

	_, e := c.client.ListBuckets(ctxT, &s3.ListBucketsInput{}) //nolint: exhaustruct // intended

	if e != nil {
		return errors.AddContext(e,
			"could not ping buckets (timeout: '%v')",
			defaultPingTimeout)
	}

	return nil
}

func (c *clientS3) UploadTest(
	ctx context.Context,
	bucketName string,
	creds types.Credentials,
) error {
	reader := bytes.NewReader([]byte("This is an upload test."))
	objectKey := "test-object.txt"

	clog.Infof(ctx, "Starting to upload test object.")

	s3Creds, ok := creds.(*S3Credentials)
	if !ok {
		return errors.New("credentials must be 'S3Credentials'")
	}

	_, err := c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
		Body:   reader,
	}, func(o *s3.Options) {
		o.Credentials = credentials.NewStaticCredentialsProvider(
			string(s3Creds.AccessKeyID),
			string(s3Creds.SecretAccessKey),
			string(s3Creds.SessionToken),
		)
	})
	if err != nil {
		return errors.AddContext(err, "failed to upload test object")
	}

	err = s3.NewObjectExistsWaiter(c.client).Wait(
		ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey)},
		time.Minute)
	if err != nil {
		return errors.AddContext(
			err,
			"failed attempt to wait for object '%s' to exist",
			objectKey,
		)
	}

	clog.Info(ctx, "Test object was successfully uploaded.", "bn", bucketName)

	_, err = c.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(o *s3.Options) {
		o.Credentials = credentials.NewStaticCredentialsProvider(
			string(s3Creds.AccessKeyID),
			string(s3Creds.SecretAccessKey),
			string(s3Creds.SessionToken),
		)
	})
	if err != nil {
		return errors.AddContext(err, "Failed to delete object '%s'", objectKey)
	}

	err = s3.NewObjectNotExistsWaiter(c.client).Wait(
		ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey)},
		time.Minute)
	if err != nil {
		return errors.AddContext(
			err,
			"failed attempt to wait for object %s to be delteted",
			objectKey,
		)
	}

	clog.Infof(ctx, "Test object successfully deleted.")

	return nil
}

// Credentials implements [types.Client].
func (c *clientS3) NewCredentials(
	ctx context.Context,
	permissions types.BucketPermissions,
	duration time.Duration,
) (creds types.Credentials, err error) {
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

	in := sts.AssumeRoleInput{ //nolint: exhaustruct // intended
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
	clog.Info(ctx, "Created credentials successfully.")

	return &S3Credentials{
		AccessKeyID:     secret.RedactedString(*res.Credentials.AccessKeyId),
		SessionToken:    secret.RedactedString(*res.Credentials.SessionToken),
		SecretAccessKey: secret.RedactedString(*res.Credentials.SecretAccessKey),
		Expiration:      *res.Credentials.Expiration,
	}, nil
}
