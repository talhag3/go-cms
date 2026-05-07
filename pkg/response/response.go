package response

import "github.com/gofiber/fiber/v2"

// Success sends a successful JSON response
// In PHP: return response()->json(['success' => true, 'data' => $data])
//
// GO CONCEPT - INTERFACE{} / ANY:
// interface{} (or 'any' in Go 1.18+) accepts any type
// In PHP: You'd use 'mixed' type hint
func Success(c *fiber.Ctx, data interface{}) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

// SuccessWithMessage sends a success response with a message
func SuccessWithMessage(c *fiber.Ctx, message string, data interface{}) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// Error sends an error JSON response
func Error(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"code":    status,
			"message": message,
		},
	})
}

// ValidationError sends validation error response
func ValidationError(c *fiber.Ctx, errs []string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"code":    400,
			"message": "Validation failed",
			"errors":  errs,
		},
	})
}

// Paginated sends a paginated response
func Paginated(c *fiber.Ctx, data interface{}, total int64, page, perPage, totalPages int) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"items":       data,
			"total":       total,
			"page":        page,
			"per_page":    perPage,
			"total_pages": totalPages,
		},
	})
}
