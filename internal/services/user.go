// internal/services/user.go
package services

import (
	"context"
	"fmt"
	"time"

	"github.com/talhag3/go-cms/internal/models"
	"github.com/talhag3/go-cms/internal/repositories"
)

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.userRepo.GetAll(ctx)
}

func (s *UserService) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.User, error) {
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Role:     "author",
	}

	err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.userRepo.Authenticate(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	token := fmt.Sprintf("mock_token_%d_%d", user.ID, time.Now().Unix())

	return &models.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return s.userRepo.Delete(ctx, id)
}
