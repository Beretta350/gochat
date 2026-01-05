package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

// Setup configures all middlewares for the application
func Setup(app *fiber.App) {
	// Recover from panics
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	// Request ID - adds unique ID to each request
	app.Use(requestid.New(requestid.Config{
		Header: "X-Request-ID",
	}))

	// Logger - logs all requests
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))

	// Helmet - security headers
	app.Use(helmet.New(helmet.Config{
		XSSProtection:         "0",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSMaxAge:            31536000,
		HSTSExcludeSubdomains: false,
		ContentSecurityPolicy: "default-src 'self'",
		ReferrerPolicy:        "strict-origin-when-cross-origin",
	}))

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Request-ID",
		AllowCredentials: false,
		MaxAge:           86400,
	}))

	// Rate Limiter - 100 requests per minute per IP
	app.Use(limiter.New(limiter.Config{
		Max:               100,
		Expiration:        1 * time.Minute,
		LimiterMiddleware: limiter.SlidingWindow{},
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(ErrorResponse{
				Success: false,
				Error:   "Too Many Requests",
				Message: "Rate limit exceeded. Please try again later.",
			})
		},
	}))
}
