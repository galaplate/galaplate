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

	// Example routes with different policy combinations
	api := app.Group("/api")

	var logController = controllers.LogControllerInstance
	app.Get("/logs", middleware.BasicAuth(), logController.ShowLogsPage)

	// Test routes for testing framework
	var testController = controllers.TestControllerInstance
	api.Get("/health", testController.GetHealthCheck)
	api.Post("/test", testController.CreateTestData)
	api.Get("/test/:id", testController.GetTestData)
}
