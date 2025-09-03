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

func main() {
	if len(os.Args) > 1 && os.Args[1] == "console" {
		bootstrap.Init()
		kernel := console.NewKernel()
		pkgConsole.RegisterCommands(kernel)

		if err := kernel.Run(os.Args); err != nil {
			logger.Fatal("Console command failed: ", err.Error())
		}
		return
	}

	cfg := bootstrap.DefaultConfig()
	cfg.SetupRoutes = router.SetupRouter
	// Optional: Customize GORM configuration
	// cfg.DatabaseConfig = &bootstrap.DatabaseConfig{
	//     GormConfig: &bootstrap.GormConfig{
	//         Config: gorm.Config{
	//             PrepareStmt: true,
	//             DisableForeignKeyConstraintWhenMigrating: true,
	//         },
	//     },
	// }

	app := bootstrap.App(cfg)
	port := env.Get("APP_PORT")
	if port == "" {
		port = "8080"
	}

	if err := app.Listen(":" + port); err != nil {
		logger.Fatal("Server won't run: ", err.Error())
	}
}
