package server

import (
	"time"

	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/transport/grpc"
	"github.com/google/wire"

	cmtv1 "github.com/go-microservice/moment-service/api/comment/v1"
	postv1 "github.com/go-microservice/moment-service/api/post/v1"
	"github.com/go-microservice/moment-service/internal/service"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGRPCServer)

// NewGRPCServer creates a gRPC server
func NewGRPCServer(
	cfg *app.ServerConfig,
	postSvc *service.PostServiceServer,
	commentSvc *service.CommentServiceServer,
) *grpc.Server {

	grpcServer := grpc.NewServer(
		grpc.Network("tcp"),
		grpc.Address(":9090"),
		grpc.Timeout(3*time.Second),
	)

	// register biz service
	postv1.RegisterPostServiceServer(grpcServer, postSvc)
	cmtv1.RegisterCommentServiceServer(grpcServer, commentSvc)

	return grpcServer
}
