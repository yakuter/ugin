# Quick Start Guide

## First Time Setup

```bash
# Clone the repository
git clone https://github.com/yakuter/ugin.git
cd ugin

# Download dependencies
make deps

# Build the application
make build

# Run the application
./bin/ugin
```

## Most Used Commands

```bash
# Development
make run-dev          # Run without building (fast iteration)
make run              # Build then run
make build            # Build to ./bin/ugin

# Testing
make test             # Run tests
make test-coverage    # Generate coverage report

# Code Quality
make fmt              # Format code
make check            # Run all checks

# Clean
make clean            # Remove build artifacts

# Help
make help             # Show all commands
```

## Environment Variables

```bash
# Set Gin mode
export GIN_MODE=debug     # Development (verbose logs)
export GIN_MODE=release   # Production (minimal logs)
export GIN_MODE=test      # Testing
```

## Configuration

Edit `config.yml`:

```yaml
database:
  driver: "sqlite"        # sqlite, mysql, or postgres
  dbname: "ugin"
  logmode: true

server:
  port: "8081"
  secret: "mySecretKey"
```

## Quick API Test

```bash
# Start server
make run-dev

# In another terminal:

# Get all posts
curl http://localhost:8081/api/v1/posts

# Create a post
curl -X POST http://localhost:8081/api/v1/posts \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","description":"Test post"}'

# Sign up
curl -X POST http://localhost:8081/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","master_password":"password123"}'

# Sign in (get JWT token)
curl -X POST http://localhost:8081/api/v1/auth/signin \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","master_password":"password123"}'
```

## Project Structure Quick Reference

```
cmd/ugin/main.go              → Application entry point (15 lines!)
internal/core/                → Application core
  ├── app.go                  → App lifecycle & DI
  ├── database.go             → DB initialization
  └── router.go               → Route setup
internal/domain/              → Domain models
internal/repository/          → Data access layer
internal/service/             → Business logic
internal/handler/http/        → HTTP handlers
internal/config/              → Configuration
pkg/logger/                   → Logging utilities
```

## Docker Quick Start

```bash
# Build image
make build-image

# Run with PostgreSQL
make run-app-postgres

# Stop
make clean-app-postgres
```

## Common Issues

**Port already in use:**
```bash
# Change port in config.yml
server:
  port: "8082"  # or any available port
```

**Database connection failed:**
- Check `config.yml` database settings
- For SQLite: No setup needed
- For MySQL/PostgreSQL: Ensure database is running

**Module errors:**
```bash
make deps
go mod tidy
```

## Next Steps

1. Read the [README.md](README.md) for detailed documentation
2. Check [MIGRATION.md](MIGRATION.md) if upgrading from old structure
3. Browse the code in `cmd/`, `controller/`, and `service/`
4. Customize `config.yml` for your needs
5. Start building your API!

## Resources

- [Gin Framework Docs](https://gin-gonic.com/docs/)
- [GORM Docs](https://gorm.io/docs/)
- [Go Project Layout](https://github.com/golang-standards/project-layout)

---

**Pro Tip:** Use `make run-dev` during development for fast iteration without building!

