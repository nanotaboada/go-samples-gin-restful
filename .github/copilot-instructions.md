# Copilot Instructions

## Project Overview

REST API for managing football players built with Go 1.25+ and Gin Web Framework. Implements CRUD operations backed by SQLite with GORM ORM, includes Swagger documentation, in-memory caching, and comprehensive testing. Part of a cross-language comparison study (Java, .NET, TypeScript, Python, Go, Rust).

## Structure

```text
main.go           - application entry point (Gin setup, DB connection, route registration)
go.mod            - Go module dependencies
/route            - route registration with caching middleware
/controller       - HTTP handlers (request/response logic)
/service          - business logic (GORM interactions)
/data             - database connection setup
/model            - domain models (Player struct)
/storage          - SQLite database file (players-sqlite3.db, pre-seeded)
/docs             - auto-generated Swagger docs (DO NOT EDIT manually)
/tests            - integration tests with testify assertions
/scripts          - Docker entrypoint and healthcheck scripts
/.github          - CI/CD workflows (go-ci.yml for build/test, go-cd.yml for releases)
```

## Build, Run, Test Commands

All commands validated and used in CI/CD. Always run from repository root.

Development:
- Install dependencies: `go mod download`
- Run server: `go run .` (starts on port 9000)
- Format code: `go fmt ./...` (automatic in most editors)
- Build binary: `go build` or `go build -v ./...`
- Clean dependencies: `go mod tidy`

Testing:
- All tests: `go test ./...`
- With coverage: `go test -v ./... -coverpkg=github.com/nanotaboada/go-samples-gin-restful/service,github.com/nanotaboada/go-samples-gin-restful/controller,github.com/nanotaboada/go-samples-gin-restful/route -covermode=atomic -coverprofile=coverage.out`
- View coverage: `go tool cover -html=coverage.out`
- Target: 80%+ coverage for service, controller, route packages

Docker:
- Build and run: `docker compose up --build`
- Stop: `docker compose down`
- Reset database: `docker compose down -v` (removes volume)

Documentation:
- Regenerate Swagger: `swag init` (after updating API comments)
- View docs: http://localhost:9000/swagger/index.html (when server running)

## CI/CD Validation

CI runs on push/PR to master (go-ci.yml):
1. Build: `go build -v ./...`
2. Test with race detector: `go test -v ./... -coverpkg=... -covermode=atomic -coverprofile=coverage.out`
3. Lint commits: commitlint validates Conventional Commits format
4. Upload coverage to Codecov (target: 80%+)

CD runs on version tags (go-cd.yml):
- Validates player name in tag (e.g., v1.0.0-ademir)
- Builds Docker image
- Publishes to GitHub Container Registry with three tags (semver, player name, latest)
- Creates GitHub Release with auto-generated changelog

## Stack

- Language: Go 1.25+
- Framework: Gin Web Framework
- Database: SQLite + GORM ORM
- Caching: gin-contrib/cache (in-memory, 1 hour for GET requests)
- Testing: Go testing package + testify/assert
- Linting: golangci-lint
- API Documentation: Swaggo (Swagger generation from comments)

## Project Patterns

- Architecture: Layered (Routes → Controllers → Services → Data)
- Dependency Injection: Explicit passing (Go doesn't have built-in DI)
- Error Handling: Explicit error returns, check errors immediately
- Logging: Standard library `log` package (structured logging with slog for production)

## Code Conventions

- File Names: snake_case for all files (`player_controller.go`, `player_service.go`)
- Package Structure: Flat structure with feature-based packages (`controller/`, `service/`, `data/`, `model/`, `route/`)
- Naming:
  - camelCase for private (unexported) identifiers
  - PascalCase for public (exported) identifiers
  - Short variable names in small scopes (`p` for player, `err` for error)
- Error Handling: Always check errors immediately after function calls
  ```go
  player, err := service.GetByID(id)
  if err != nil {
      return nil, err
  }
  ```
- Pointers: Use pointers for structs in function signatures to avoid copying

## Testing

- Test Structure: `*_test.go` files in the same package
- Naming: `Test*` for test functions, descriptive names (e.g., `TestGetPlayerByID`)
- Table-Driven Tests: Use for multiple test cases
- Coverage: Target 80%+ for service, controller, route packages

## Avoid

- Ignoring errors (using `_` discard) - always handle errors explicitly
- Panic in library code - return errors instead
- Global mutable state - pass dependencies explicitly
- Using `interface{}` without type assertions - prefer concrete types or generics
- Complex goroutines for simple CRUD operations - keep it simple

## Commits

- Format: Follow Conventional Commits with issue number suffix
- Pattern: `type(scope): description (#issue)` (max 80 chars)
- Examples: `feat(api): add player stats endpoint (#42)`, `fix(db): resolve connection pool leak (#88)`
