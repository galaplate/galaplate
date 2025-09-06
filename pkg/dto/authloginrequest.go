package dto

import (
	"github.com/galaplate/core/supports"
	"github.com/gofiber/fiber/v2"
)

// AuthLoginRequest - Generated on 2025-09-06 06:23:36
type AuthLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (s *AuthLoginRequest) Validate(c *fiber.Ctx) (u *AuthLoginRequest, err error) {
	if err = supports.NewValidator(c).Validate(s); err != nil {
		return nil, err
	}

	return s, nil
}
