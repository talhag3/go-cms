// internal/handlers/post.go
package handlers

import (
	"strconv"

	"github.com/talhag3/go-cms/internal/models"
	"github.com/talhag3/go-cms/internal/services"
	"github.com/talhag3/go-cms/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// PostHandler handles HTTP requests for posts
// In PHP: This is like app/Http/Controllers/PostController.php
//
// GO CONCEPT - STRUCT AS CONTAINER:
// We store the service as a struct field
// In PHP: The service is injected via constructor
type PostHandler struct {
	postService *services.PostService
}

// NewPostHandler creates a new PostHandler with dependency injection
func NewPostHandler(postService *services.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// GetAll handles GET /api/v1/posts
// In PHP: public function index(Request $request) { ... }
func (h *PostHandler) GetAll(c *fiber.Ctx) error {
	// Parse query parameters
	// In PHP: $request->query('page', 1)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	status := c.Query("status", "")

	result, err := h.postService.GetAllPosts(page, perPage, status)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch posts")
	}

	return response.Paginated(c, result.Data, result.Total, result.Page, result.PerPage, result.TotalPages)
}

// GetByID handles GET /api/v1/posts/:id
func (h *PostHandler) GetByID(c *fiber.Ctx) error {
	// Parse URL parameter
	// In PHP: $request->route('id')
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid post ID")
	}

	post, err := h.postService.GetPostByID(uint(id))
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "Post not found")
	}

	return response.Success(c, post)
}

// GetBySlug handles GET /api/v1/posts/slug/:slug
func (h *PostHandler) GetBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return response.Error(c, fiber.StatusBadRequest, "Slug is required")
	}

	post, err := h.postService.GetPostBySlug(slug)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "Post not found")
	}

	return response.Success(c, post)
}

// Create handles POST /api/v1/posts
func (h *PostHandler) Create(c *fiber.Ctx) error {
	// Parse request body into struct
	// In PHP: $request->validate([...]) or $request->all()
	var req models.CreatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate input
	if errs := validateCreatePostRequest(req); len(errs) > 0 {
		return response.ValidationError(c, errs)
	}

	// Get user ID from context (set by Auth middleware)
	// In PHP: auth()->id()
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		// Default to user 1 if no auth (for demo purposes)
		userID = 1
	}

	post, err := h.postService.CreatePost(req, userID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.SuccessWithMessage(c, "Post created successfully", post)
}

// Update handles PUT /api/v1/posts/:id
func (h *PostHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid post ID")
	}

	var req models.UpdatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errs := validateUpdatePostRequest(req); len(errs) > 0 {
		return response.ValidationError(c, errs)
	}

	post, err := h.postService.UpdatePost(uint(id), req)
	if err != nil {
		if err.Error() == "post not found" {
			return response.Error(c, fiber.StatusNotFound, "Post not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.SuccessWithMessage(c, "Post updated successfully", post)
}

// Delete handles DELETE /api/v1/posts/:id
func (h *PostHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid post ID")
	}

	if err := h.postService.DeletePost(uint(id)); err != nil {
		if err.Error() == "post not found" {
			return response.Error(c, fiber.StatusNotFound, "Post not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.SuccessWithMessage(c, "Post deleted successfully", nil)
}

// GetByAuthor handles GET /api/v1/posts/author/:id
func (h *PostHandler) GetByAuthor(c *fiber.Ctx) error {
	authorID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid author ID")
	}

	posts, err := h.postService.GetPostsByAuthor(uint(authorID))
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch posts")
	}

	return response.Success(c, posts)
}

// --- Validation Helpers ---
// In production, use a library like go-playground/validator

func validateCreatePostRequest(req models.CreatePostRequest) []string {
	var errs []string

	if req.Title == "" {
		errs = append(errs, "Title is required")
	} else if len(req.Title) < 3 {
		errs = append(errs, "Title must be at least 3 characters")
	} else if len(req.Title) > 200 {
		errs = append(errs, "Title must be at most 200 characters")
	}

	if req.Content == "" {
		errs = append(errs, "Content is required")
	} else if len(req.Content) < 10 {
		errs = append(errs, "Content must be at least 10 characters")
	}

	if req.Category == "" {
		errs = append(errs, "Category is required")
	}

	if req.Status != "" && req.Status != "draft" && req.Status != "published" {
		errs = append(errs, "Status must be 'draft' or 'published'")
	}

	return errs
}

func validateUpdatePostRequest(req models.UpdatePostRequest) []string {
	var errs []string

	if req.Title != nil {
		if len(*req.Title) < 3 {
			errs = append(errs, "Title must be at least 3 characters")
		} else if len(*req.Title) > 200 {
			errs = append(errs, "Title must be at most 200 characters")
		}
	}

	if req.Content != nil && len(*req.Content) < 10 {
		errs = append(errs, "Content must be at least 10 characters")
	}

	if req.Status != nil {
		validStatuses := map[string]bool{"draft": true, "published": true, "archived": true}
		if !validStatuses[*req.Status] {
			errs = append(errs, "Status must be 'draft', 'published', or 'archived'")
		}
	}

	return errs
}
