package handler

import (
	"context"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"

	"github.com/Beretta350/gochat/internal/app/chat"
	"github.com/Beretta350/gochat/pkg/logger"
)

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	chatService *chat.Service
	ctx         context.Context
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(ctx context.Context, chatService *chat.Service) *WebSocketHandler {
	return &WebSocketHandler{
		chatService: chatService,
		ctx:         ctx,
	}
}

// Upgrade is a middleware that checks if the request is a WebSocket upgrade
func (h *WebSocketHandler) Upgrade(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

// Handle handles WebSocket connections
func (h *WebSocketHandler) Handle(c *websocket.Conn) {
	userToken := c.Query("token")
	if userToken == "" {
		logger.Error("Missing user token")
		_ = c.Close()
		return
	}

	requestID := c.Query("request_id", "unknown")
	logger.Infof("[%s] WebSocket connection from user: %s", requestID, userToken)

	h.chatService.HandleConnection(h.ctx, c, userToken)
}
