package middleware

import (
	"fmt"
	"strings"

	"github.com/talhag3/go-cms/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// Auth middleware checks for valid authorization header
//
// GO CONCEPT - HIGHER-ORDER FUNCTION:
// This function RETURNS another function (closure)
// In PHP: This is like a middleware class with handle() method
//
// PHP Version:
//
//	public function handle($request, Closure $next) {
//	    $token = $request->header('Authorization');
//	    // validate...
//	    return $next($request);
//	}
func Auth() fiber.Handler {
	// This returned function is the actual middleware
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		// In PHP: $request->header('Authorization')
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return response.Error(c, fiber.StatusUnauthorized, "Authorization header is required")
		}

		// Split "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid authorization format. Use: Bearer <token>")
		}

		token := parts[1]
		if token == "" {
			return response.Error(c, fiber.StatusUnauthorized, "Token is required")
		}

		// Mock token validation
		// PRODUCTION: Parse and validate JWT
		if !strings.HasPrefix(token, "mock_token_") {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid token")
		}

		// Extract user ID from mock token
		var userID uint
		_, err := fmt.Sscanf(token, "mock_token_%d_", &userID)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid token format")
		}

		// Store user ID in context
		// In PHP: $request->attributes->set('userID', $userID)
		// In Go: c.Locals() stores values in request context
		c.Locals("userID", userID)

		// Continue to next handler
		// In PHP: return $next($request);
		return c.Next()
	}
}
