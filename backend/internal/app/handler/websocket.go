package handler

import (
	"context"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"

	"github.com/Beretta350/gochat/internal/app/auth"
	"github.com/Beretta350/gochat/internal/app/chat"
	"github.com/Beretta350/gochat/pkg/logger"
)

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	chatService *chat.Service
	jwtService  *auth.JWTService
}

// NewWebSocketHandler creates a new WebSocket handler (Fx provider)
func NewWebSocketHandler(chatService *chat.Service, jwtService *auth.JWTService) *WebSocketHandler {
	logger.Info("WebSocket handler initialized")
	return &WebSocketHandler{
		chatService: chatService,
		jwtService:  jwtService,
	}
}

// Upgrade is a middleware that checks if the request is a WebSocket upgrade
func (h *WebSocketHandler) Upgrade(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		// Try to get token from cookie first (preferred)
		token := c.Cookies(AccessTokenCookie)

		// Fallback to query string (for clients that can't use cookies)
		if token == "" {
			token = c.Query("token")
		}

		if token == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Missing token")
		}

		// Validate JWT token
		claims, err := h.jwtService.ValidateAccessToken(token)
		if err != nil {
			if err == auth.ErrExpiredToken {
				return fiber.NewError(fiber.StatusUnauthorized, "Token expired")
			}
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
		}

		// Store user info in locals for the WebSocket handler
		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("username", claims.Username)

		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

// Handle handles WebSocket connections
func (h *WebSocketHandler) Handle(ctx context.Context) func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		// Get user info from locals (set by Upgrade middleware)
		userID := c.Locals("user_id")
		username := c.Locals("username")

		if userID == nil {
			logger.Error("Missing user info in WebSocket connection")
			_ = c.Close()
			return
		}

		userIDStr, ok := userID.(string)
		if !ok || userIDStr == "" {
			logger.Error("Invalid user ID in WebSocket connection")
			_ = c.Close()
			return
		}

		requestID := c.Query("request_id", "unknown")
		logger.Infof("[%s] WebSocket connection: %s (%v)", requestID, userIDStr, username)

		h.chatService.HandleConnection(ctx, c, userIDStr)
	}
}
