package service

import (
	"github.com/google/wire"
)

// Svc global var
var Svc Service

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewPostServiceServer, NewCommentServiceServer)

// Service define all service
type Service interface {
	Greeter() IGreeterService
}

// service struct
type service struct {
}

// New init service
func New() Service {
	return &service{}
}

func (s *service) Greeter() IGreeterService {
	return newGreeterService(s)
}
