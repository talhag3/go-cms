package repositories

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/talhag3/go-cms/internal/models"
)

// UserRepository defines user data operations
type UserRepository interface {
	GetByID(ctx context.Context, id uint) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
	Authenticate(ctx context.Context, email, password string) (*models.User, error)
}

// MockUserRepository implements UserRepository with in-memory storage
type MockUserRepository struct {
	mu     sync.RWMutex
	users  map[uint]models.User
	nextID uint
}

// NewMockUserRepository creates a new mock user repository
func NewMockUserRepository() *MockUserRepository {
	repo := &MockUserRepository{
		users:  make(map[uint]models.User),
		nextID: 1,
	}
	repo.seedData()
	return repo
}

func (r *MockUserRepository) seedData() {
	now := time.Now()
	users := []models.User{
		{
			Username:  "admin",
			Email:     "admin@cms.com",
			Role:      "admin",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Username:  "editor",
			Email:     "editor@cms.com",
			Role:      "editor",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Username:  "author",
			Email:     "author@cms.com",
			Role:      "author",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	for _, user := range users {
		user.ID = r.nextID
		r.users[user.ID] = user
		r.nextID++
	}
}

func (r *MockUserRepository) GetByID(id uint) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (r *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return &user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *MockUserRepository) GetAll() ([]models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Pre-allocate slice for better performance
	users := make([]models.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}

func (r *MockUserRepository) Create(user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check for duplicate email
	for _, u := range r.users {
		if u.Email == user.Email {
			return errors.New("email already exists")
		}
	}

	user.ID = r.nextID
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	r.users[user.ID] = *user
	r.nextID++
	return nil
}

func (r *MockUserRepository) Update(user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return errors.New("user not found")
	}

	user.UpdatedAt = time.Now()
	r.users[user.ID] = *user
	return nil
}

func (r *MockUserRepository) Delete(id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return errors.New("user not found")
	}

	delete(r.users, id)
	return nil
}

// Authenticate checks credentials and returns user if valid
// In real app, you'd hash and compare passwords with bcrypt
func (r *MockUserRepository) Authenticate(email, password string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Mock authentication - accepts any non-empty password
	// PRODUCTION: Use bcrypt.CompareHashAndPassword()
	for _, user := range r.users {
		if user.Email == email {
			if password == "" {
				return nil, errors.New("password required")
			}
			return &user, nil
		}
	}
	return nil, errors.New("invalid credentials")
}
