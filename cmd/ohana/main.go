package main

import (
	"fmt"

	"go.uber.org/fx"

	"github.com/OhanaFS/ohana/boundary"
	"github.com/OhanaFS/ohana/config"
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
			config.LoadConfig,
			config.NewLogger,
			boundary.NewRouter,
		),
		fx.Invoke(
			boundary.NewServer,
		),
	).Run()
}
