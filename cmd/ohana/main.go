package main

import (
	"fmt"
	"github.com/OhanaFS/ohana/controller/inc"
	"github.com/OhanaFS/ohana/controller/middleware"
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/selfsign"
	"go.uber.org/fx"
	"os"

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

	// Run normally or generate certs

	flagsConfig := config.LoadFlagsConfig()

	if *flagsConfig.GenCA || *flagsConfig.GenCerts {
		err := selfsign.ProcessFlags(flagsConfig, false)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		os.Exit(0)
	}

	fx.New(
		fx.Provide(
			// Shared providers
			config.LoadConfig,
			config.NewLogger,
			config.NewDatabase,
			middleware.Provide,
			controller.NewRouter,
			inc.NewInc,

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
			controller.NewBackend,

			// DB
			dbfs.InitDB,

			// Inc
			inc.RegisterIncServices,
		),
	).Run()
}
