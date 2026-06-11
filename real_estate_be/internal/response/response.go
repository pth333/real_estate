package response

import "github.com/gofiber/fiber/v2"

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func Success(c *fiber.Ctx, status int, message string, data interface{}, meta interface{}) error {
	return c.Status(status).JSON(APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func OK(c *fiber.Ctx, data interface{}) error {
	return Success(c, fiber.StatusOK, "", data, nil)
}

func Created(c *fiber.Ctx, message string, data interface{}) error {
	return Success(c, fiber.StatusCreated, message, data, nil)
}

func Error(c *fiber.Ctx, status int, message string, err interface{}) error {
	return c.Status(status).JSON(APIResponse{
		Success: false,
		Message: message,
		Error:   err,
	})
}

func BadRequest(c *fiber.Ctx, message string, err interface{}) error {
	return Error(c, fiber.StatusBadRequest, message, err)
}

func Unauthorized(c *fiber.Ctx, message string, err interface{}) error {
	return Error(c, fiber.StatusUnauthorized, message, err)
}

func InternalServerError(c *fiber.Ctx, message string, err interface{}) error {
	return Error(c, fiber.StatusInternalServerError, message, err)
}
