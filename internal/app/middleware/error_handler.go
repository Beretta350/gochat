package middleware

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Beretta350/gochat/pkg/logger"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Success   bool   `json:"success"`
	Error     string `json:"error"`
	Message   string `json:"message,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

// CustomErrorHandler handles errors in a standardized way
func CustomErrorHandler(c *fiber.Ctx, err error) error {
	// Default to 500 Internal Server Error
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	// Check if it's a Fiber error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	// Get request ID from locals
	requestID, _ := c.Locals("requestid").(string)

	// Log the error
	logger.Errorf("[%s] Error %d: %s", requestID, code, err.Error())

	// Return JSON error response
	return c.Status(code).JSON(ErrorResponse{
		Success:   false,
		Error:     message,
		Message:   err.Error(),
		RequestID: requestID,
	})
}
