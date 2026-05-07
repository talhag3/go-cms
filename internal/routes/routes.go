package routes

import (
	"github.com/talhag3/go-cms/internal/handlers"
	"github.com/talhag3/go-cms/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// Setup configures all routes for the application
// In PHP: This is like routes/api.php
//
// ROUTE STRUCTURE:
// ┌─────────────────────────────────────────────────────────────┐
// │  GET    /api/v1/posts          → GetAll (public)           │
// │  GET    /api/v1/posts/:id      → GetByID (public)          │
// │  GET    /api/v1/posts/slug/:slug → GetBySlug (public)      │
// │  GET    /api/v1/posts/author/:id → GetByAuthor (public)    │
// │  POST   /api/v1/auth/login     → Login (public)            │
// │  GET    /api/v1/users          → GetAll (public)           │
// │  GET    /api/v1/users/:id      → GetByID (public)          │
// │                                                              │
// │  POST   /api/v1/posts          → Create (protected)        │
// │  PUT    /api/v1/posts/:id      → Update (protected)        │
// │  DELETE  /api/v1/posts/:id     → Delete (protected)        │
// │  POST   /api/v1/users          → Create (protected)        │
// │  DELETE  /api/v1/users/:id     → Delete (protected)        │
// └─────────────────────────────────────────────────────────────┘
func Setup(app *fiber.App, postHandler *handlers.PostHandler, userHandler *handlers.UserHandler) {

	// Health check (no auth required)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"message": "CMS API is running",
		})
	})

	// API v1 base group
	// In PHP: Route::prefix('api/v1')->group(function () { ... });
	api := app.Group("/api")

	// ─── PUBLIC ROUTES ───────────────────────────────────────
	public := api.Group("/v1")

	// Auth routes
	auth := public.Group("/auth")
	auth.Post("/login", userHandler.Login)

	// IMPORTANT: Order matters! More specific routes must come first
	// /posts/slug/:slug must be before /posts/:id
	public.Get("/posts/slug/:slug", postHandler.GetBySlug)
	public.Get("/posts/author/:id", postHandler.GetByAuthor)
	public.Get("/posts", postHandler.GetAll)
	public.Get("/posts/:id", postHandler.GetByID)

	// User public routes
	public.Get("/users", userHandler.GetAll)
	public.Get("/users/:id", userHandler.GetByID)

	// ─── PROTECTED ROUTES ────────────────────────────────────
	// In PHP: Route::middleware('auth:api')->group(function () { ... });
	protected := api.Group("/v1")
	protected.Use(middleware.Auth())

	// Post protected routes
	protected.Post("/posts", postHandler.Create)
	protected.Put("/posts/:id", postHandler.Update)
	protected.Delete("/posts/:id", postHandler.Delete)

	// User protected routes
	protected.Post("/users", userHandler.Create)
	protected.Delete("/users/:id", userHandler.Delete)
}
