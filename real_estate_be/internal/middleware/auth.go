package middleware

import (
	"time"

	"real_estate_be/internal/response"
	"real_estate_be/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware check access token từ Authorization header
func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return response.Unauthorized(c, "Missing authorization header", nil)
	}

	// Parse "Bearer <token>"
	tokenStr := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenStr = authHeader[7:]
	}

	claims, err := jwt.ParseAccessToken(tokenStr)
	if err != nil {
		return response.Unauthorized(c, "Invalid or expired token", err.Error())
	}

	// Check expiration
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return response.Unauthorized(c, "Token expired", nil)
	}

	// Lưu email vào user context để dùng ở controller
	c.Locals("email", claims.Email)
	c.Locals("token", tokenStr)

	return c.Next()
}
