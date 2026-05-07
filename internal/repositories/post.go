package repositories

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/talhag3/go-cms/internal/models"
)

type PostRepository interface {
	GetAll(page, perPage int, status string) ([]models.Post, int64, error)
	GetByID(id uint) (*models.Post, error)
	GetBySlug(slug string) (*models.Post, error)
	Create(post *models.Post) error
	Update(post *models.Post) error
	Delete(id uint) error
	GetByAuthor(authorID uint) ([]models.Post, error)
}

type MockPostRepository struct {
	mu     sync.RWMutex
	posts  map[uint]models.Post
	nextID uint
}

func NewMockPostRepository() *MockPostRepository {
	repo := &MockPostRepository{
		posts:  make(map[uint]models.Post),
		nextID: 1,
	}
	repo.seedData()
	return repo
}

func (r *MockPostRepository) seedData() {
	now := time.Now()
	samplePosts := []models.Post{
		{
			Title:     "Getting Started with Go",
			Slug:      "getting-started-with-go",
			Content:   "Go is a statically typed, compiled language designed for simplicity and efficiency. It was created at Google by Robert Griesemer, Rob Pike, and Ken Thompson.",
			Status:    "published",
			AuthorID:  1,
			Category:  "Programming",
			Tags:      []string{"go", "tutorial", "beginner"},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Title:     "Understanding Fiber Framework",
			Slug:      "understanding-fiber-framework",
			Content:   "Fiber is an Express.js inspired web framework built on top of Fasthttp, the fastest HTTP engine for Go. Designed for ease of use with zero memory allocation in hot paths.",
			Status:    "published",
			AuthorID:  1,
			Category:  "Programming",
			Tags:      []string{"fiber", "go", "web-framework"},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Title:     "Draft: Advanced Go Patterns",
			Slug:      "advanced-go-patterns",
			Content:   "This post covers advanced design patterns in Go including worker pools, fan-out/fan-in, pipelines, and more.",
			Status:    "draft",
			AuthorID:  2,
			Category:  "Programming",
			Tags:      []string{"go", "patterns", "advanced"},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Title:     "Building REST APIs in Go",
			Slug:      "building-rest-apis-go",
			Content:   "Learn how to build professional REST APIs using Go and Fiber framework. This comprehensive guide covers routing, middleware, validation, and more.",
			Status:    "published",
			AuthorID:  2,
			Category:  "API",
			Tags:      []string{"api", "rest", "go", "fiber"},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Title:     "Archived: Old Tutorial",
			Slug:      "old-tutorial",
			Content:   "This is an old tutorial that has been archived and is no longer maintained.",
			Status:    "archived",
			AuthorID:  1,
			Category:  "Old",
			Tags:      []string{"old"},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	for _, post := range samplePosts {
		post.ID = r.nextID
		r.posts[post.ID] = post
		r.nextID++
	}
}

// GetAll retrieves posts with pagination and optional status filter

func (r *MockPostRepository) GetAll(page, perPage int, status string) ([]models.Post, int64, error) {
	// RLock = Read Lock (allows multiple readers, no writers)
	// defer = Execute this when function returns (like finally{})
	r.mu.RLock()
	defer r.mu.RUnlock()

	// GO CONCEPT - SLICE:
	// var filtered []models.Post = nil slice
	// We use []models.Post{} for empty slice (not nil)
	var filtered []models.Post

	for _, post := range r.posts {
		// Filter by status if provided
		if status != "" && post.Status != status {
			continue
		}
		filtered = append(filtered, post)
	}

	total := int64(len(filtered))

	// Apply pagination defaults
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	// Calculate slice boundaries
	start := (page - 1) * perPage
	if start >= len(filtered) {
		return []models.Post{}, total, nil
	}

	end := start + perPage
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[start:end], total, nil
}

// GetByID retrieves a single post by its ID
func (r *MockPostRepository) GetByID(id uint) (*models.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	post, exists := r.posts[id]
	// GO CONCEPT - COMMA OK IDIOM:
	// exists is a boolean indicating if key was found
	// In PHP: if (isset($posts[$id]))
	if !exists {
		return nil, errors.New("post not found")
	}
	return &post, nil
}

// GetBySlug retrieves a post by its slug
func (r *MockPostRepository) GetBySlug(slug string) (*models.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, post := range r.posts {
		if post.Slug == slug {
			return &post, nil
		}
	}
	return nil, errors.New("post not found")
}

// Create adds a new post to the store
func (r *MockPostRepository) Create(post *models.Post) error {
	// Lock() = Write Lock (exclusive access)
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check for duplicate slug
	for _, p := range r.posts {
		if p.Slug == post.Slug {
			return errors.New("slug already exists")
		}
	}

	post.ID = r.nextID
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	r.posts[post.ID] = *post
	r.nextID++

	return nil
}

// Update modifies an existing post
func (r *MockPostRepository) Update(post *models.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.posts[post.ID]; !exists {
		return errors.New("post not found")
	}

	// Check for duplicate slug (excluding current post)
	for _, p := range r.posts {
		if p.Slug == post.Slug && p.ID != post.ID {
			return errors.New("slug already exists")
		}
	}

	post.UpdatedAt = time.Now()
	r.posts[post.ID] = *post
	return nil
}

// Delete removes a post from the store
func (r *MockPostRepository) Delete(id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.posts[id]; !exists {
		return errors.New("post not found")
	}

	delete(r.posts, id)
	return nil
}

// GetByAuthor retrieves all posts by a specific author
func (r *MockPostRepository) GetByAuthor(authorID uint) ([]models.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var posts []models.Post
	for _, post := range r.posts {
		if post.AuthorID == authorID {
			posts = append(posts, post)
		}
	}
	return posts, nil
}

// GenerateSlug creates a URL-friendly slug from a title

func GenerateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1 // -1 means "remove this character"
	}, slug)

	return slug
}
