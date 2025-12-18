package app

import (
	"context"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/Beretta350/gochat/internal/app/chat"
	"github.com/Beretta350/gochat/internal/config"
	applogger "github.com/Beretta350/gochat/pkg/logger"
	"github.com/Beretta350/gochat/pkg/redisclient"
)

func Run() {
	serverConfig := config.GetServerConfig()

	// Create chat service
	chatService := chat.NewService()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "GoChat API v1.0",
	})

	// Middlewares
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "GoChat API is running",
		})
	})

	// WebSocket upgrade middleware
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// WebSocket route
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		userToken := c.Query("token")
		if userToken == "" {
			applogger.Error("Missing user token")
			_ = c.Close()
			return
		}

		// Create context for this connection
		ctx := context.Background()

		// Handle the WebSocket connection
		chatService.HandleConnection(ctx, c, userToken)
	}))

	// Graceful shutdown
	app.Hooks().OnShutdown(func() error {
		applogger.Info("Shutting down...")
		return redisclient.Close()
	})

	// Start server
	applogger.Infof("🚀 Fiber server starting on port %s", serverConfig.Port)
	if err := app.Listen(":" + serverConfig.Port); err != nil {
		applogger.Fatal("Failed to start server:", err)
	}
}
