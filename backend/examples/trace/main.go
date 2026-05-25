/**
 *
 *    ____          __
 *   / __/__ ____ _/ /__
 *  / _// _ `/ _ `/ / -_)
 * /___/\_,_/\_, /_/\__/
 *         /___/
 *
 *
 * generate by http://patorjk.com/software/taag/#p=display&f=Small%20Slant&t=Eagle
 */
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/insight/backend/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	_ "go.uber.org/automaxprocs"

	"github.com/insight/backend/internal/model"
	"github.com/insight/backend/internal/repository"
	"github.com/insight/backend/internal/server"
	"github.com/insight/backend/internal/service"
	eagle "github.com/insight/backend/pkg/app"
	logger "github.com/insight/backend/pkg/log"
	"github.com/insight/backend/pkg/trace"
	v "github.com/insight/backend/pkg/version"
)

var (
	cfgDir  = pflag.StringP("config dir", "c", "config", "config path.")
	env     = pflag.StringP("env name", "e", "", "env var name.")
	version = pflag.BoolP("version", "v", false, "show version info.")
)

// @title eagle docs api
// @version 1.0
// @description eagle demo

// @host localhost:8080
// @BasePath /v1
func main() {
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

	// init tracer
	if cfg.EnableTrace {
		var traceCfg trace.Config
		err := config.Load("trace", &traceCfg)
		_, err = trace.InitTracerProvider(traceCfg.ServiceName, traceCfg.CollectorEndpoint)
		if err != nil {
			panic(err)
		}
	}

	// init service
	db, _ := model.GetDB()
	service.Svc = service.New(repository.New(db))

	gin.SetMode(cfg.Mode)

	// start app
	app := eagle.New(
		eagle.WithName(cfg.Name),
		eagle.WithVersion(cfg.Version),
		eagle.WithLogger(logger.GetLogger()),
		eagle.WithServer(
			// init http server
			server.NewHTTPServer(&cfg.HTTP),
		),
	)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
