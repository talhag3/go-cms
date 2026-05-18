// internal/repositories/user_pg.go
package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/talhag3/go-cms/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserPgRepository implements UserRepository using PostgreSQL
type UserPgRepository struct {
	pool *pgxpool.Pool
}

// NewUserPgRepository creates a new PostgreSQL-backed user repository
func NewUserPgRepository(pool *pgxpool.Pool) *UserPgRepository {
	return &UserPgRepository{pool: pool}
}

func (r *UserPgRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User

	query := `
        SELECT id, username, email, role, created_at, updated_at
        FROM users
        WHERE id = $1
    `

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Role,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

func (r *UserPgRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	query := `
        SELECT id, username, email, role, created_at, updated_at
        FROM users
        WHERE email = $1
    `

	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Role,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

func (r *UserPgRepository) GetAll(ctx context.Context) ([]models.User, error) {
	query := `
        SELECT id, username, email, role, created_at, updated_at
        FROM users
        ORDER BY id
    `

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.Role,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	return users, nil
}

func (r *UserPgRepository) Create(ctx context.Context, user *models.User) error {
	query := `
        INSERT INTO users (username, email, password_hash, role, created_at, updated_at)
        VALUES ($1, $2, $3, $4, NOW(), NOW())
        RETURNING id, created_at, updated_at
    `

	// In production, password would be hashed before reaching here
	passwordHash := "mock_hash_" + time.Now().Format("20060102150405")

	err := r.pool.QueryRow(
		ctx,
		query,
		user.Username,
		user.Email,
		passwordHash,
		user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "unique") {
			if strings.Contains(err.Error(), "email") {
				return errors.New("email already exists")
			}
			return errors.New("username already exists")
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserPgRepository) Update(ctx context.Context, user *models.User) error {
	query := `
        UPDATE users
        SET username = $1, email = $2, role = $3, updated_at = NOW()
        WHERE id = $4
        RETURNING updated_at
    `

	err := r.pool.QueryRow(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Role,
		user.ID,
	).Scan(&user.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *UserPgRepository) Delete(ctx context.Context, id uint) error {
	query := `DELETE FROM users WHERE id = $1 RETURNING id`

	var deletedID uint
	err := r.pool.QueryRow(ctx, query, id).Scan(&deletedID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// Authenticate checks credentials and returns user if valid
func (r *UserPgRepository) Authenticate(ctx context.Context, email, password string) (*models.User, error) {
	// First, get the user by email
	user, err := r.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// In production, you'd do:
	// err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	//
	// For mock: accept any non-empty password
	if password == "" {
		return nil, errors.New("password required")
	}

	return user, nil
}
