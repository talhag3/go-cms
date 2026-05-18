// cmd/server/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/talhag3/go-cms/internal/config"
	"github.com/talhag3/go-cms/internal/database"
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
	// ═══════════════════════════════════════════════════════════════════════
	// STEP 1: Load Configuration
	// ═══════════════════════════════════════════════════════════════════════
	cfg := config.Load()

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 2: Connect to Database
	// ═══════════════════════════════════════════════════════════════════════
	//
	// GO CONCEPT - CONTEXT WITH TIMEOUT:
	// context.WithTimeout creates a context that auto-cancels after duration
	// This prevents hanging forever if database is unreachable
	//
	// PHP: Connection attempt throws exception after default timeout
	// Go: We explicitly control the timeout
	//
	dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer dbCancel() // Always cancel to free resources

	log.Printf("🗄️  Connecting to PostgreSQL at %s:%s/%s...",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)

	pool, err := database.NewPool(dbCtx, cfg.Database)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer database.Close(pool) // Close pool when main() exits

	log.Println("✅ Database connected successfully")

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 3: Run Migrations
	// ═══════════════════════════════════════════════════════════════════════
	if err := database.RunMigrations(dbCtx, pool); err != nil {
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 4: Initialize Repositories (NOW USING POSTGRESQL!)
	// ═══════════════════════════════════════════════════════════════════════
	//
	// BEFORE: postRepo := repositories.NewMockPostRepository()
	// AFTER:  postRepo := repositories.NewPostPgRepository(pool)
	//
	// This is the ONLY change needed to switch from mock to real DB!
	// Services, handlers, routes - ALL stay the same!
	//
	postRepo := repositories.NewPostPgRepository(pool)
	userRepo := repositories.NewUserPgRepository(pool)

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 5: Initialize Services
	// ═══════════════════════════════════════════════════════════════════════
	postService := services.NewPostService(postRepo, userRepo)
	userService := services.NewUserService(userRepo)

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 6: Initialize Handlers
	// ═══════════════════════════════════════════════════════════════════════
	postHandler := handlers.NewPostHandler(postService)
	userHandler := handlers.NewUserHandler(userService)

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 7: Create Fiber Application
	// ═══════════════════════════════════════════════════════════════════════
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorHandler: customErrorHandler,
	})

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 8: Register Middleware
	// ═══════════════════════════════════════════════════════════════════════
	app.Use(recover.New())
	app.Use(middleware.Logger())
	app.Use(middleware.RequestID())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 9: Setup Routes
	// ═══════════════════════════════════════════════════════════════════════
	routes.Setup(app, postHandler, userHandler)

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 10: Graceful Shutdown
	// ═══════════════════════════════════════════════════════════════════════
	//
	// GO CONCEPT - GRACEFUL SHUTDOWN:
	// In PHP: Process killed, connections lost
	// In Go: We listen for SIGINT/SIGTERM and:
	//   1. Stop accepting new requests
	//   2. Wait for in-flight requests to complete
	//   3. Close database connections
	//   4. Exit cleanly
	//
	// This is ESSENTIAL for production!
	//
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh // Block until signal received

		log.Println("🛑 Shutdown signal received...")

		// Shutdown with 5 second timeout
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := app.ShutdownWithContext(shutdownCtx); err != nil {
			log.Printf("❌ Error during shutdown: %v", err)
		}

		log.Println("✅ Server shutdown complete")
	}()

	// ═══════════════════════════════════════════════════════════════════════
	// STEP 11: Start Server
	// ═══════════════════════════════════════════════════════════════════════
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	log.Println("╔══════════════════════════════════════════════════════════╗")
	log.Printf("║  🚀 %s is running", cfg.App.Name)
	log.Printf("║  📍 URL:  http://%s", addr)
	log.Printf("║  🌍 Env:  %s", cfg.App.Env)
	log.Printf("║  🗄️  DB:   PostgreSQL (%s:%s/%s)",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	log.Println("╠══════════════════════════════════════════════════════════╣")
	log.Println("║  Available Endpoints:                                    ║")
	log.Println("║  GET  /health              Health Check                  ║")
	log.Println("║  POST /api/v1/auth/login   Login                        ║")
	log.Println("║  GET  /api/v1/posts        List Posts                   ║")
	log.Println("║  GET  /api/v1/posts/:id    Get Post                     ║")
	log.Println("║  POST /api/v1/posts        Create Post 🔒              ║")
	log.Println("║  PUT  /api/v1/posts/:id    Update Post 🔒              ║")
	log.Println("║  DEL  /api/v1/posts/:id    Delete Post 🔒              ║")
	log.Println("║  GET  /api/v1/users        List Users                   ║")
	log.Println("╚══════════════════════════════════════════════════════════╝")

	if err := app.Listen(addr); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	if e, ok := err.(*fiber.Error); ok {
		return c.Status(e.Code).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    e.Code,
				"message": e.Message,
			},
		})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"code":    500,
			"message": "Internal server error",
		},
	})
}
