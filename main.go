package main

import (
	"fmt"
	"os"

	"github.com/galaplate/core/bootstrap"
	"github.com/galaplate/core/console"
	"github.com/galaplate/core/env"
	"github.com/galaplate/core/logger"
	pkgConsole "github.com/galaplate/galaplate/console"
	_ "github.com/galaplate/galaplate/db/migrations"
	"github.com/galaplate/galaplate/router"
)

func withSetupRoutes(ac *bootstrap.AppConfig) {
	ac.SetupRoutes = router.SetupRouter
}

func main() {
	app := bootstrap.NewApp(withSetupRoutes)

	if len(os.Args) > 1 && os.Args[1] == "console" {
		// bootstrap.Init(cfg)
		kernel := console.NewKernel()
		pkgConsole.RegisterCommands(kernel)

		if err := kernel.Run(os.Args); err != nil {
			logger.Fatal(fmt.Sprintf("Console command failed: %s", err.Error()))
		}
		return
	}

	port := env.Get("APP_PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info(fmt.Sprintf("Server will started in port %s", port))

	if err := app.Listen(":" + port); err != nil {
		logger.Fatal(fmt.Sprintf("Server won't run: %s", err.Error()))
	}
}
