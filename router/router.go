package router

import (
	"github.com/galaplate/galaplate/middleware"
	"github.com/galaplate/galaplate/pkg/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupRouter(app *fiber.App) {

	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello world")
	})

	var logController = controllers.LogControllerInstance
	app.Get("/logs", middleware.BasicAuth(), logController.ShowLogsPage)
}
