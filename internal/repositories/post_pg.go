// internal/repositories/post_pg.go
package repositories

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/talhag3/go-cms/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostPgRepository implements PostRepository using PostgreSQL
//
// ═══════════════════════════════════════════════════════════════════════════
// GO CONCEPT - IMPLICIT INTERFACE IMPLEMENTATION:
// ═══════════════════════════════════════════════════════════════════════════
//
// Remember PostRepository interface from post.go?
//
//	type PostRepository interface {
//	    GetAll(page, perPage int, status string) ([]models.Post, int64, error)
//	    GetByID(id uint) (*models.Post, error)
//	    ...
//	}
//
// This struct AUTOMATICALLY implements PostRepository because it has
// ALL the same methods. No "implements" keyword needed!
//
// PHP: class PostPgRepository implements PostRepository { }
// Go:  (just define the methods - compiler checks automatically)
//
// This is why interfaces are so powerful in Go - you can swap
// implementations without changing ANY other code.
// ═══════════════════════════════════════════════════════════════════════════
type PostPgRepository struct {
	pool *pgxpool.Pool
}

// NewPostPgRepository creates a new PostgreSQL-backed post repository
//
// PHP COMPARISON:
// In Laravel with dependency injection:
//
//	public function __construct(protected PDO $db) {}
//
// In Go:
//
//	func NewPostPgRepository(pool *pgxpool.Pool) *PostPgRepository {
//	    return &PostPgRepository{pool: pool}
//	}
func NewPostPgRepository(pool *pgxpool.Pool) *PostPgRepository {
	return &PostPgRepository{pool: pool}
}

// GetAll retrieves posts with pagination and optional status filter
//
// ═══════════════════════════════════════════════════════════════════════════
// PGX QUERY PATTERNS EXPLAINED:
// ═══════════════════════════════════════════════════════════════════════════
//
// 1. pool.Query(ctx, sql, args...)
//   - Returns multiple rows
//   - MUST close with defer rows.Close()
//   - PHP equivalent: $stmt = $pdo->query($sql); while ($row = $stmt->fetch())
//
// 2. pool.QueryRow(ctx, sql, args...)
//   - Returns exactly ONE row
//   - NO need to close
//   - PHP equivalent: $stmt = $pdo->query($sql); $row = $stmt->fetch()
//
// 3. pool.Exec(ctx, sql, args...)
//   - Executes without returning rows (INSERT, UPDATE, DELETE)
//   - Returns command tag (e.g., "INSERT 0 1")
//   - PHP equivalent: $pdo->exec($sql) or $stmt->execute()
//
// 4. $1, $2, $3... in SQL
//   - These are POSITIONAL PARAMETERS (prevents SQL injection)
//   - PHP equivalent: $pdo->prepare("SELECT * WHERE id = ?")
//   - Go uses $1 not ? (PostgreSQL native syntax)
//
// ═══════════════════════════════════════════════════════════════════════════
func (r *PostPgRepository) GetAll(ctx context.Context, page, perPage int, status string) ([]models.Post, int64, error) {
	// Build the query dynamically based on filters
	// In PHP: You might use a query builder or conditional concatenation
	//
	// GO CONCEPT - DYNAMIC SQL BUILDING:
	// We build the SQL string and arguments slice together
	// This is MANUAL compared to Eloquent but gives full control

	baseQuery := `
        SELECT 
            p.id, p.title, p.slug, p.content, p.status, 
            p.author_id, p.category, p.tags,
            p.created_at, p.updated_at
        FROM posts p
    `
	countQuery := "SELECT COUNT(*) FROM posts p"

	var args []interface{}
	var conditions []string
	argNum := 1 // Track argument position ($1, $2, etc.)

	// Add status filter if provided
	if status != "" {
		conditions = append(conditions, fmt.Sprintf("p.status = $%d", argNum))
		args = append(args, status)
		argNum++
	}

	// Apply WHERE clause if we have conditions
	if len(conditions) > 0 {
		whereClause := " WHERE " + strings.Join(conditions, " AND ")
		baseQuery += whereClause
		countQuery += whereClause
	}

	// Get total count first
	// QueryRow expects exactly one row
	var total int64
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posts: %w", err)
	}

	// Apply pagination defaults
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	if perPage > 100 {
		perPage = 100
	}

	// Add ORDER BY and LIMIT/OFFSET
	// IMPORTANT: Always use LIMIT to prevent fetching millions of rows
	baseQuery += fmt.Sprintf(" ORDER BY p.created_at DESC LIMIT $%d OFFSET $%d", argNum, argNum+1)
	args = append(args, perPage, (page-1)*perPage)

	// Execute the query
	rows, err := r.pool.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch posts: %w", err)
	}
	// ALWAYS defer Close() on rows to prevent connection leaks!
	// In PHP: Statement closes automatically at end of request
	// In Go: We MUST close explicitly
	defer rows.Close()

	// Scan rows into slice
	// In PHP: $posts = $stmt->fetchAll(PDO::FETCH_CLASS, Post::class);
	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var tags []string // PostgreSQL array -> Go slice

		// ═══════════════════════════════════════════════════════════════
		// GO CONCEPT - ROW SCANNING:
		// ═══════════════════════════════════════════════════════════════
		//
		// Scan() maps columns to variables BY POSITION
		// The order must match the SELECT order!
		//
		// PHP: $post->title = $row['title']
		// Go:  rows.Scan(&post.ID, &post.Title, &post.Slug, ...)
		//
		// Common pitfall: Wrong number of variables = error
		// Solution: Be careful with SELECT order, or use pgx.NamedArgs
		//
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Slug,
			&post.Content,
			&post.Status,
			&post.AuthorID,
			&post.Category,
			&tags, // PostgreSQL text[] -> Go []string
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan post row: %w", err)
		}

		// Handle nil tags (PostgreSQL NULL array)
		if tags != nil {
			post.Tags = tags
		} else {
			post.Tags = []string{}
		}

		posts = append(posts, post)
	}

	// ALWAYS check for errors after iteration
	// This catches errors that occurred during row scanning
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating post rows: %w", err)
	}

	// Return empty slice instead of nil (consistent JSON response)
	if posts == nil {
		posts = []models.Post{}
	}

	return posts, total, nil
}

// GetByID retrieves a single post by its ID
func (r *PostPgRepository) GetByID(ctx context.Context, id uint) (*models.Post, error) {
	// QueryRow for single row - no need to close
	var post models.Post
	var tags []string

	// In PHP: $post = Post::find($id);
	query := `
        SELECT id, title, slug, content, status, author_id, category, tags,
               created_at, updated_at
        FROM posts
        WHERE id = $1
    `

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Slug,
		&post.Content,
		&post.Status,
		&post.AuthorID,
		&post.Category,
		&tags,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	// ═══════════════════════════════════════════════════════════════
	// GO CONCEPT - pgx.ErrNoRows:
	// ═══════════════════════════════════════════════════════════════
	//
	// When QueryRow finds no rows, it returns pgx.ErrNoRows
	// This is DIFFERENT from a real error!
	//
	// PHP: Post::find($id) returns null
	// Go:  QueryRow returns pgx.ErrNoRows (must check explicitly)
	//
	// We convert it to our domain error "post not found"
	//
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("post not found")
		}
		return nil, fmt.Errorf("failed to get post by id: %w", err)
	}

	if tags != nil {
		post.Tags = tags
	} else {
		post.Tags = []string{}
	}

	return &post, nil
}

// GetBySlug retrieves a post by its slug
func (r *PostPgRepository) GetBySlug(ctx context.Context, slug string) (*models.Post, error) {
	var post models.Post
	var tags []string

	query := `
        SELECT id, title, slug, content, status, author_id, category, tags,
               created_at, updated_at
        FROM posts
        WHERE slug = $1
    `

	err := r.pool.QueryRow(ctx, query, slug).Scan(
		&post.ID, &post.Title, &post.Slug, &post.Content, &post.Status,
		&post.AuthorID, &post.Category, &tags, &post.CreatedAt, &post.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("post not found")
		}
		return nil, fmt.Errorf("failed to get post by slug: %w", err)
	}

	if tags != nil {
		post.Tags = tags
	} else {
		post.Tags = []string{}
	}

	return &post, nil
}

// Create inserts a new post and returns it with the generated ID
func (r *PostPgRepository) Create(ctx context.Context, post *models.Post) error {
	// ═══════════════════════════════════════════════════════════════
	// GO CONCEPT - RETURNING CLAUSE:
	// ═══════════════════════════════════════════════════════════════
	//
	// PostgreSQL's RETURNING clause returns the inserted row
	// This is MUCH better than: INSERT then SELECT to get ID
	//
	// PHP: $post = Post::create([...]); $post->id (auto-filled)
	// Go:  We use RETURNING id to get it in one query
	//
	// Also: NOW() for timestamps instead of Go's time.Now()
	// This ensures database and app timestamps match
	//
	query := `
        INSERT INTO posts (title, slug, content, status, author_id, category, tags, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
        RETURNING id, created_at, updated_at
    `

	// Handle nil tags
	tags := post.Tags
	if tags == nil {
		tags = []string{}
	}

	err := r.pool.QueryRow(
		ctx,
		query,
		post.Title,
		post.Slug,
		post.Content,
		post.Status,
		post.AuthorID,
		post.Category,
		tags, // Go []string -> PostgreSQL text[]
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		// Check for unique constraint violation (slug already exists)
		// PostgreSQL error code 23505 = unique_violation
		if strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "unique") {
			return errors.New("slug already exists")
		}
		return fmt.Errorf("failed to create post: %w", err)
	}

	return nil
}

// Update modifies an existing post
func (r *PostPgRepository) Update(ctx context.Context, post *models.Post) error {
	// In PHP: $post->update([...]) or Post::where('id', $id)->update([...])
	//
	// For PARTIAL updates, we have two options:
	// Option 1: Build dynamic SQL (complex but efficient)
	// Option 2: Update all fields (simpler, slightly wasteful)
	//
	// We use Option 2 here for clarity
	query := `
        UPDATE posts
        SET title = $1, slug = $2, content = $3, status = $4,
            category = $5, tags = $6, updated_at = NOW()
        WHERE id = $7
        RETURNING updated_at
    `

	tags := post.Tags
	if tags == nil {
		tags = []string{}
	}

	err := r.pool.QueryRow(
		ctx,
		query,
		post.Title,
		post.Slug,
		post.Content,
		post.Status,
		post.Category,
		tags,
		post.ID,
	).Scan(&post.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("post not found")
		}
		if strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "unique") {
			return errors.New("slug already exists")
		}
		return fmt.Errorf("failed to update post: %w", err)
	}

	return nil
}

// Delete removes a post from the database
func (r *PostPgRepository) Delete(ctx context.Context, id uint) error {
	// In PHP: Post::destroy($id)
	//
	// We use RETURNING id to verify the row was actually deleted
	query := `DELETE FROM posts WHERE id = $1 RETURNING id`

	var deletedID uint
	err := r.pool.QueryRow(ctx, query, id).Scan(&deletedID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("post not found")
		}
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}

// GetByAuthor retrieves all posts by a specific author
func (r *PostPgRepository) GetByAuthor(ctx context.Context, authorID uint) ([]models.Post, error) {
	query := `
        SELECT id, title, slug, content, status, author_id, category, tags,
               created_at, updated_at
        FROM posts
        WHERE author_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.pool.Query(ctx, query, authorID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch posts by author: %w", err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var tags []string

		err := rows.Scan(
			&post.ID, &post.Title, &post.Slug, &post.Content, &post.Status,
			&post.AuthorID, &post.Category, &tags, &post.CreatedAt, &post.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post row: %w", err)
		}

		if tags != nil {
			post.Tags = tags
		} else {
			post.Tags = []string{}
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating post rows: %w", err)
	}

	if posts == nil {
		posts = []models.Post{}
	}

	return posts, nil
}

// Count returns total count of posts (helper for pagination)
func (r *PostPgRepository) Count(ctx context.Context, status string) (int64, error) {
	var count int64
	var err error

	if status != "" {
		err = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM posts WHERE status = $1", status).Scan(&count)
	} else {
		err = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM posts").Scan(&count)
	}

	if err != nil {
		return 0, fmt.Errorf("failed to count posts: %w", err)
	}

	return count, nil
}

// math import needed for ceil in pagination
var _ = math.Ceil // Prevent unused import error
