// internal/database/migrations.go
package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RunMigrations creates database tables
//
// PHP COMPARISON:
// Laravel: php artisan migrate (reads database/migrations/)
// Here: We run SQL directly (simple approach)
//
// For production, consider:
// - golang-migrate/migrate (versioned migrations)
// - goose (another migration tool)
// - Atlas (modern schema management)
//
// GO CONCEPT - CONTEXT:
// context.Context is passed through the call chain
// It carries: deadlines, cancellation signals, request-scoped values
// In PHP: Not really applicable (each request is isolated)
// In Go: One process handles many requests, so we need context
func RunMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	log.Println("🔄 Running database migrations...")

	// Create users table
	// In PHP/Laravel: Schema::create('users', function (Blueprint $table) { ... })
	usersTableSQL := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(100) NOT NULL UNIQUE,
        email VARCHAR(255) NOT NULL UNIQUE,
        password_hash VARCHAR(255) NOT NULL DEFAULT '',
        role VARCHAR(50) NOT NULL DEFAULT 'author',
        created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
    );

    -- Create index on email for faster login queries
    CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
    `

	if _, err := pool.Exec(ctx, usersTableSQL); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}
	log.Println("  ✅ Users table ready")

	// Create posts table
	postsTableSQL := `
    CREATE TABLE IF NOT EXISTS posts (
        id SERIAL PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        slug VARCHAR(255) NOT NULL UNIQUE,
        content TEXT NOT NULL DEFAULT '',
        status VARCHAR(50) NOT NULL DEFAULT 'draft',
        author_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        category VARCHAR(100) NOT NULL DEFAULT '',
        tags TEXT[] NOT NULL DEFAULT '{}',
        created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
    );

    -- Index for slug lookups (public API)
    CREATE INDEX IF NOT EXISTS idx_posts_slug ON posts(slug);

    -- Index for author posts
    CREATE INDEX IF NOT EXISTS idx_posts_author ON posts(author_id);

    -- Index for status filtering
    CREATE INDEX IF NOT EXISTS idx_posts_status ON posts(status);

    -- Composite index for common query pattern
    CREATE INDEX IF NOT EXISTS idx_posts_status_created ON posts(status, created_at DESC);
    `

	if _, err := pool.Exec(ctx, postsTableSQL); err != nil {
		return fmt.Errorf("failed to create posts table: %w", err)
	}
	log.Println("  ✅ Posts table ready")

	// Insert seed data (only if tables are empty)
	if err := seedData(ctx, pool); err != nil {
		return fmt.Errorf("failed to seed data: %w", err)
	}

	log.Println("🎉 Migrations completed!")
	return nil
}

// seedData inserts initial data
// In PHP: DatabaseSeeder.php
func seedData(ctx context.Context, pool *pgxpool.Pool) error {
	// Check if we already have data
	var userCount int
	err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&userCount)
	if err != nil {
		return err
	}

	if userCount > 0 {
		log.Println("  ℹ️  Seed data already exists, skipping...")
		return nil
	}

	log.Println("  🌱 Inserting seed data...")

	// Insert users
	// In PHP: User::create(['username' => 'admin', ...])
	//
	// GO CONCEPT - MULTIPLE ROW INSERT:
	// We can insert multiple rows in one query (faster than individual inserts)
	userInsertSQL := `
    INSERT INTO users (username, email, password_hash, role) VALUES
        ('admin', 'admin@cms.com', '$2b$12$mock_hash_for_admin', 'admin'),
        ('editor', 'editor@cms.com', '$2b$12$mock_hash_for_editor', 'editor'),
        ('author', 'author@cms.com', '$2b$12$mock_hash_for_author', 'author')
    RETURNING id, username
    `

	rows, err := pool.Query(ctx, userInsertSQL)
	if err != nil {
		return fmt.Errorf("failed to insert users: %w", err)
	}
	defer rows.Close()

	log.Println("  ✅ Users seeded")

	// Insert posts
	postInsertSQL := `
    INSERT INTO posts (title, slug, content, status, author_id, category, tags) VALUES
        ('Getting Started with Go', 'getting-started-with-go',
         'Go is a statically typed, compiled language designed for simplicity and efficiency.',
         'published', 1, 'Programming', ARRAY['go', 'tutorial', 'beginner']),

        ('Understanding Fiber Framework', 'understanding-fiber-framework',
         'Fiber is an Express.js inspired web framework built on top of Fasthttp.',
         'published', 1, 'Programming', ARRAY['fiber', 'go', 'web-framework']),

        ('Draft: Advanced Go Patterns', 'advanced-go-patterns',
         'This post covers advanced design patterns in Go including worker pools.',
         'draft', 2, 'Programming', ARRAY['go', 'patterns', 'advanced']),

        ('Building REST APIs in Go', 'building-rest-apis-go',
         'Learn how to build professional REST APIs using Go and Fiber framework.',
         'published', 2, 'API', ARRAY['api', 'rest', 'go', 'fiber']),

        ('Archived: Old Tutorial', 'old-tutorial',
         'This is an old tutorial that has been archived.',
         'archived', 1, 'Old', ARRAY['old'])
    `

	if _, err := pool.Exec(ctx, postInsertSQL); err != nil {
		return fmt.Errorf("failed to insert posts: %w", err)
	}

	log.Println("  ✅ Posts seeded")

	return nil
}
