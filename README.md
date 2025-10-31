# 🚀 UGin - Ultimate Gin API Boilerplate

[![Go Version](https://img.shields.io/badge/Go-1.23-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Gin Version](https://img.shields.io/badge/Gin-1.10.0-00ADD8?style=flat)](https://github.com/gin-gonic/gin)
[![GORM Version](https://img.shields.io/badge/GORM-1.30.0-00ADD8?style=flat)](https://gorm.io/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A production-ready REST API boilerplate written in Go with Gin Framework, featuring JWT authentication, GORM ORM, comprehensive middleware support, and **interactive Swagger documentation**.

## ✨ Features

- 🎯 **Modern Go** - Built with Go 1.23
- ⚡ **Gin Framework** - Fast HTTP web framework
- 🗄️ **Multi-Database Support** - SQLite, MySQL, PostgreSQL via GORM
- 🔐 **JWT Authentication** - Secure token-based authentication
- 📝 **Comprehensive Logging** - Application, database, and access logs
- 🔍 **Advanced Querying** - Built-in filtering, search, and pagination
- 🛡️ **Security Middleware** - CORS, rate limiting, and more
- 📦 **Gzip Compression** - Automatic response compression
- 🔄 **Hot Reload** - Development mode with auto-reload
- 📊 **Structured Logging** - Using logrus
- 🏗️ **Clean Architecture** - Repository pattern with dependency injection
- 🧪 **Fully Testable** - Interface-based design for easy mocking
- 🌐 **Context Propagation** - Proper context handling throughout the stack
- ♻️ **Graceful Shutdown** - Proper resource cleanup on exit
- 📚 **Swagger/OpenAPI** - Interactive API documentation with Swagger UI

## 📋 Table of Contents

- [Quick Start](#-quick-start)
- [Project Structure](#-project-structure)
- [Configuration](#-configuration)
- [API Endpoints](#-api-endpoints)
  - [Interactive API Documentation (Swagger)](#-interactive-api-documentation)
- [Database](#-database)
- [Middleware](#-middleware)
- [Logging](#-logging)
- [Development](#-development)
  - [Swagger Documentation](#swagger-documentation)
- [Docker Support](#-docker-support)

## 📚 Additional Documentation

- **[QUICK_START.md](QUICK_START.md)** - Fast setup guide with common commands
- **[SWAGGER_GUIDE.md](SWAGGER_GUIDE.md)** - Complete Swagger/OpenAPI documentation guide

## 🚀 Quick Start

### Prerequisites

- Go 1.23 or higher
- Git

### Installation

```bash
# Clone the repository
git clone https://github.com/yakuter/ugin.git
cd ugin

# Download dependencies
go mod download

# Build the application
make build

# Run the application
./bin/ugin
```

Or use the Makefile for a simpler workflow:

```bash
# Build and run in one command
make run

# Or run directly without building (development mode)
make run-dev
```

The server will start at `http://127.0.0.1:8081`

**🎉 Access Swagger UI:** `http://127.0.0.1:8081/swagger/index.html`

## 📁 Project Structure

```
ugin/
├── cmd/                      # Application entry points
│   └── ugin/
│       └── main.go           # Main entry point (simple!)
├── internal/                 # Private application code
│   ├── core/                 # Application core
│   │   ├── app.go            # Application lifecycle
│   │   ├── database.go       # Database initialization
│   │   └── router.go         # Router setup
│   ├── domain/               # Domain models (entities)
│   │   ├── post.go
│   │   ├── user.go
│   │   └── auth.go
│   ├── repository/           # Data access layer
│   │   ├── repository.go     # Repository interfaces
│   │   └── gormrepo/         # GORM implementations
│   │       ├── post.go
│   │       └── user.go
│   ├── service/              # Business logic layer
│   │   ├── interfaces.go     # Service interfaces
│   │   ├── post.go
│   │   ├── auth.go
│   │   └── post_test.go      # Example tests
│   ├── handler/              # HTTP handlers
│   │   └── http/
│   │       ├── post.go
│   │       ├── auth.go
│   │       ├── admin.go
│   │       └── middleware.go
│   └── config/               # Configuration management
│       └── config.go
├── pkg/                      # Public reusable packages
│   └── logger/               # Logging utilities
│       └── logger.go
├── containers/               # Docker configuration
│   ├── composes/             # Docker compose files
│   └── images/               # Dockerfiles
├── bin/                      # Compiled binaries (gitignored)
├── config.yml                # Application configuration
├── Makefile                  # Build automation
└── go.mod                    # Go module definition
```

This structure follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout) with **Clean Architecture** principles.

### Architecture Layers

1. **Core Layer** (`internal/core/`) - Application lifecycle and wiring
2. **Domain Layer** (`internal/domain/`) - Pure business entities
3. **Repository Layer** (`internal/repository/`) - Data access interfaces and implementations
4. **Service Layer** (`internal/service/`) - Business logic and use cases
5. **Handler Layer** (`internal/handler/`) - HTTP request/response handling
6. **Infrastructure** (`pkg/`, `internal/config/`) - External concerns

**Key Principles:**
- ✅ **Dependency Injection** - No global state
- ✅ **Interface-based** - Easy to mock and test
- ✅ **Context Propagation** - Proper timeout and cancellation
- ✅ **Clean Separation** - Each layer has a single responsibility
- ✅ **Simple main.go** - Entry point is just 15 lines!

## ⚙️ Configuration

Edit `config.yml` to configure your application:

```yaml
database:
  driver: "sqlite"      # Options: sqlite, mysql, postgres
  dbname: "ugin"
  username: "user"      # Not required for SQLite
  password: "password"  # Not required for SQLite
  host: "localhost"     # Not required for SQLite
  port: "5432"          # Not required for SQLite
  logmode: true         # Enable SQL query logging

server:
  port: "8081"
  secret: "mySecretKey"                    # JWT secret key
  accessTokenExpireDuration: 1             # Hours
  refreshTokenExpireDuration: 1            # Hours
  limitCountPerRequest: 1                  # Rate limit per request
```

### Database Drivers

**SQLite** (Default - No setup required):
```yaml
database:
  driver: "sqlite"
  dbname: "ugin"
  logmode: true
```

**MySQL**:
```yaml
database:
  driver: "mysql"
  dbname: "ugin"
  username: "root"
  password: "password"
  host: "localhost"
  port: "3306"
```

**PostgreSQL**:
```yaml
database:
  driver: "postgres"
  dbname: "ugin"
  username: "user"
  password: "password"
  host: "localhost"
  port: "5432"
  sslmode: "disable"
```

## 📡 API Endpoints

All API endpoints are versioned with `/api/v1` prefix.

### 📚 Interactive API Documentation

**Swagger UI** is available at: **`http://localhost:8081/swagger/index.html`**

- 📖 View all endpoints with detailed documentation
- 🧪 Test API endpoints directly from browser
- 🔐 Test authentication with JWT tokens
- 📋 See request/response examples with real data

See [SWAGGER_GUIDE.md](SWAGGER_GUIDE.md) for detailed usage.

### Authentication Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/signup` | Register a new user |
| POST | `/api/v1/auth/signin` | Sign in and get JWT tokens |
| POST | `/api/v1/auth/refresh` | Refresh access token |
| POST | `/api/v1/auth/check` | Validate token |

### Posts Endpoints (Public)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/posts` | Get all posts (supports pagination) |
| GET | `/api/v1/posts/:id` | Get a single post by ID |
| POST | `/api/v1/posts` | Create a new post |
| PUT | `/api/v1/posts/:id` | Update an existing post |
| DELETE | `/api/v1/posts/:id` | Delete a post |

### Posts Endpoints (JWT Protected)

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/v1/postsjwt` | Get all posts | JWT |
| GET | `/api/v1/postsjwt/:id` | Get a single post | JWT |
| POST | `/api/v1/postsjwt` | Create a new post | JWT |
| PUT | `/api/v1/postsjwt/:id` | Update a post | JWT |
| DELETE | `/api/v1/postsjwt/:id` | Delete a post | JWT |

### Admin Endpoints (Basic Auth)

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/admin/dashboard` | Admin dashboard | Basic Auth |

**Default credentials**: `username1:password1`, `username2:password2`, `username3:password3`

### Query Parameters

All list endpoints support advanced querying:

```
GET /posts/?Limit=10&Offset=0&Sort=ID&Order=DESC&Search=keyword
```

| Parameter | Description | Example |
|-----------|-------------|---------|
| `Limit` | Number of records to return | `Limit=25` |
| `Offset` | Number of records to skip | `Offset=0` |
| `Sort` | Field to sort by | `Sort=ID` |
| `Order` | Sort order (ASC/DESC) | `Order=DESC` |
| `Search` | Search keyword | `Search=hello` |

### Example API Requests

#### Create a Post

```bash
curl -X POST http://localhost:8081/api/v1/posts \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Hello World",
    "description": "This is a sample post",
    "tags": [
      {
        "name": "golang",
        "description": "Go programming language"
      },
      {
        "name": "api",
        "description": "REST API"
      }
    ]
  }'
```

#### Get Posts with Pagination

```bash
curl "http://localhost:8081/api/v1/posts?Limit=10&Offset=0&Sort=id&Order=DESC"
```

#### Sign Up

```bash
curl -X POST http://localhost:8081/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "master_password": "password123"
  }'
```

#### Sign In

```bash
curl -X POST http://localhost:8081/api/v1/auth/signin \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "master_password": "password123"
  }'
```

**Response:**
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "transmission_key": "...",
  "access_token_expires_at": "2025-11-01T10:00:00Z",
  "refresh_token_expires_at": "2025-11-02T10:00:00Z"
}
```

#### Access Protected Endpoint

```bash
curl http://localhost:8081/api/v1/postsjwt \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

#### Refresh Token

```bash
curl -X POST http://localhost:8081/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```

## 🗄️ Database

### Domain Models

UGin includes example domain models demonstrating relationships:

**Post Model** (`internal/domain/post.go`):
```go
type Post struct {
    ID          uint       `json:"id" gorm:"primarykey"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
    DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`
    Name        string     `json:"name" gorm:"type:varchar(255);not null"`
    Description string     `json:"description" gorm:"type:text"`
    Tags        []Tag      `json:"tags,omitempty" gorm:"foreignKey:PostID"`
}
```

**Tag Model** (`internal/domain/post.go`):
```go
type Tag struct {
    ID          uint       `json:"id" gorm:"primarykey"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
    DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`
    PostID      uint       `json:"post_id" gorm:"index;not null"`
    Name        string     `json:"name" gorm:"type:varchar(255);not null"`
    Description string     `json:"description" gorm:"type:text"`
}
```

**User Model** (`internal/domain/user.go`):
```go
type User struct {
    ID             uint       `json:"id" gorm:"primarykey"`
    CreatedAt      time.Time  `json:"created_at"`
    UpdatedAt      time.Time  `json:"updated_at"`
    DeletedAt      *time.Time `json:"deleted_at,omitempty" gorm:"index"`
    Email          string     `json:"email" gorm:"uniqueIndex;not null"`
    MasterPassword string     `json:"-" gorm:"not null"` // Never exposed in JSON
}
```

### Repository Pattern

The application uses the Repository pattern for data access:

```go
// Repository interface (internal/repository/repository.go)
type PostRepository interface {
    GetByID(ctx context.Context, id string) (*domain.Post, error)
    List(ctx context.Context, filter ListFilter) ([]*domain.Post, *ListResult, error)
    Create(ctx context.Context, post *domain.Post) error
    Update(ctx context.Context, post *domain.Post) error
    Delete(ctx context.Context, id string) error
}

// GORM implementation (internal/repository/gormrepo/post.go)
type postRepository struct {
    db *gorm.DB
}
```

### Migrations

Migrations run automatically on application startup in `cmd/ugin/main.go`:

```go
func autoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &domain.Post{},
        &domain.Tag{},
        &domain.User{},
    )
}
```

## 🛡️ Middleware

### Built-in Middleware

1. **Logger** - Request logging (Gin)
2. **Recovery** - Panic recovery (Gin)
3. **CORS** - Cross-Origin Resource Sharing
4. **Gzip** - Response compression
5. **Security** - Security headers
6. **Rate Limiting** - Request rate limiting (per IP)
7. **JWT Auth** - Token validation

### Using JWT Authentication

The JWT middleware is applied to protected routes:

```go
// In main.go
postsJWT := v1.Group("/postsjwt")
postsJWT.Use(httpHandler.JWTAuth(authService))
{
    postsJWT.GET("", postHandler.List)
    // ... other protected routes
}
```

Protected endpoints require an `Authorization` header:
```
Authorization: Bearer YOUR_ACCESS_TOKEN
```

### Custom Middleware

Add custom middleware in `internal/handler/http/middleware.go`:

```go
func CustomMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Your logic here
        c.Next()
    }
}
```

Then register it in `cmd/ugin/main.go`:
```go
router.Use(httpHandler.CustomMiddleware())
```

## 🧪 Testing

The new architecture makes testing easy with dependency injection and interfaces.

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test -v ./internal/service/...
```

### Example Test

Here's how to test a service with mocked dependencies (`internal/service/post_test.go`):

```go
func TestPostService_GetByID(t *testing.T) {
    // Create mock repository
    mockRepo := &mockPostRepository{
        getByIDFunc: func(ctx context.Context, id string) (*domain.Post, error) {
            return &domain.Post{
                ID:          1,
                Name:        "Test Post",
                Description: "Test Description",
            }, nil
        },
    }

    // Create service with mock
    svc := service.NewPostService(mockRepo, &mockLogger{})

    // Test
    post, err := svc.GetByID(context.Background(), "1")
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    if post.Name != "Test Post" {
        t.Errorf("expected 'Test Post', got '%s'", post.Name)
    }
}
```

### Benefits of This Architecture

✅ **Easy to Mock** - All dependencies are interfaces  
✅ **Isolated Tests** - No global state to manage  
✅ **Fast Tests** - No database required for service tests  
✅ **Reliable** - Tests don't affect each other

## 📝 Logging

UGin provides three types of logs:

### Application Log (`ugin.log`)
General application events and errors:
```
INFO 2025-10-31T10:05:53+03:00 Server is starting at 127.0.0.1:8081
ERROR 2025-10-31T10:06:15+03:00 Failed to connect to database
```

### Database Log (`ugin.db.log`)
SQL queries and database operations:
```
2025/10/31 10:05:53 /Users/user/ugin/pkg/database/database.go:80
[0.017ms] [rows:-] SELECT count(*) FROM sqlite_master WHERE type='table'
```

### Access Log (`ugin.access.log`)
HTTP request logs:
```
[GIN] 2025/10/31 - 10:05:53 | 200 | 9.255625ms | 127.0.0.1 | GET "/posts/"
```

### Log Levels

Configure log verbosity using the `GIN_MODE` environment variable:

```bash
# Development mode (verbose)
export GIN_MODE=debug

# Test mode
export GIN_MODE=test

# Production mode (minimal logging)
export GIN_MODE=release
```

## 🔧 Development

### Available Make Commands

View all available commands:
```bash
make help
```

### Run in Development Mode

```bash
# Set debug mode
export GIN_MODE=debug

# Run directly (no build step)
make run-dev

# Or run with build
make run
```

### Build Commands

```bash
# Build for development
make build

# Build for production (optimized, smaller binary)
make build-prod

# The binary will be created in ./bin/ugin
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage
# This generates coverage.html that you can open in a browser
```

### Code Quality

```bash
# Format code
make fmt

# Run go vet
make vet

# Run linter (requires golangci-lint)
make lint

# Run all checks (format + vet + test)
make check
```

### Dependency Management

```bash
# Download dependencies
make deps

# Update dependencies to latest versions
make deps-update
```

### Clean Build Artifacts

```bash
# Remove binaries and log files
make clean
```

### Swagger Documentation

```bash
# Generate Swagger docs
make swagger

# Generate docs and run
make run-swagger

# View documentation
# Open http://localhost:8081/swagger/index.html
```

**Note:** You need to install `swag` CLI first:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

See [SWAGGER_GUIDE.md](SWAGGER_GUIDE.md) for detailed Swagger usage.

## 🐳 Docker Support

### Build Docker Image

```bash
make build-image
```

### Run with Docker Compose

**With MySQL**:
```bash
# Start application with MySQL
make run-app-mysql

# Stop MySQL containers
make clean-app-mysql
```

**With PostgreSQL**:
```bash
# Start application with PostgreSQL
make run-app-postgres

# Stop PostgreSQL containers
make clean-app-postgres
```

### Manual Docker Commands

```bash
# Build image
docker build -t ugin:latest -f containers/images/Dockerfile .

# Run container
docker run -p 8081:8081 -v $(pwd)/config.yml:/app/config.yml ugin:latest
```

## 📦 Dependencies

Core dependencies:

- [gin-gonic/gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [gorm.io/gorm](https://gorm.io/) - ORM library
- [spf13/viper](https://github.com/spf13/viper) - Configuration management
- [golang-jwt/jwt](https://github.com/golang-jwt/jwt) - JWT implementation
- [sirupsen/logrus](https://github.com/sirupsen/logrus) - Structured logging
- [didip/tollbooth](https://github.com/didip/tollbooth) - Rate limiting
- [swaggo/swag](https://github.com/swaggo/swag) - Swagger documentation

See `go.mod` for the complete list.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Gin](https://github.com/gin-gonic/gin) - Amazing HTTP web framework
- [GORM](https://gorm.io/) - Fantastic ORM library
- [Viper](https://github.com/spf13/viper) - Complete configuration solution

## 📞 Support

If you have any questions or need help, please open an issue on GitHub.

---

Made with ❤️ using Go
