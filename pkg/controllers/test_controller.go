package controllers

import (
	"github.com/galaplate/core/supports"
	"github.com/gofiber/fiber/v2"
)

type TestController struct{}

type CreateTestRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

func (s *CreateTestRequest) Validate(c *fiber.Ctx) (u *CreateTestRequest, err error) {
	if err = supports.NewValidator(c).Validate(s); err != nil {
		return nil, err
	}

	return s, nil

}

func NewTestController() *TestController {
	return &TestController{}
}

func (tc *TestController) GetHealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "API is working",
	})
}

func (tc *TestController) CreateTestData(c *fiber.Ctx) error {
	req, err := new(CreateTestRequest).Validate(c)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Test data created successfully",
		"data": fiber.Map{
			"id":          1,
			"name":        req.Name,
			"description": req.Description,
		},
	})
}

func (tc *TestController) GetTestData(c *fiber.Ctx) error {
	id := c.Params("id")

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"id":          id,
			"name":        "Test Item",
			"description": "This is a test item",
		},
	})
}

var TestControllerInstance = NewTestController()
