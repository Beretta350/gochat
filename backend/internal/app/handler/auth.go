package handler

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/Beretta350/gochat/internal/app/auth"
	"github.com/Beretta350/gochat/internal/config"
	"github.com/Beretta350/gochat/pkg/logger"
	"github.com/Beretta350/gochat/pkg/validator"
)

const (
	AccessTokenCookie  = "access_token"
	RefreshTokenCookie = "refresh_token"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *auth.Service
	config      *config.Config
}

// NewAuthHandler creates a new auth handler (Fx provider)
func NewAuthHandler(authService *auth.Service, cfg *config.Config) *AuthHandler {
	logger.Info("Auth handler initialized")
	return &AuthHandler{
		authService: authService,
		config:      cfg,
	}
}

// setTokenCookies sets HttpOnly cookies for access and refresh tokens
func (h *AuthHandler) setTokenCookies(c *fiber.Ctx, tokens *auth.TokenPair) {
	sameSite := fiber.CookieSameSiteLaxMode
	switch h.config.Cookie.SameSite {
	case "Strict":
		sameSite = fiber.CookieSameSiteStrictMode
	case "None":
		sameSite = fiber.CookieSameSiteNoneMode
	}

	// Access token cookie
	c.Cookie(&fiber.Cookie{
		Name:     AccessTokenCookie,
		Value:    tokens.AccessToken,
		Path:     "/",
		Domain:   h.config.Cookie.Domain,
		MaxAge:   int(h.config.JWT.AccessExpiry.Seconds()),
		Secure:   h.config.Cookie.Secure,
		HTTPOnly: true,
		SameSite: sameSite,
	})

	// Refresh token cookie
	c.Cookie(&fiber.Cookie{
		Name:     RefreshTokenCookie,
		Value:    tokens.RefreshToken,
		Path:     "/api/v1/auth", // Only sent to auth endpoints
		Domain:   h.config.Cookie.Domain,
		MaxAge:   int(h.config.JWT.RefreshExpiry.Seconds()),
		Secure:   h.config.Cookie.Secure,
		HTTPOnly: true,
		SameSite: sameSite,
	})
}

// clearTokenCookies clears the auth cookies (for logout)
func (h *AuthHandler) clearTokenCookies(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     AccessTokenCookie,
		Value:    "",
		Path:     "/",
		Domain:   h.config.Cookie.Domain,
		MaxAge:   -1,
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
	})

	c.Cookie(&fiber.Cookie{
		Name:     RefreshTokenCookie,
		Value:    "",
		Path:     "/api/v1/auth",
		Domain:   h.config.Cookie.Domain,
		MaxAge:   -1,
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
	})
}

// Register handles user registration
// POST /api/v1/auth/register
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req auth.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Sanitize inputs
	req.Email = validator.SanitizeString(req.Email)
	req.Username = validator.SanitizeString(req.Username)

	// Validate struct using go-playground/validator
	if validationErrors := validator.Struct(&req); len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Validation failed",
			"errors": validationErrors,
		})
	}

	response, err := h.authService.Register(c.Context(), &req)
	if err != nil {
		if errors.Is(err, auth.ErrEmailAlreadyExists) {
			return fiber.NewError(fiber.StatusConflict, "Email already exists")
		}
		logger.Errorf("Register error: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to register user")
	}

	// Set HttpOnly cookies
	h.setTokenCookies(c, response.Tokens)

	logger.Infof("User registered: %s", response.User.Email)

	// Return user info (tokens are in cookies, not in response body)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user": response.User,
	})
}

// Login handles user login
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req auth.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Sanitize email
	req.Email = validator.SanitizeString(req.Email)

	// Validate struct using go-playground/validator
	if validationErrors := validator.Struct(&req); len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Validation failed",
			"errors": validationErrors,
		})
	}

	response, err := h.authService.Login(c.Context(), &req)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
		}
		logger.Errorf("Login error: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to login")
	}

	// Set HttpOnly cookies
	h.setTokenCookies(c, response.Tokens)

	logger.Infof("User logged in: %s", response.User.Email)

	// Return user info (tokens are in cookies, not in response body)
	return c.JSON(fiber.Map{
		"user": response.User,
	})
}

// Refresh handles token refresh
// POST /api/v1/auth/refresh
func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	// Get refresh token from cookie
	refreshToken := c.Cookies(RefreshTokenCookie)
	if refreshToken == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Missing refresh token")
	}

	req := &auth.RefreshRequest{
		RefreshToken: refreshToken,
	}

	tokens, err := h.authService.Refresh(c.Context(), req)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidToken) || errors.Is(err, auth.ErrExpiredToken) {
			h.clearTokenCookies(c)
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired refresh token")
		}
		if errors.Is(err, auth.ErrUserNotFound) {
			h.clearTokenCookies(c)
			return fiber.NewError(fiber.StatusUnauthorized, "User not found")
		}
		logger.Errorf("Refresh error: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to refresh token")
	}

	// Set new HttpOnly cookies
	h.setTokenCookies(c, tokens)

	logger.Info("Token refreshed successfully")

	return c.JSON(fiber.Map{
		"success": true,
	})
}

// Logout handles user logout
// POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	h.clearTokenCookies(c)
	return c.JSON(fiber.Map{
		"success": true,
	})
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
