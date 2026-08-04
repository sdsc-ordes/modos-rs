package s3

import (
	"context"

	bst "github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage/types"
	clog "gitlab.com/data-custodian/custodian/components/lib-common/pkg/log/context"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type (
	clientS3 struct {
		client *s3.Client
		sts    *sts.Client
	}
)

func NewClient(
	ctx context.Context,
	conf *bst.S3Connection,
) (*clientS3, error) {
	clog.Info(ctx, "Create connection to S3 storage.",
		"endpoint", conf.Endpoint.String())

	cfg, err := awsconfig.LoadDefaultConfig(
		ctx,
		awsconfig.WithBaseEndpoint(conf.Endpoint.String()),
		awsconfig.WithRegion(conf.Region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				string(conf.AccessKey),
				string(conf.SecretKey),
				""),
		),
	)

	addOpts := func(o *s3.Options) {
		o.UsePathStyle = conf.UsePathStyle
	}

	if err != nil {
		clog.ErrorE(ctx, err, "Failed to load S3 config.")
	}

	client := s3.NewFromConfig(cfg, addOpts)
	sts := sts.NewFromConfig(cfg)

	return &clientS3{client, sts}, nil
}
