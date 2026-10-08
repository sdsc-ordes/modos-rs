package service

import (
	"context"

	kms "github.com/sdsc-ordes/modos-rs/components/kms/pkg/api/kms/v1"
)

// Assert that service implements the kmsv1.Server.
var _ kms.KMSServer = (*Service)(nil)

// GetCredential implements [kms.KmsServer].
func (s *Service) GetCredential(
	context.Context, *kms.GetCredentialRequest) (*kms.GetCredentialResponse, error) {
	return nil, nil
}
