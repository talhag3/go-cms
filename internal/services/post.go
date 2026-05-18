package services

import (
	"context"
	"math"

	"github.com/talhag3/go-cms/internal/models"
	"github.com/talhag3/go-cms/internal/repositories"
)

type PostService struct {
	postRepo repositories.PostRepository
	userRepo repositories.UserRepository
}

func NewPostService(postRepo repositories.PostRepository, userRepo repositories.UserRepository) *PostService {
	return &PostService{
		postRepo: postRepo,
		userRepo: userRepo,
	}
}

// NOTE: All methods now accept context.Context as first parameter

func (s *PostService) GetAllPosts(ctx context.Context, page, perPage int, status string) (*models.PaginatedResponse[models.Post], error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	if perPage > 100 {
		perPage = 100
	}

	posts, total, err := s.postRepo.GetAll(ctx, page, perPage, status)
	if err != nil {
		return nil, err
	}

	for i := range posts {
		if posts[i].AuthorID > 0 {
			author, err := s.userRepo.GetByID(ctx, posts[i].AuthorID)
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

func (s *PostService) GetPostByID(ctx context.Context, id uint) (*models.Post, error) {
	post, err := s.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if post.AuthorID > 0 {
		author, err := s.userRepo.GetByID(ctx, post.AuthorID)
		if err == nil {
			post.Author = author
		}
	}

	return post, nil
}

func (s *PostService) GetPostBySlug(ctx context.Context, slug string) (*models.Post, error) {
	post, err := s.postRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	if post.AuthorID > 0 {
		author, err := s.userRepo.GetByID(ctx, post.AuthorID)
		if err == nil {
			post.Author = author
		}
	}

	return post, nil
}

func (s *PostService) CreatePost(ctx context.Context, req models.CreatePostRequest, authorID uint) (*models.Post, error) {
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

	if post.Tags == nil {
		post.Tags = []string{}
	}

	err := s.postRepo.Create(ctx, post)
	if err != nil {
		return nil, err
	}

	author, _ := s.userRepo.GetByID(ctx, authorID)
	post.Author = author

	return post, nil
}

func (s *PostService) UpdatePost(ctx context.Context, id uint, req models.UpdatePostRequest) (*models.Post, error) {
	post, err := s.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

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

	err = s.postRepo.Update(ctx, post)
	if err != nil {
		return nil, err
	}

	if post.AuthorID > 0 {
		author, err := s.userRepo.GetByID(ctx, post.AuthorID)
		if err == nil {
			post.Author = author
		}
	}

	return post, nil
}

func (s *PostService) DeletePost(ctx context.Context, id uint) error {
	_, err := s.postRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return s.postRepo.Delete(ctx, id)
}

func (s *PostService) GetPostsByAuthor(ctx context.Context, authorID uint) ([]models.Post, error) {
	posts, err := s.postRepo.GetByAuthor(ctx, authorID)
	if err != nil {
		return nil, err
	}

	author, err := s.userRepo.GetByID(ctx, authorID)
	for i := range posts {
		if err == nil {
			posts[i].Author = author
		}
	}

	return posts, nil
}
