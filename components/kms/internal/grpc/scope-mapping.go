package grpc

import (
	kms "github.com/sdsc-ordes/modos-rs/components/kms/pkg/api/kms/v1"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/errors"
)

// ExpectedScopes maps method names to needed scopes.
func ExpectedScopes(fullMethod string) ([]string, error) {
	switch fullMethod {
	case kms.KMS_GetCredential_FullMethodName:
		return []string{"permissions"}, nil
	default:
		return nil, errors.New("function '%s' not implemented", fullMethod)
	}
}
