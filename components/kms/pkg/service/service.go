package service

import (
	"context"

	kmsv1 "github.com/sdsc-ordes/modos-rs/components/kms/pkg/api/v1"
	"github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage/types"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/auth"
)

type Service struct {
	// We want to implement all functions and get a compile error if proto file drifts.
	kmsv1.UnsafeKmsServer

	Storage     types.Client
	JWTVerifier *auth.JWTVerifier
}

// Assert that service implements the kmsv1.Server
var _ kmsv1.KmsServer = (*Service)(nil)

// GetCredential implements [kmsv1.KmsServer].
func (s *Service) GetCredential(
	context.Context, *kmsv1.GetCredentialRequest) (*kmsv1.GetCredentialResponse, error) {

	return nil, nil
}
