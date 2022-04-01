package server

import (
	"time"

	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/transport/grpc"
	"github.com/google/wire"

	v1 "github.com/go-microservice/moment-service/api/post/v1"
	"github.com/go-microservice/moment-service/internal/service"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGRPCServer)

// NewGRPCServer creates a gRPC server
func NewGRPCServer(cfg *app.ServerConfig, svc *service.PostServiceServer) *grpc.Server {

	grpcServer := grpc.NewServer(
		grpc.Network("tcp"),
		grpc.Address(":9090"),
		grpc.Timeout(3*time.Second),
	)

	// register biz service
	v1.RegisterPostServiceServer(grpcServer, svc)

	return grpcServer
}
