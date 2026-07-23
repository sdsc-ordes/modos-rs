package s3

import (
	"context"

	bst "github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage/types"
	clog "gitlab.com/data-custodian/custodian/components/lib-common/pkg/log/context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type (
	clientS3 struct {
		client *s3.Client

		// FIXME: Do we need this. Probably we should gather all buckets on startup.
		// or from time to time.
		buckets []string
	}
)

func NewClient(
	ctx context.Context,
	conf *bst.S3Connection,
) (*clientS3, error) {
	clog.Info(ctx, "Create connection to blob storage.",
		"endpoint", conf.Endpoint.String())

	cfg, err := awsconfig.LoadDefaultConfig(
		ctx,
		awsconfig.WithRegion(conf.Region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				string(conf.AccessKey),
				string(conf.SecretKey),
				""),
		),
	)
	if err != nil {
		clog.ErrorE(ctx, err, "Failed to load S3 config.")
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(conf.Endpoint.String())
		o.UsePathStyle = conf.UsePathStyle

		// We upload only encrypted data with `age` which provides integrity.
		o.ResponseChecksumValidation = aws.ResponseChecksumValidationUnset
	})

	return &clientS3{client, nil}, nil
}
