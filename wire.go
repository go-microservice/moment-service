//go:build wireinject
// +build wireinject

package main

import (
	eagle "github.com/go-eagle/eagle/pkg/app"
	"github.com/go-microservice/moment-service/internal/cache"
	"github.com/go-microservice/moment-service/internal/repository"
	"github.com/go-microservice/moment-service/internal/server"
	"github.com/go-microservice/moment-service/internal/service"
	"github.com/google/wire"
)

func InitApp(cfg *eagle.Config, config *eagle.ServerConfig) (*eagle.App, func(), error) {
	panic(wire.Build(server.ProviderSet, service.ProviderSet, repository.ProviderSet, cache.ProviderSet, newApp))
}
