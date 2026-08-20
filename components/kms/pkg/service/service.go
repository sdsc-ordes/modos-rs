package service

import (
	"github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage/types"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/auth"
)

type Service struct {
	Storage     types.Client
	JWTVerifier *auth.JWTVerifier
}
