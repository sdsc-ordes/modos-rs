package service

import (
	"context"

	kms "github.com/sdsc-ordes/modos-rs/components/kms/pkg/api/v1"
	"github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage/types"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/auth"
)

type Service struct {
	// We want to implement all functions and get a compile error if proto file drifts.
	kms.UnsafeKMSServer

	Storage     types.Client
	JWTVerifier *auth.JWTVerifier
}

// Assert that service implements the kmsv1.Server
var _ kms.KMSServer = (*Service)(nil)

// GetCredential implements [kms.KmsServer].
func (s *Service) GetCredential(
	context.Context, *kms.GetCredentialRequest) (*kms.GetCredentialResponse, error) {

	return nil, nil
}
