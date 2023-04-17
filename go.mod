module github.com/go-microservice/moment-service

go 1.16

require (
	github.com/gin-gonic/gin v1.7.3
	github.com/go-eagle/eagle v1.4.1-0.20220530122715-ee4d43d37e53
	github.com/go-microservice/relation-service v0.0.0-20220615144835-aa0a69fbee93
	github.com/go-redis/redis/v8 v8.11.4
	github.com/google/wire v0.5.0
	github.com/hibiken/asynq v0.22.0
	github.com/jinzhu/copier v0.3.5
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/spf13/cast v1.4.1
	github.com/spf13/pflag v1.0.5
	github.com/swaggo/gin-swagger v1.2.0
	go.opentelemetry.io/otel v1.3.0
	go.opentelemetry.io/otel/trace v1.3.0
	go.uber.org/automaxprocs v1.4.0
	google.golang.org/grpc v1.43.0
	google.golang.org/protobuf v1.27.1
	gorm.io/gorm v1.22.4
)
