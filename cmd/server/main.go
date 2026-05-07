// cmd/server/main.go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/talhag3/go-cms/internal/config"
	"github.com/talhag3/go-cms/internal/handlers"
	"github.com/talhag3/go-cms/internal/middleware"
	"github.com/talhag3/go-cms/internal/repositories"
	"github.com/talhag3/go-cms/internal/routes"
	"github.com/talhag3/go-cms/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// ═══════════════════════════════════════════════════════════
	// STEP 1: Load Configuration
	// In PHP: This happens automatically via bootstrap/app.php
	// ═══════════════════════════════════════════════════════════
	cfg := config.Load()

	// ═══════════════════════════════════════════════════════════
	// STEP 2: Initialize Repositories (Data Access Layer)
	// In PHP: These would be resolved by the service container
	// ═══════════════════════════════════════════════════════════
	postRepo := repositories.NewMockPostRepository()
	userRepo := repositories.NewMockUserRepository()

	// ═══════════════════════════════════════════════════════════
	// STEP 3: Initialize Services (Business Logic Layer)
	// Services receive their dependencies via constructor
	// This is MANUAL DEPENDENCY INJECTION
	//
	// In PHP (Laravel):
	//   $postService = app(PostService::class);
	//   // Container automatically injects dependencies
	//
	// In Go:
	//   postService := services.NewPostService(postRepo, userRepo)
	//   // We explicitly pass dependencies
	// ═══════════════════════════════════════════════════════════
	postService := services.NewPostService(postRepo, userRepo)
	userService := services.NewUserService(userRepo)

	// ═══════════════════════════════════════════════════════════
	// STEP 4: Initialize Handlers (HTTP/Controller Layer)
	// ═══════════════════════════════════════════════════════════
	postHandler := handlers.NewPostHandler(postService)
	userHandler := handlers.NewUserHandler(userService)

	// ═══════════════════════════════════════════════════════════
	// STEP 5: Create Fiber Application
	// In PHP: This is like bootstrapping Laravel
	// ═══════════════════════════════════════════════════════════
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorHandler: customErrorHandler,
	})

	// ═══════════════════════════════════════════════════════════
	// STEP 6: Register Global Middleware
	// In PHP: These go in app/Http/Kernel.php $middleware array
	// ═══════════════════════════════════════════════════════════
	app.Use(recover.New())          // Recover from panics (like try-catch)
	app.Use(middleware.Logger())    // Log requests
	app.Use(middleware.RequestID()) // Add request ID
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// ═══════════════════════════════════════════════════════════
	// STEP 7: Setup Routes
	// In PHP: This loads routes/api.php
	// ═══════════════════════════════════════════════════════════
	routes.Setup(app, postHandler, userHandler)

	// ═══════════════════════════════════════════════════════════
	// STEP 8: Start Server
	// In PHP: php artisan serve --port=3000
	// ═══════════════════════════════════════════════════════════
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	log.Println("╔══════════════════════════════════════════════╗")
	log.Printf("║  🚀 %s is running", cfg.App.Name)
	log.Printf("║  📍 URL:  http://%s", addr)
	log.Printf("║  🌍 Env:  %s", cfg.App.Env)
	log.Println("╠══════════════════════════════════════════════╣")
	log.Println("║  Available Endpoints:                        ║")
	log.Println("║  GET  /health              Health Check      ║")
	log.Println("║  POST /api/v1/auth/login   Login            ║")
	log.Println("║  GET  /api/v1/posts        List Posts       ║")
	log.Println("║  GET  /api/v1/posts/:id    Get Post         ║")
	log.Println("║  POST /api/v1/posts        Create Post 🔒   ║")
	log.Println("║  PUT  /api/v1/posts/:id    Update Post 🔒   ║")
	log.Println("║  DEL  /api/v1/posts/:id    Delete Post 🔒   ║")
	log.Println("║  GET  /api/v1/users        List Users       ║")
	log.Println("╚══════════════════════════════════════════════╝")

	if err := app.Listen(addr); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

// customErrorHandler provides consistent error responses
// In PHP: This is like app/Exceptions/Handler.php
func customErrorHandler(c *fiber.Ctx, err error) error {
	// Handle known Fiber errors
	if e, ok := err.(*fiber.Error); ok {
		return c.Status(e.Code).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    e.Code,
				"message": e.Message,
			},
		})
	}

	// Handle unexpected errors
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"code":    500,
			"message": "Internal server error",
		},
	})
}
