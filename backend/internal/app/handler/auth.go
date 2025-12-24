package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/Beretta350/gochat/internal/app/auth"
	"github.com/Beretta350/gochat/pkg/logger"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *auth.Service
}

// NewAuthHandler creates a new auth handler (Fx provider)
func NewAuthHandler(authService *auth.Service) *AuthHandler {
	logger.Info("Auth handler initialized")
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
// POST /api/v1/auth/register
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req auth.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Basic validation
	if req.Email == "" || req.Username == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Email, username and password are required")
	}

	if len(req.Password) < 8 {
		return fiber.NewError(fiber.StatusBadRequest, "Password must be at least 8 characters")
	}

	response, err := h.authService.Register(c.Context(), &req)
	if err != nil {
		if errors.Is(err, auth.ErrEmailAlreadyExists) {
			return fiber.NewError(fiber.StatusConflict, "Email already exists")
		}
		logger.Errorf("Register error: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to register user")
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// Login handles user login
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req auth.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if req.Email == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Email and password are required")
	}

	response, err := h.authService.Login(c.Context(), &req)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
		}
		logger.Errorf("Login error: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to login")
	}

	return c.JSON(response)
}

// Refresh handles token refresh
// POST /api/v1/auth/refresh
func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req auth.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if req.RefreshToken == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Refresh token is required")
	}

	tokens, err := h.authService.Refresh(c.Context(), &req)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidToken) || errors.Is(err, auth.ErrExpiredToken) {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired refresh token")
		}
		if errors.Is(err, auth.ErrUserNotFound) {
			return fiber.NewError(fiber.StatusUnauthorized, "User not found")
		}
		logger.Errorf("Refresh error: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to refresh token")
	}

	return c.JSON(tokens)
}

// Me returns the current user info
// GET /api/v1/auth/me
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	// User info is set by auth middleware
	userID := c.Locals("user_id")
	email := c.Locals("email")
	username := c.Locals("username")

	if userID == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Not authenticated")
	}

	return c.JSON(fiber.Map{
		"id":       userID,
		"email":    email,
		"username": username,
	})
}
