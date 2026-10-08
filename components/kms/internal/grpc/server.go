package grpc

import (
	"context"
	"fmt"
	"net"

	"golang.org/x/sync/errgroup"

	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	protoval "buf.build/go/protovalidate"
	protovalMiddlewear "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/auth"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/errors"
	clog "gitlab.com/data-custodian/custodian/components/lib-common/pkg/log/context"
	"google.golang.org/grpc"

	"github.com/sdsc-ordes/modos-rs/components/kms/internal/config"
	"github.com/sdsc-ordes/modos-rs/components/kms/internal/middlewear"
)

// Server represents the GRPC server.
type Server struct {
	// The actual server with public endpoints.
	Public *grpc.Server

	// The health server.
	management *grpc.Server
	ctx        context.Context
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
		grpc.ChainUnaryInterceptor(
			authInterceptor,
			validationInterceptor,
		),
	)

	serverMgmt := grpc.NewServer()
	healthSrv := health.NewServer()
	healthSrv.SetServingStatus("", healthgrpc.HealthCheckResponse_SERVING)
	healthgrpc.RegisterHealthServer(serverMgmt, healthSrv)
	reflection.Register(serverMgmt)

	return &Server{Public: server, management: serverMgmt}, nil
}

// Serve serves all endpoints registered.
// This function returns whenever the `ctx` is canceled.
func (s *Server) Serve(
	ctx context.Context,
	cfg *config.Server,
	cfgMgmt *config.Server,
) error {
	addr := fmt.Sprintf("%v:%v", cfg.Hostname, cfg.Port)
	addrMgmt := fmt.Sprintf("%v:%v", cfgMgmt.Hostname, cfgMgmt.Port)
	clog.Infof(ctx, "Start serving public endpoints at '%s'.", addr)
	clog.Infof(ctx, "Start serving management endpoints at '%s'.", addrMgmt)

	// Create TCP socket.
	lc := net.ListenConfig{} //nolint:exhaustruct // all fields optional
	listener, err := lc.Listen(ctx, "tcp", addr)
	if err != nil {
		return errors.AddContext(err, "could not create listener at '%s'", addr)
	}

	listenerMgmt, err := lc.Listen(ctx, "tcp", addrMgmt)
	if err != nil {
		return errors.AddContext(err, "could not create listener at '%s'", addr)
	}

	s.ctx = ctx

	// NOTE: `grpc.Server.Serve` does not observe `ctx`, it only returns once
	// `Stop`/`GracefulStop` is called. Bridge the cancellation ourselves.
	go func() {
		<-ctx.Done() // Wait till ctx is canceled.
		clog.Info(ctx, "Context cancelled, shutting down GRPC server.")
		s.close()
	}()

	// Spawn all servers and wait on them.
	eg, _ := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return s.Public.Serve(listener)
	})

	eg.Go(func() error {
		return s.management.Serve(listenerMgmt)
	})

	return eg.Wait()
}

// Close cleans up all resources.
func (s *Server) close() {
	clog.Info(s.ctx, "Shutting down GRPC server.")
	s.Public.GracefulStop()
	s.management.GracefulStop()
}
