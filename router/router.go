package router

import (
	"github.com/galaplate/galaplate/pkg/controllers"
	"github.com/galaplate/galaplate/pkg/middleware"
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

	logViewer := app.Group("/admin/logs")
	var logController = controllers.LogController{}
	logViewer.Get("/", logController.Index)
	logViewer.Get("/export", logController.Export)
	logViewer.Post("/cleanup", logController.CleanupLogs)
	logViewer.Get("/stats", logController.GetLogStats)

	// Auth routes
	var authController = controllers.AuthControllerInstance
	api.Post("/register", authController.Register)
	api.Post("/login", authController.Login)

	// Test routes for testing framework
	var testController = controllers.TestControllerInstance
	api.Get("/health", testController.GetHealthCheck)
	api.Post("/test", testController.CreateTestData)
	api.Get("/test/:id", testController.GetTestData)

	// Protected routes (require JWT authentication)
	api.Get("/profile", middleware.JWTAuth(), func(c *fiber.Ctx) error {
		user := c.Locals("user")
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Profile data",
			"data":    user,
		})
	})
}
