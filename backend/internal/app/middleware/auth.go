package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/Beretta350/gochat/internal/app/auth"
)

const (
	AccessTokenCookie = "access_token"
)

// getTokenFromRequest extracts the JWT token from cookie or Authorization header
func getTokenFromRequest(c *fiber.Ctx) string {
	// First, try to get token from cookie (preferred, more secure)
	token := c.Cookies(AccessTokenCookie)
	if token != "" {
		return token
	}

	// Fallback to Authorization header (for backward compatibility / API clients)
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// Check Bearer prefix
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}

// AuthMiddleware creates a JWT authentication middleware
func AuthMiddleware(jwtService *auth.JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := getTokenFromRequest(c)
		if tokenString == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Missing authentication")
		}

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
		tokenString := getTokenFromRequest(c)
		if tokenString == "" {
			return c.Next()
		}

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
