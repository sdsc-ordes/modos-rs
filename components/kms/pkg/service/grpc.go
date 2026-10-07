package service

import (
	"context"
	"fmt"
	"net"

	protoval "buf.build/go/protovalidate"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	protovalMiddlewear "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"

	"github.com/sdsc-ordes/modos-rs/components/kms/internal/config"
	"github.com/sdsc-ordes/modos-rs/components/kms/internal/middlewear"
	kms "github.com/sdsc-ordes/modos-rs/components/kms/pkg/api/v1"
	"github.com/sdsc-ordes/quitsh/pkg/errors"
)

// Assert that service implements the kmsv1.Server
var _ kms.KMSServer = (*Service)(nil)

// GetCredential implements [kms.KmsServer].
func (s *Service) GetCredential(
	context.Context, *kms.GetCredentialRequest) (*kms.GetCredentialResponse, error) {

	return nil, nil
}

func (s *Service) Serve(ctx context.Context, cfg *config.Server, cfgOIDC *config.OIDC) error {
	// Create TCP socket.
	lc := net.ListenConfig{} //nolint:exhaustruct // all fields optional
	lis, err := lc.Listen(ctx, "tcp", fmt.Sprintf("%v:%v", cfg.Hostname, cfg.Port))
	if err != nil {
		return errors.AddContext(err, "could not create listener")
	}

	protoValidator, err := protoval.New()
	if err != nil {
		return errors.AddContext(err, "failed to create protovalidate interceptor", err)
	}
	validationInterceptor := protovalMiddlewear.UnaryServerInterceptor(protoValidator)

	// Create auth interceptor
	authInterceptor := middlewear.AuthenticationInterceptor(s.JWTVerifier, cfgOIDC)

	// Create gRPC server with middleware chain (matching main.go exactly)
	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(authInterceptor, validationInterceptor),
		),
	)
}
