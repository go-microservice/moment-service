package repository

import (
	"context"
	"time"

	"github.com/go-eagle/eagle/pkg/client/consulclient"
	"github.com/go-eagle/eagle/pkg/registry"
	"github.com/go-eagle/eagle/pkg/registry/consul"
	"github.com/go-eagle/eagle/pkg/transport/grpc"
	"github.com/go-microservice/moment-service/internal/model"
	"github.com/google/wire"

	relationV1 "github.com/go-microservice/relation-service/api/relation/v1"
)

// ProviderSet is repo providers.
var ProviderSet = wire.NewSet(
	model.Init,
	NewPostInfo,
	NewPostLatest,
	NewPostHot,
	NewUserPost,
	NewCommentInfo,
	NewCommentContent,
	NewCommentLatest,
	NewCommentHot,
	NewUserComment,
	NewUserLike,
)

func getConsulDiscovery() registry.Discovery {
	client, err := consulclient.New()
	if err != nil {
		panic(err)
	}
	return consul.New(client)
}

func NewRelationClient() relationV1.RelationServiceClient {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	endpoint := "discovery:///relation-svc"
	conn, err := grpc.DialInsecure(
		ctx,
		grpc.WithEndpoint(endpoint),
		grpc.WithDiscovery(getConsulDiscovery()),
	)
	if err != nil {
		panic(err)
	}
	c := relationV1.NewRelationServiceClient(conn)
	return c
}
