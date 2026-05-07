package services

import (
	"math"

	"github.com/talhag3/go-cms/internal/models"
	"github.com/talhag3/go-cms/internal/repositories"
)

// PostService handles business logic for posts
//
// GO CONCEPT - DEPENDENCY INJECTION:
// The service receives its dependencies through the constructor
// This makes the code TESTABLE and DECOUPLED
//
// In PHP (Laravel): The service container does this automatically
//
//	public function __construct(PostRepository $postRepo)
//
// In Go: We do it explicitly (manual DI)
type PostService struct {
	postRepo repositories.PostRepository // Interface, not concrete type!
	userRepo repositories.UserRepository
}

// NewPostService creates a new PostService with its dependencies
// This is the CONSTRUCTOR
func NewPostService(postRepo repositories.PostRepository, userRepo repositories.UserRepository) *PostService {
	return &PostService{
		postRepo: postRepo,
		userRepo: userRepo,
	}
}

// GetAllPosts returns paginated posts with enriched author data
func (s *PostService) GetAllPosts(page, perPage int, status string) (*models.PaginatedResponse[models.Post], error) {
	// Set defaults
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	if perPage > 100 {
		perPage = 100
	}

	posts, total, err := s.postRepo.GetAll(page, perPage, status)
	if err != nil {
		return nil, err
	}

	// Enrich posts with author data
	// In PHP: $post->load('author') or eager loading
	for i := range posts {
		if posts[i].AuthorID > 0 {
			author, err := s.userRepo.GetByID(posts[i].AuthorID)
			if err == nil {
				posts[i].Author = author
			}
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &models.PaginatedResponse[models.Post]{
		Data:       posts,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

// GetPostByID returns a single post by ID
func (s *PostService) GetPostByID(id uint) (*models.Post, error) {
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Enrich with author
	if post.AuthorID > 0 {
		author, err := s.userRepo.GetByID(post.AuthorID)
		if err == nil {
			post.Author = author
		}
	}

	return post, nil
}

// GetPostBySlug returns a post by its slug
func (s *PostService) GetPostBySlug(slug string) (*models.Post, error) {
	post, err := s.postRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}

	if post.AuthorID > 0 {
		author, err := s.userRepo.GetByID(post.AuthorID)
		if err == nil {
			post.Author = author
		}
	}

	return post, nil
}

// CreatePost creates a new post with business validation
func (s *PostService) CreatePost(req models.CreatePostRequest, authorID uint) (*models.Post, error) {
	// Business logic: set default status
	status := req.Status
	if status == "" {
		status = "draft"
	}

	post := &models.Post{
		Title:    req.Title,
		Slug:     repositories.GenerateSlug(req.Title),
		Content:  req.Content,
		Status:   status,
		AuthorID: authorID,
		Category: req.Category,
		Tags:     req.Tags,
	}

	// Handle nil tags
	if post.Tags == nil {
		post.Tags = []string{}
	}

	err := s.postRepo.Create(post)
	if err != nil {
		return nil, err
	}

	// Return created post with author
	author, _ := s.userRepo.GetByID(authorID)
	post.Author = author

	return post, nil
}

// UpdatePost updates an existing post (partial update)
//
// GO CONCEPT - DEREFERENCING POINTERS:
// When you have *string and want the string value, use *ptr
// if req.Title != nil { post.Title = *req.Title }
// In PHP: if ($req->title !== null) { $post->title = $req->title; }
func (s *PostService) UpdatePost(id uint, req models.UpdatePostRequest) (*models.Post, error) {
	// Fetch existing post first
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Partial update - only update provided fields
	// In PHP: $post->fill($request->only(['title', 'content']))
	if req.Title != nil {
		post.Title = *req.Title
		post.Slug = repositories.GenerateSlug(*req.Title)
	}
	if req.Content != nil {
		post.Content = *req.Content
	}
	if req.Category != nil {
		post.Category = *req.Category
	}
	if req.Status != nil {
		post.Status = *req.Status
	}
	if req.Tags != nil {
		post.Tags = req.Tags
	}

	err = s.postRepo.Update(post)
	if err != nil {
		return nil, err
	}

	// Enrich with author
	if post.AuthorID > 0 {
		author, err := s.userRepo.GetByID(post.AuthorID)
		if err == nil {
			post.Author = author
		}
	}

	return post, nil
}

// DeletePost removes a post
func (s *PostService) DeletePost(id uint) error {
	_, err := s.postRepo.GetByID(id)
	if err != nil {
		return err
	}
	return s.postRepo.Delete(id)
}

// GetPostsByAuthor returns all posts by a specific author
func (s *PostService) GetPostsByAuthor(authorID uint) ([]models.Post, error) {
	posts, err := s.postRepo.GetByAuthor(authorID)
	if err != nil {
		return nil, err
	}

	author, err := s.userRepo.GetByID(authorID)
	for i := range posts {
		if err == nil {
			posts[i].Author = author
		}
	}

	return posts, nil
}
