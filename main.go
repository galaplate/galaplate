package main

import (
	"os"

	"github.com/galaplate/core/bootstrap"
	"github.com/galaplate/core/config"
	"github.com/galaplate/core/console"
	"github.com/galaplate/core/logger"
	"github.com/galaplate/galaplate/router"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "console" {
		bootstrap.Init()
		kernel := console.NewKernel()

		if err := kernel.Run(os.Args); err != nil {
			logger.Fatal("Console command failed: ", err.Error())
		}
		return
	}

	// Configure bootstrap with our router
	cfg := bootstrap.DefaultConfig()
	cfg.SetupRoutes = router.SetupRouter

	app := bootstrap.App(cfg)
	port := config.Get("APP_PORT")
	if port == "" {
		port = "8080"
	}

	if err := app.Listen(":" + port); err != nil {
		logger.Fatal("Server won't run: ", err.Error())
	}
}
