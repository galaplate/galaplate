package main

import (
	"os"

	"github.com/galaplate/core/bootstrap"
	"github.com/galaplate/core/console"
	"github.com/galaplate/core/env"
	"github.com/galaplate/core/logger"
	pkgConsole "github.com/galaplate/galaplate/console"
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
			logger.Fatal("Console command failed: ", err.Error())
		}
		return
	}

	port := env.Get("APP_PORT")
	if port == "" {
		port = "8080"
	}

	if err := app.Listen(":" + port); err != nil {
		logger.Fatal("Server won't run: ", err.Error())
	}
}
