package services

import (
	"fmt"
	"time"

	"github.com/talhag3/go-cms/internal/models"
	"github.com/talhag3/go-cms/internal/repositories"
)

// UserService handles business logic for users
type UserService struct {
	userRepo repositories.UserRepository
}

// NewUserService creates a new UserService
func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetAllUsers returns all users
func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.GetAll()
}

// GetUserByID returns a single user
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

// CreateUser creates a new user
func (s *UserService) CreateUser(req models.CreateUserRequest) (*models.User, error) {
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Role:     "author", // Default role for new users
	}

	err := s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and returns a token
func (s *UserService) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.userRepo.Authenticate(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	// Generate mock token
	// PRODUCTION: Use JWT library like golang-jwt
	token := fmt.Sprintf("mock_token_%d_%d", user.ID, time.Now().Unix())

	return &models.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

// DeleteUser removes a user
func (s *UserService) DeleteUser(id uint) error {
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}
	return s.userRepo.Delete(id)
}
