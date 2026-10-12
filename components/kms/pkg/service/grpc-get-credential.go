package service

import (
	"context"

	"github.com/sdsc-ordes/modos-rs/components/kms/internal/grpc"
	kms "github.com/sdsc-ordes/modos-rs/components/kms/pkg/api/kms/v1"
	clog "gitlab.com/data-custodian/custodian/components/lib-common/pkg/log/context"
)

// Assert that service implements the kmsv1.Server.
var _ kms.KMSServer = (*Service)(nil)

// GetCredential implements [kms.KmsServer].
func (s *Service) GetCredential(
	ctx context.Context,
	_ *kms.GetCredentialRequest,
) (*kms.GetCredentialResponse, error) {
	clog.Info(ctx, "Get credential start.")

	_ = grpc.Claims(ctx)

	return nil, nil
}
