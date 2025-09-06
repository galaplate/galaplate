package dto

import (
	"github.com/galaplate/core/supports"
	"github.com/gofiber/fiber/v2"
)

// AuthRegisterRequest - Generated on 2025-09-06 06:23:44
type AuthRegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (s *AuthRegisterRequest) Validate(c *fiber.Ctx) (u *AuthRegisterRequest, err error) {
	myValidator := &supports.XValidator{}
	if err := c.BodyParser(s); err != nil {
		return nil, err
	}

	if err := myValidator.Validate(s); err != nil {
		return nil, err
	}

	return s, nil
}
