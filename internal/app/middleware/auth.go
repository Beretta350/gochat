package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/Beretta350/gochat/internal/app/auth"
)

// AuthMiddleware creates a JWT authentication middleware
func AuthMiddleware(jwtService *auth.JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Missing authorization header")
		}

		// Check Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization format")
		}

		tokenString := parts[1]

		// Validate token
		claims, err := jwtService.ValidateAccessToken(tokenString)
		if err != nil {
			if err == auth.ErrExpiredToken {
				return fiber.NewError(fiber.StatusUnauthorized, "Token expired")
			}
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
		}

		// Set user info in context
		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("username", claims.Username)

		return c.Next()
	}
}

// OptionalAuthMiddleware extracts user info if token is present, but doesn't require it
func OptionalAuthMiddleware(jwtService *auth.JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Next()
		}

		tokenString := parts[1]

		claims, err := jwtService.ValidateAccessToken(tokenString)
		if err != nil {
			return c.Next()
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("username", claims.Username)

		return c.Next()
	}
}
