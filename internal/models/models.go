package models

import "time"

type User struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Post struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Content   string    `json:"content"`
	Status    string    `json:"status"` // draft, published, archived
	AuthorID  uint      `json:"author_id"`
	Author    *User     `json:"author,omitempty"` // Pointer + omitempty
	Category  string    `json:"category"`
	Tags      []string  `json:"tags,omitempty"` // omitempty = don't include if empty
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreatePostRequest - DTO for creating a post
// In PHP: This would be a FormRequest class
// In Go: We use separate struct for input validation
type CreatePostRequest struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
	Status   string   `json:"status"`
}

// UpdatePostRequest - DTO for updating a post
//
// GO CONCEPT - POINTER FIELDS FOR OPTIONAL VALUES:
// Title *string means the field is OPTIONAL
// In PHP: You'd use nullable types ?string
// In Go: We use *string so we can distinguish between:
//   - nil (not provided)
//   - "" (explicitly set to empty)
//   - "value" (set to a value)
type UpdatePostRequest struct {
	Title    *string  `json:"title"`
	Content  *string  `json:"content"`
	Category *string  `json:"category"`
	Tags     []string `json:"tags"`
	Status   *string  `json:"status"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type PaginatedResponse[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	TotalPages int   `json:"total_pages"`
}
