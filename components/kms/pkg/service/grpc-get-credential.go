package service

import (
	"context"
	"time"

	"github.com/sdsc-ordes/modos-rs/components/kms/internal/grpc"
	kms "github.com/sdsc-ordes/modos-rs/components/kms/pkg/api/kms/v1"
	"github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage/s3"
	clog "gitlab.com/data-custodian/custodian/components/lib-common/pkg/log/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Assert that service implements the kmsv1.Server.
var _ kms.KMSServer = (*Service)(nil)

// GetCredential implements [kms.KMSServer].
func (s *Service) GetCredential(
	ctx context.Context,
	_ *kms.GetCredentialRequest,
) (*kms.GetCredentialResponse, error) {
	clog.Info(ctx, "Get credential start.")

	cl := grpc.Claims(ctx)

	creds, err := s.Storage.NewCredentials(
		ctx,
		cl.BucketPermissions,
		1*time.Hour,
	)
	if err != nil {
		return nil, status.Errorf(
			codes.Unauthenticated,
			"retrieving credential failed")
	}

	switch creds.Type() {
	case s3.CredentialType:
		c, ok := creds.(*s3.S3Credentials)
		if !ok {
			clog.Errorf(ctx, "Could not convert to S3 credential")

			return nil, status.Errorf(
				codes.Internal,
				"storage type '%v' not implemented",
				creds.Type())
		}

		return &kms.GetCredentialResponse{
			AccessKeyId:     (string)(c.AccessKeyID),
			SecretAccessKey: (string)(c.SecretAccessKey),
			SessionToken:    (string)(c.SessionToken),
			Expiration:      timestamppb.New(c.Expiration),
		}, nil

	default:
		return nil, status.Errorf(
			codes.Internal,
			"credential type '%v' not implemented", creds.Type(),
		)
	}
}
