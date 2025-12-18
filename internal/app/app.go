package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"

	"github.com/Beretta350/gochat/internal/app/chat"
	"github.com/Beretta350/gochat/internal/app/handler"
	"github.com/Beretta350/gochat/internal/app/middleware"
	"github.com/Beretta350/gochat/internal/app/repository"
	"github.com/Beretta350/gochat/internal/app/router"
	"github.com/Beretta350/gochat/internal/app/worker"
	"github.com/Beretta350/gochat/internal/config"
	"github.com/Beretta350/gochat/pkg/logger"
	"github.com/Beretta350/gochat/pkg/redisclient"
)

// Run starts the application
func Run() {
	serverConfig := config.GetServerConfig()

	// Create context that cancels on shutdown signals
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	go handleShutdown(cancel)

	// Initialize services
	messageRepo := repository.NewInMemoryMessageRepository()
	chatService := chat.NewService()

	// Start message worker
	messageWorker := worker.NewMessageWorker(messageRepo, "worker-1")
	go messageWorker.Start(ctx)

	// Create Fiber app with custom error handler
	app := fiber.New(fiber.Config{
		AppName:      "GoChat API v1.0",
		ErrorHandler: middleware.CustomErrorHandler,
	})

	// Setup middlewares (recover, requestid, logger, helmet, cors, rate limiter)
	middleware.Setup(app)

	// Initialize handlers
	healthHandler := handler.NewHealthHandler()
	wsHandler := handler.NewWebSocketHandler(ctx, chatService)

	// Setup routes
	router.Setup(app, &router.Config{
		HealthHandler:    healthHandler,
		WebSocketHandler: wsHandler,
	})

	// Graceful shutdown hook
	app.Hooks().OnShutdown(func() error {
		logger.Info("Shutting down gracefully...")
		cancel()
		return redisclient.Close()
	})

	// Start server
	logger.Infof("🚀 GoChat server starting on port %s", serverConfig.Port)
	logger.Info("📊 Metrics available at /metrics")
	logger.Info("🔌 WebSocket endpoint: /ws?token=<user_token>")
	logger.Info("❤️  Health check: /api/v1/health")

	if err := app.Listen(":" + serverConfig.Port); err != nil {
		logger.Fatal("Failed to start server:", err)
	}
}

func handleShutdown(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	logger.Info("Shutdown signal received")
	cancel()
}
