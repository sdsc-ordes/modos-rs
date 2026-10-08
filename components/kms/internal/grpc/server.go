package grpc

import (
	"context"
	"fmt"
	"net"

	protoval "buf.build/go/protovalidate"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	protovalMiddlewear "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/auth"
	clog "gitlab.com/data-custodian/custodian/components/lib-common/pkg/log/context"
	"google.golang.org/grpc"

	"github.com/sdsc-ordes/modos-rs/components/kms/internal/config"
	"github.com/sdsc-ordes/modos-rs/components/kms/internal/middlewear"
	"github.com/sdsc-ordes/quitsh/pkg/errors"
)

// Server represents the GRPC server.
type Server struct {
	S        *grpc.Server
	listener net.Listener
	ctx      context.Context
}

// NewServer returns a new GRPC server listening.
func NewServer(verifier *auth.JWTVerifier, cfgOIDC *config.OIDC) (*Server, error) {
	protoValidator, err := protoval.New()
	if err != nil {
		return nil, errors.AddContext(err,
			"failed to create protovalidate interceptor")
	}
	validationInterceptor := protovalMiddlewear.UnaryServerInterceptor(protoValidator)

	// Create auth interceptor
	authInterceptor := middlewear.AuthenticationInterceptor(verifier, cfgOIDC)

	// Create gRPC server with middleware chain (matching main.go exactly)
	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				authInterceptor,
				validationInterceptor),
		),
	)

	return &Server{S: server}, nil
}

// Serve serves all endpoints registered.
// This function returns whenever the `ctx` is canceled.
func (s *Server) Serve(ctx context.Context, cfg *config.Server) error {
	addr := fmt.Sprintf("%v:%v", cfg.Hostname, cfg.Port)
	clog.Infof(ctx, "Start serving at '%s'.", addr)

	// Create TCP socket.
	lc := net.ListenConfig{} //nolint:exhaustruct // all fields optional
	listener, err := lc.Listen(ctx, "tcp", addr)
	if err != nil {
		return errors.AddContext(err, "could not create listener at '%s'", addr)
	}

	s.listener = listener
	s.ctx = ctx

	// NOTE: `grpc.Server.Serve` does not observe `ctx`, it only returns once
	// `Stop`/`GracefulStop` is called. Bridge the cancellation ourselves.
	go func() {
		<-ctx.Done()
		clog.Info(ctx, "Context cancelled, shutting down GRPC server.")
		s.close()
	}()

	return s.S.Serve(listener)
}

// Close cleans up all resources.
func (s *Server) close() {
	clog.Info(s.ctx, "Shutting down GRPC server.")
	s.S.GracefulStop()
}
