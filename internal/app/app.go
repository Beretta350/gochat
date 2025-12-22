package app

import (
	"context"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"go.uber.org/fx"

	"github.com/Beretta350/gochat/internal/app/auth"
	appfx "github.com/Beretta350/gochat/internal/app/fx"
	"github.com/Beretta350/gochat/internal/app/handler"
	"github.com/Beretta350/gochat/internal/app/middleware"
	"github.com/Beretta350/gochat/internal/app/worker"
	"github.com/Beretta350/gochat/internal/config"
	"github.com/Beretta350/gochat/pkg/logger"
	"github.com/Beretta350/gochat/pkg/postgres"
	"github.com/Beretta350/gochat/pkg/redisclient"
)

// Run starts the application with Fx dependency injection
func Run() {
	fx.New(
		// Provide all dependencies
		appfx.Module,

		// Invoke the server
		fx.Invoke(startServer),
	).Run()
}

// ServerParams holds all dependencies needed to start the server
type ServerParams struct {
	fx.In

	Lifecycle    fx.Lifecycle
	Config       *config.Config
	Postgres     *postgres.Client
	Redis        *redisclient.Client
	JWTService   *auth.JWTService
	Health       *handler.HealthHandler
	Auth         *handler.AuthHandler
	Conversation *handler.ConversationHandler
	WebSocket    *handler.WebSocketHandler
	Worker       *worker.MessageWorker
}

func startServer(p ServerParams) {
	app := fiber.New(fiber.Config{
		AppName:      "GoChat API v1.0",
		ErrorHandler: middleware.CustomErrorHandler,
	})

	// Setup middlewares
	middleware.Setup(app)

	// Setup routes
	setupRoutes(app, p)

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Infof("🚀 Starting GoChat on port %s", p.Config.Server.Port)
			logger.Info("📊 Metrics: /metrics")
			logger.Info("🔐 Auth: /api/v1/auth/*")
			logger.Info("💬 Conversations: /api/v1/conversations/*")
			logger.Info("🔌 WebSocket: /ws?token=<jwt>")
			logger.Info("❤️  Health: /api/v1/health")

			// Start worker in background
			go p.Worker.Start(ctx)

			// Start server in background
			go func() {
				if err := app.Listen(":" + p.Config.Server.Port); err != nil {
					logger.Errorf("Server error: %v", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down...")
			_ = p.Redis.Close()
			p.Postgres.Close()
			return app.Shutdown()
		},
	})
}

func setupRoutes(app *fiber.App, p ServerParams) {
	// API v1 routes
	api := app.Group("/api/v1")

	// Health check (public)
	api.Get("/health", p.Health.Check)

	// Auth routes (public)
	authGroup := api.Group("/auth")
	authGroup.Post("/register", p.Auth.Register)
	authGroup.Post("/login", p.Auth.Login)
	authGroup.Post("/refresh", p.Auth.Refresh)

	// Protected auth routes
	authGroup.Get("/me", middleware.AuthMiddleware(p.JWTService), p.Auth.Me)

	// Conversation routes (protected)
	convGroup := api.Group("/conversations", middleware.AuthMiddleware(p.JWTService))
	convGroup.Post("/", p.Conversation.Create)
	convGroup.Get("/", p.Conversation.List)
	convGroup.Get("/:id", p.Conversation.Get)
	convGroup.Get("/:id/messages", p.Conversation.GetMessages)

	// WebSocket routes (JWT in query string)
	ws := app.Group("/ws")
	ws.Use(p.WebSocket.Upgrade)
	ws.Get("/", websocket.New(p.WebSocket.Handle(context.Background())))

	// Monitoring
	app.Get("/metrics", monitor.New(monitor.Config{
		Title: "GoChat Metrics",
	}))
}
