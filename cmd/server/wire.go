//go:build wireinject
// +build wireinject

package main

import (
	eagle "github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/client/consulclient"
	logger "github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/registry"
	"github.com/go-eagle/eagle/pkg/registry/consul"
	"github.com/go-eagle/eagle/pkg/transport/grpc"
	"github.com/go-microservice/moment-service/internal/cache"
	"github.com/go-microservice/moment-service/internal/repository"
	"github.com/go-microservice/moment-service/internal/server"
	"github.com/go-microservice/moment-service/internal/service"
	"github.com/google/wire"
)

func InitApp(cfg *eagle.Config, config *eagle.ServerConfig) (*eagle.App, func(), error) {
	panic(wire.Build(server.ProviderSet, service.ProviderSet, repository.ProviderSet, cache.ProviderSet, newApp))
}

func newApp(cfg *eagle.Config, gs *grpc.Server) *eagle.App {
	return eagle.New(
		eagle.WithName(cfg.Name),
		eagle.WithVersion(cfg.Version),
		eagle.WithLogger(logger.GetLogger()),
		eagle.WithServer(
			// init HTTP server
			server.NewHTTPServer(&cfg.HTTP),
			// init gRPC server
			gs,
		),
		//eagle.WithRegistry(getConsulRegistry()),
	)
}

// create a consul register
func getConsulRegistry() registry.Registry {
	client, err := consulclient.New()
	if err != nil {
		panic(err)
	}
	return consul.New(client)
}
