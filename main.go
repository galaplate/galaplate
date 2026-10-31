package main

import (
	"fmt"
	"os"

	"github.com/galaplate/core/bootstrap"
	"github.com/galaplate/core/config"
	"github.com/galaplate/core/console"
	"github.com/galaplate/core/logger"
	_ "github.com/galaplate/galaplate/db/migrations"
	_ "github.com/galaplate/galaplate/providers"
)

func main() {
	app := bootstrap.NewApp()

	if len(os.Args) > 1 && os.Args[1] == "console" {
		kernel := console.NewKernel()

		if err := kernel.Run(os.Args); err != nil {
			logger.Fatal(fmt.Sprintf("Console command failed: %s", err.Error()))
		}
		return
	}

	port := config.ConfigString("app.port")
	if port == "" {
		port = "8080"
	}

	logger.Info(fmt.Sprintf("Server will started in port %s", port))

	if err := app.Listen(":" + port); err != nil {
		logger.Fatal(fmt.Sprintf("Server won't run: %s", err.Error()))
	}
}
