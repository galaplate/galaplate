package controllers

import (
	"errors"
	"fmt"

	"github.com/galaplate/core/database"
	"github.com/galaplate/core/supports"
	"github.com/galaplate/galaplate/pkg/dto"
	"github.com/galaplate/galaplate/pkg/middleware"
	"github.com/galaplate/galaplate/pkg/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthController struct{}

type AuthResponse struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (ac *AuthController) Register(c *fiber.Ctx) error {
	req, err := new(dto.AuthRegisterRequest).Validate(c)
	if err != nil {
		return err
	}

	db := database.Connect

	// Check if user with email already exists
	var existingUser models.User
	if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"message": "User with this email already exists",
		})
	}

	// Check if user with username already exists
	if err := db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"message": "User with this username already exists",
		})
	}

	// Hash password
	hashedPassword, err := new(supports.Bcrypt).HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to process password",
		})
	}

	// Create user
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Status:   false,
	}

	if err := db.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create user",
			"error":   err.Error(),
		})
	}

	// Generate JWT token
	jwtService := middleware.NewJWTService()
	token, err := jwtService.GenerateToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to generate token",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "User registered successfully",
		"data": AuthResponse{
			User:  &user,
			Token: token,
		},
	})
}

func (ac *AuthController) Login(c *fiber.Ctx) error {
	req, err := new(dto.AuthLoginRequest).Validate(c)
	if err != nil {
		return err
	}

	db := database.Connect

	// Find user by email
	var user models.User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Invalid credentials",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": fmt.Sprintf("Database error: %s", err.Error()),
		})
	}

	// Verify password
	bcryptService := new(supports.Bcrypt)
	if !bcryptService.DoPasswordsMatch(user.Password, req.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid credentials",
		})
	}

	// Generate JWT token
	jwtService := middleware.NewJWTService()
	token, err := jwtService.GenerateToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to generate token",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Login successful",
		"data": AuthResponse{
			User:  &user,
			Token: token,
		},
	})
}

var AuthControllerInstance = NewAuthController()
