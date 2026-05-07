# Go CMS

A lightweight Content Management System built with Go and Fiber, designed as a learning project to demonstrate clean architecture principles and best practices for developers transitioning from PHP/Laravel to Go.

## Features

- **RESTful API** for managing posts and users
- **JWT Authentication** for secure endpoints
- **Clean Architecture** with separation of concerns:
  - Handlers (HTTP layer)
  - Services (Business logic)
  - Repositories (Data access)
  - Models (Data structures)
- **Middleware** for logging, request ID, CORS, and recovery
- **DTO Patterns** for request/response validation
- **Paginated Responses** for list endpoints
- **Mock Repositories** for easy testing
- **Configuration Management** with environment support
- **Comprehensive Error Handling**
- **Health Check Endpoint**

## Architecture Overview

The project follows a layered architecture inspired by Domain-Driven Design and Clean Architecture principles:

```
HTTP REQUEST
    │
    ▼
MIDDLEWARE LAYER
    │
    ▼
HANDLER LAYER (Controllers)
    │
    ▼
SERVICE LAYER (Business Logic)
    │
    ▼
REPOSITORY LAYER (Data Access)
    │
    ▼
DATA SOURCES (Mock/Database)
```

Each layer has a specific responsibility:
- **Handlers**: Handle HTTP requests/responses, validation, and routing
- **Services**: Contain business logic and orchestrate repository calls
- **Repositories**: Abstract data access operations
- **Models**: Define data structures and DTOs

## Technology Stack

- **Language**: Go 1.22+
- **Web Framework**: Fiber v2
- **JSON Handling**: Built-in encoding/json
- **Environment Variables**: godotenv (via config package)
- **Routing**: Fiber router
- **Mock Data**: In-memory repositories for development/testing

## Project Structure

```
go-cms/
├── cmd/
│   └── server/                 # Application entry point
├── internal/
│   ├── config/                 # Configuration management
│   ├── handlers/               # HTTP handlers (controllers)
│   ├── middleware/             # Custom middleware
│   ├── models/                 # Data models and DTOs
│   ├── repositories/           # Data access interfaces and implementations
│   ├── routes/                 # Route definitions
│   └── services/               # Business logic services
├── go.mod                      # Go module definition
├── go.sum                      # Dependency checksums
└── README.md                   # This file
```

## Getting Started

### Prerequisites

- Go 1.22 or higher
- Git

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/talhag3/go-cms.git
   cd go-cms
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run the application:
   ```bash
   go run cmd/server/main.go
   ```

4. The server will start on `http://localhost:3000` (default configuration)

### Configuration

Configuration is managed through the `internal/config` package. Default values are set in the code, but you can override them using environment variables:

- `APP_NAME`: Application name (default: "Go CMS")
- `APP_ENV`: Environment (development/staging/production) (default: "development")
- `SERVER_HOST`: Server host (default: "localhost")
- `SERVER_PORT`: Server port (default: "3000")

## API Endpoints

### Authentication

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/auth/login` | Login user | No |

### Posts

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/posts` | List all posts (paginated) | No |
| GET | `/api/v1/posts/:id` | Get post by ID | No |
| POST | `/api/v1/posts` | Create new post | Yes |
| PUT | `/api/v1/posts/:id` | Update post | Yes |
| DELETE | `/api/v1/posts/:id` | Delete post | Yes |

### Users

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/users` | List all users (paginated) | Yes |

### Health Check

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/health` | Health check endpoint | No |

## Best Practices Demonstrated

This project illustrates several Go best practices for developers coming from PHP/Laravel backgrounds:

### 1. Folder Structure = Namespace
- PHP: `namespace App\Services;`
- Go: `package services` (folder name = package name)

### 2. Capitalization = Visibility
- PHP: `public/private/protected`
- Go: Capitalized = exported, lowercase = unexported

### 3. Error Handling
- PHP: `try/catch` with exceptions
- Go: `if err != nil { return err }` (explicit, no hidden control flow)

### 4. Dependency Injection
- PHP: Constructor injection via container
- Go: Pass interfaces via constructor functions (manual DI)

### 5. Interfaces for Abstraction
- PHP: Explicit `implements`
- Go: Implicit (if it has the methods, it implements)
- Benefit: Easy swapping of implementations (mock → real database)

### 6. Keep main() Thin
- Only: config loading, DI wiring, server start
- All logic in packages

### 7. Internal/ for Private Code
- Code in `internal/` can't be imported by other projects
- Like "private" packages

### 8. Pkg/ for Reusable Code
- Code in `pkg/` is designed to be reused
- Like "public" packages

### 9. Use Pointers Wisely
- Use `*T` when you need nil, to modify original, or avoid copying large structs
- Don't use for simple values/primitives

### 10. Defer for Cleanup
- PHP: `try/finally` or `__destruct`
- Go: `defer func() { ... }` (executes when function returns)

## Development

### Running Tests

Currently, the project uses mock repositories for demonstration purposes. To run tests:

```bash
go test ./...
```

### Adding Features

1. Follow the existing patterns in the codebase
2. Add new models in `internal/models/`
3. Create repository interfaces and implementations in `internal/repositories/`
4. Add business logic in `internal/services/`
5. Create handlers in `internal/handlers/`
6. Register routes in `internal/routes/`
7. Add any necessary middleware in `internal/middleware/`

## License

This project is open source and available under the [MIT License](LICENSE).

## Acknowledgments

- Inspired by Laravel's architecture and conventions
- Built with [Fiber](https://gofiber.io/) - Express inspired web framework for Go
- Designed to help PHP developers transition to Go