package main

import (
	"fmt"
	"github.com/OhanaFS/ohana/dbfs"

	"go.uber.org/fx"

	"github.com/OhanaFS/ohana/config"
	"github.com/OhanaFS/ohana/controller"
	"github.com/OhanaFS/ohana/service"
)

var (
	Version   = "0.0.1"
	BuildTime string
	GitCommit string
)

func main() {
	fmt.Printf("Ohana v%s (built %s, commit %s)\n", Version, BuildTime, GitCommit)

	fx.New(
		fx.Provide(
			// Shared providers
			config.LoadConfig,
			config.NewLogger,
			config.NewDatabase,
			config.NewRedis,
			controller.NewRouter,

			// Services
			service.NewHealth,
		),
		fx.Invoke(
			// HTTP Server
			controller.NewServer,

			// Register routes
			controller.RegisterHealth,

			// DB
			dbfs.InitDB,
			dbfs.RedisTest,
		),
	).Run()
}
