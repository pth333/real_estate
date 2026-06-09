package controller

import (
	"real_estate_be/internal/controller/dto"
	"real_estate_be/internal/response"
	"real_estate_be/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service usecase.AuthServiceInterface
}

func NewUserHandler(service usecase.AuthServiceInterface) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	if err := h.service.Register(req); err != nil {
		return response.BadRequest(c, "Register failed", err.Error())
	}

	return response.Created(c, "Dang ky thanh cong", nil)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	token, err := h.service.Login(req)
	if err != nil {
		return response.Unauthorized(c, "Login failed", err.Error())
	}

	return response.OK(c, fiber.Map{
		"token": token,
	})
}
