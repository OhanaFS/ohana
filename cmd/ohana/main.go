package main

import (
	"fmt"

	"github.com/OhanaFS/ohana/controller/middleware"
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
			middleware.Provide,
			controller.NewRouter,

			// Services
			service.NewHealth,
			service.NewSession,
			service.NewAuth,
			service.NewUploadService,
		),
		fx.Invoke(
			// HTTP Server
			controller.NewServer,

			// Register routes
			controller.RegisterHealth,
			controller.RegisterAuth,
			controller.RegisterUpload,

			// DB
			dbfs.InitDB,
		),
	).Run()
}
