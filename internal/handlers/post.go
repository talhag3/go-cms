// internal/handlers/post.go
package handlers

import (
	"strconv"

	"github.com/talhag3/go-cms/internal/models"
	"github.com/talhag3/go-cms/internal/services"
	"github.com/talhag3/go-cms/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type PostHandler struct {
	postService *services.PostService
}

func NewPostHandler(postService *services.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// ═══════════════════════════════════════════════════════════════════════════
// GO CONCEPT - FIBER CONTEXT TO STANDARD CONTEXT:
// ═══════════════════════════════════════════════════════════════════════════
//
// Fiber uses its own *fiber.Ctx (different from standard context.Context)
// pgx needs standard context.Context
//
// c.Context() converts Fiber context to standard context
// This allows request cancellation to propagate to database queries!
//
// PHP: No equivalent (each request is isolated)
// Go: Context carries cancellation from HTTP request to DB query
//
// Example: If client disconnects mid-request, the DB query gets cancelled
// This saves resources on your database!
func (h *PostHandler) GetAll(c *fiber.Ctx) error {
	ctx := c.Context() // Convert to standard context

	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	status := c.Query("status", "")

	result, err := h.postService.GetAllPosts(ctx, page, perPage, status)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch posts")
	}

	return response.Paginated(c, result.Data, result.Total, result.Page, result.PerPage, result.TotalPages)
}

func (h *PostHandler) GetByID(c *fiber.Ctx) error {
	ctx := c.Context()

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid post ID")
	}

	post, err := h.postService.GetPostByID(ctx, uint(id))
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "Post not found")
	}

	return response.Success(c, post)
}

func (h *PostHandler) GetBySlug(c *fiber.Ctx) error {
	ctx := c.Context()

	slug := c.Params("slug")
	if slug == "" {
		return response.Error(c, fiber.StatusBadRequest, "Slug is required")
	}

	post, err := h.postService.GetPostBySlug(ctx, slug)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "Post not found")
	}

	return response.Success(c, post)
}

func (h *PostHandler) Create(c *fiber.Ctx) error {
	ctx := c.Context()

	var req models.CreatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errs := validateCreatePostRequest(req); len(errs) > 0 {
		return response.ValidationError(c, errs)
	}

	userID, ok := c.Locals("userID").(uint)
	if !ok {
		userID = 1
	}

	post, err := h.postService.CreatePost(ctx, req, userID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.SuccessWithMessage(c, "Post created successfully", post)
}

func (h *PostHandler) Update(c *fiber.Ctx) error {
	ctx := c.Context()

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

	post, err := h.postService.UpdatePost(ctx, uint(id), req)
	if err != nil {
		if err.Error() == "post not found" {
			return response.Error(c, fiber.StatusNotFound, "Post not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.SuccessWithMessage(c, "Post updated successfully", post)
}

func (h *PostHandler) Delete(c *fiber.Ctx) error {
	ctx := c.Context()

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid post ID")
	}

	if err := h.postService.DeletePost(ctx, uint(id)); err != nil {
		if err.Error() == "post not found" {
			return response.Error(c, fiber.StatusNotFound, "Post not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.SuccessWithMessage(c, "Post deleted successfully", nil)
}

func (h *PostHandler) GetByAuthor(c *fiber.Ctx) error {
	ctx := c.Context()

	authorID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid author ID")
	}

	posts, err := h.postService.GetPostsByAuthor(ctx, uint(authorID))
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch posts")
	}

	return response.Success(c, posts)
}

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
