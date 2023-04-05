package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	eagle "github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/config"
	logger "github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/redis"
	v "github.com/go-eagle/eagle/pkg/version"
	"github.com/go-microservice/moment-service/internal/model"
	"github.com/go-microservice/moment-service/internal/tasks"
	"github.com/spf13/pflag"

	"github.com/hibiken/asynq"
)

var (
	cfgDir  = pflag.StringP("config dir", "c", "config", "config path.")
	env     = pflag.StringP("env name", "e", "", "env var name.")
	version = pflag.BoolP("version", "v", false, "show version info.")
)

func init() {
	pflag.Parse()
	if *version {
		ver := v.Get()
		marshaled, err := json.MarshalIndent(&ver, "", "  ")
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(marshaled))
		return
	}

	// init config
	c := config.New(*cfgDir, config.WithEnv(*env))
	var cfg eagle.Config
	if err := c.Load("app", &cfg); err != nil {
		panic(err)
	}
	// set global
	eagle.Conf = &cfg

	// -------------- init resource -------------
	logger.Init()
	// init db
	model.Init()
	// init redis
	redis.Init()
}

func main() {
	// load config
	c := config.New(*cfgDir, config.WithEnv(*env))
	var cfg tasks.Config
	if err := c.Load("cron", &cfg); err != nil {
		panic(err)
	}

	// -------------- Run worker server ------------
	go func() {
		srv := asynq.NewServer(
			asynq.RedisClientOpt{Addr: cfg.Addr},
			asynq.Config{
				// Specify how many concurrent workers to use
				Concurrency: cfg.Concurrency,
				// Optionally specify multiple queues with different priority.
				Queues: map[string]int{
					tasks.QueueCritical: 6,
					tasks.QueueDefault:  3,
					tasks.QueueLow:      1,
				},
				// See the godoc for other configuration options
			},
		)

		// mux maps a type to a handler
		mux := asynq.NewServeMux()
		// register handlers...
		mux.HandleFunc(tasks.TypePublishPost, tasks.HandlePublishPostTask)
		mux.HandleFunc(tasks.TypeDispatchPost, tasks.HandleDispatchPostTask)

		if err := srv.Run(mux); err != nil {
			log.Fatalf("could not run server: %v", err)
		}
	}()

	// ------------- Run schedule server ------------
	scheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{Addr: cfg.Addr},
		&asynq.SchedulerOpts{Location: time.Local},
	)

	// Register crontab task...
	t, _ := tasks.NewEmailWelcomeTask(5)
	if _, err := scheduler.Register("@every 5s", t); err != nil {
		log.Fatal(err)
	}

	// Run blocks and waits for os signal to terminate the program.
	if err := scheduler.Run(); err != nil {
		log.Fatal(err)
	}
}
