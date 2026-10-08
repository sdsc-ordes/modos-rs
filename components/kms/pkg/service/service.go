package service

import (
	kms "github.com/sdsc-ordes/modos-rs/components/kms/pkg/api/v1"
	"github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage/types"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/auth"
	"google.golang.org/grpc"
)

type Service struct {
	// We want to implement all functions and get a compile error if proto file drifts.
	kms.UnsafeKMSServer

	Storage     types.Client
	JWTVerifier *auth.JWTVerifier

	server grpc.Server
}
