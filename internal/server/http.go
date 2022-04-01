package server

import (
	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/transport/http"
	v1 "github.com/go-microservice/moment-service/api/post/v1"
	"github.com/go-microservice/moment-service/internal/routers"
	"github.com/go-microservice/moment-service/internal/service"
)

// NewHTTPServer creates a HTTP server
func NewHTTPServer(c *app.ServerConfig) *http.Server {
	router := routers.NewRouter()

	srv := http.NewServer(
		http.WithAddress(c.Addr),
		http.WithReadTimeout(c.ReadTimeout),
		http.WithWriteTimeout(c.WriteTimeout),
	)

	srv.Handler = router

	v1.RegisterPostServiceHTTPServer(router, &service.PostServiceServer{})

	return srv
}
