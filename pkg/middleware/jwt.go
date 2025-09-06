package middleware

import (
	"errors"
	"strings"
	"time"

	"github.com/galaplate/core/database"
	"github.com/galaplate/core/env"
	"github.com/galaplate/galaplate/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct{}

type JWTClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func NewJWTService() *JWTService {
	return &JWTService{}
}

func (j *JWTService) GenerateToken(userID uint) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := j.getJWTSecret()
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JWTService) ValidateToken(tokenString string) (*JWTClaims, error) {
	secret := j.getJWTSecret()

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (j *JWTService) getJWTSecret() string {
	secret := env.Get("APP_SECRET")
	if secret == "" {
		secret = "your-secret-key"
	}
	return secret
}

func (j *JWTService) AuthMiddleware() fiber.Handler {
	db := database.Connect
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Authorization header required",
			})
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Invalid authorization header format",
			})
		}

		tokenString := authHeader[7:]
		claims, err := j.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Invalid or expired token",
				"error":   err.Error(),
			})
		}

		// Optional: Verify user still exists in database
		var user models.User
		if err := db.First(&user, claims.UserID).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "User not found",
			})
		}

		// Store user information in context for use in handlers
		c.Locals("user", &user)
		c.Locals("user_id", claims.UserID)

		return c.Next()
	}
}

var JWTServiceInstance = NewJWTService()

func JWTAuth() fiber.Handler {
	return JWTServiceInstance.AuthMiddleware()
}
