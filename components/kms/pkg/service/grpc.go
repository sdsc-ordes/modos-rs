package service

import (
	kms "github.com/sdsc-ordes/modos-rs/components/kms/pkg/api/v1"
	"google.golang.org/grpc"
)

// RegisterAtGRPCServer registers the service `s` to a GRPC server.
func (s *Service) RegisterAtGRPCServer(server *grpc.Server) {
	kms.RegisterKMSServer(server, s)
}
