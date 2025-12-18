package router

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"

	"github.com/Beretta350/gochat/internal/app/handler"
)

// Config holds all handlers needed for routing
type Config struct {
	HealthHandler    *handler.HealthHandler
	WebSocketHandler *handler.WebSocketHandler
}

// Setup configures all routes for the application
func Setup(app *fiber.App, cfg *Config) {
	// API v1 routes
	api := app.Group("/api/v1")
	setupAPIRoutes(api, cfg)

	// WebSocket routes
	ws := app.Group("/ws")
	setupWebSocketRoutes(ws, cfg)

	// Monitoring routes
	setupMonitoringRoutes(app)
}

// setupAPIRoutes configures REST API routes
func setupAPIRoutes(router fiber.Router, cfg *Config) {
	// Health check
	router.Get("/health", cfg.HealthHandler.Check)

	// Future: Add more API routes here
	// router.Get("/users/online", cfg.UserHandler.GetOnline)
	// router.Get("/conversations", cfg.ConversationHandler.List)
}

// setupWebSocketRoutes configures WebSocket routes
func setupWebSocketRoutes(router fiber.Router, cfg *Config) {
	// WebSocket upgrade middleware
	router.Use(cfg.WebSocketHandler.Upgrade)

	// WebSocket connection handler
	router.Get("/", websocket.New(cfg.WebSocketHandler.Handle))
}

// setupMonitoringRoutes configures monitoring and metrics routes
func setupMonitoringRoutes(app *fiber.App) {
	// Monitor dashboard - shows real-time metrics
	app.Get("/metrics", monitor.New(monitor.Config{
		Title: "GoChat Metrics",
	}))
}
