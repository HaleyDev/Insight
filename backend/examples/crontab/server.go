package main

import (
	eagle "github.com/insight/backend/pkg/app"
	"github.com/insight/backend/pkg/config"
	logger "github.com/insight/backend/pkg/log"
	"github.com/insight/backend/pkg/transport/crontab"
	"github.com/robfig/cron/v3"
)

// cd examples/crontab/
// go run server.go jobs.go
func main() {
	// init config
	c := config.New(".")

	var cfg eagle.Config
	if err := c.Load("app", &cfg); err != nil {
		panic(err)
	}
	// set global
	eagle.Conf = &cfg

	// init jobs
	jobs := map[string]cron.Job{
		"greet":      GreetingJob{"golang"},
		"send_email": SendEmail{" user1"},
	}

	log := logger.Init(logger.WithFilename("crontab"))
	srv := crontab.NewServer(jobs, crontab.Logger{Log: log})

	// start app
	app := eagle.New(
		eagle.WithName(cfg.Name),
		eagle.WithVersion(cfg.Version),
		eagle.WithLogger(log),
		eagle.WithServer(
			srv,
		),
	)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
