package server

import (
	"time"

	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/transport/grpc"
	"github.com/google/wire"

	v1 "github.com/go-microservice/moment-service/api/moment/v1"
	"github.com/go-microservice/moment-service/internal/service"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGRPCServer)

// NewGRPCServer creates a gRPC server
func NewGRPCServer(
	cfg *app.ServerConfig,
	postSvc *service.PostServiceServer,
	commentSvc *service.CommentServiceServer,
	likeSvc *service.LikeServiceServer,
) *grpc.Server {

	grpcServer := grpc.NewServer(
		grpc.Network("tcp"),
		grpc.Address(":9090"),
		grpc.Timeout(3*time.Second),
	)

	// register biz service
	v1.RegisterPostServiceServer(grpcServer, postSvc)
	v1.RegisterCommentServiceServer(grpcServer, commentSvc)
	v1.RegisterLikeServiceServer(grpcServer, likeSvc)

	return grpcServer
}
