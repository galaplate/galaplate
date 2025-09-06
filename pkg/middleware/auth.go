package middleware

import (
	"crypto/subtle"
	"encoding/base64"
	"strings"

	config "github.com/galaplate/core/env"
	"github.com/galaplate/core/logger"
	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
}

func (m *AuthMiddleware) BasicAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var auth = c.Get("Authorization")
		if auth == "" {
			c.Set("WWW-Authenticate", `Basic realm="Restricted"`)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized",
			})
		}

		if !strings.HasPrefix(auth, "Basic ") {
			c.Set("WWW-Authenticate", `Basic realm="Restricted"`)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Invalid authorization header",
			})
		}

		payload, err := base64.StdEncoding.DecodeString(auth[6:])
		if err != nil {
			logger.Error("Failed to decode base64 auth:", err)
			c.Set("WWW-Authenticate", `Basic realm="Restricted"`)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Invalid base64 encoding",
			})
		}

		var pair = strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 {
			c.Set("WWW-Authenticate", `Basic realm="Restricted"`)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Invalid credentials format",
			})
		}

		var username = pair[0]
		var password = pair[1]

		var expectedUsername = config.Get("BASIC_AUTH_USERNAME")
		var expectedPassword = config.Get("BASIC_AUTH_PASSWORD")

		if expectedUsername == "" || expectedPassword == "" {
			logger.Error("Basic auth credentials not configured")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Authentication not configured",
			})
		}

		if subtle.ConstantTimeCompare([]byte(username), []byte(expectedUsername)) != 1 ||
			subtle.ConstantTimeCompare([]byte(password), []byte(expectedPassword)) != 1 {
			logger.Warn("Failed basic auth attempt for username:", username)
			c.Set("WWW-Authenticate", `Basic realm="Restricted"`)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Invalid credentials",
			})
		}

		return c.Next()
	}
}

var AuthMiddlewareInstance = &AuthMiddleware{}

func BasicAuth() fiber.Handler {
	return AuthMiddlewareInstance.BasicAuth()
}
