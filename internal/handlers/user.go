// internal/handlers/user.go
package handlers

import (
	"strconv"

	"github.com/talhag3/go-cms/internal/models"
	"github.com/talhag3/go-cms/internal/services"
	"github.com/talhag3/go-cms/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetAll handles GET /api/v1/users
func (h *UserHandler) GetAll(c *fiber.Ctx) error {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to fetch users")
	}
	return response.Success(c, users)
}

// GetByID handles GET /api/v1/users/:id
func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "User not found")
	}
	return response.Success(c, user)
}

// Create handles POST /api/v1/users
func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errs := validateCreateUserRequest(req); len(errs) > 0 {
		return response.ValidationError(c, errs)
	}

	user, err := h.userService.CreateUser(req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.SuccessWithMessage(c, "User created successfully", user)
}

// Login handles POST /api/v1/auth/login
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	var errs []string
	if req.Email == "" {
		errs = append(errs, "Email is required")
	}
	if req.Password == "" {
		errs = append(errs, "Password is required")
	}
	if len(errs) > 0 {
		return response.ValidationError(c, errs)
	}

	result, err := h.userService.Login(req)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "Invalid credentials")
	}

	return response.SuccessWithMessage(c, "Login successful", result)
}

// Delete handles DELETE /api/v1/users/:id
func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	if err := h.userService.DeleteUser(uint(id)); err != nil {
		return response.Error(c, fiber.StatusNotFound, "User not found")
	}

	return response.SuccessWithMessage(c, "User deleted successfully", nil)
}

func validateCreateUserRequest(req models.CreateUserRequest) []string {
	var errs []string

	if req.Username == "" {
		errs = append(errs, "Username is required")
	} else if len(req.Username) < 3 {
		errs = append(errs, "Username must be at least 3 characters")
	}

	if req.Email == "" {
		errs = append(errs, "Email is required")
	}

	if req.Password == "" {
		errs = append(errs, "Password is required")
	} else if len(req.Password) < 6 {
		errs = append(errs, "Password must be at least 6 characters")
	}

	return errs
}
