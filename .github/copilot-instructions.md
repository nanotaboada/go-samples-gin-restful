# Copilot Instructions: go-samples-gin-restful

## Project Overview

This is a **RESTful API proof-of-concept** built with **Go 1.25** and the **Gin web framework**. The application manages a collection of football players with full CRUD operations, Swagger documentation, and in-memory caching. It uses **SQLite** for data persistence and is designed for containerized deployment with **Docker**.

## Tech Stack

- **Language**: Go 1.25
- **Web Framework**: Gin (`github.com/gin-gonic/gin`)
- **ORM**: GORM v1.31.1 (`gorm.io/gorm`)
- **Database**: SQLite (`gorm.io/driver/sqlite`)
- **API Documentation**: Swagger/OpenAPI (`github.com/swaggo/swag`, `github.com/swaggo/gin-swagger`)
- **Caching**: In-memory cache (`github.com/gin-contrib/cache`)
- **Testing**: testify (`github.com/stretchr/testify`)
- **Containerization**: Docker & Docker Compose

## Architecture & Project Structure

The project follows a **layered architecture** with clear separation of concerns:

```
/
├── main.go                 # Entry point: DB connection, route setup, server start
├── controller/             # HTTP handlers (request/response logic)
│   └── player_controller.go
├── service/                # Business logic (ORM interactions)
│   └── player_service.go
├── route/                  # Route configuration and middleware
│   ├── player_route.go     # Route setup with caching middleware
│   └── path.go             # Path constants
├── model/                  # Data structures
│   └── player_model.go
├── data/                   # Database connection
│   └── player_data.go
├── swagger/                # Swagger configuration
│   └── swagger.go
├── docs/                   # Auto-generated Swagger docs (DO NOT EDIT)
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── tests/                  # Integration tests
│   ├── main_test.go
│   ├── player_fake.go
│   └── players.json
├── storage/                # SQLite database file (pre-seeded)
├── scripts/                # Container entrypoint & healthcheck
└── .github/workflows/      # CI/CD pipelines
```

**Data Flow**: HTTP Request → Controller (validation, status codes) → Service (business logic) → Data (ORM/DB) → Response

**Dependency Injection**: main → data.Connect() → service.NewPlayerService(db) → controller.NewPlayerController(service) → route.Setup(controller)

## Coding Guidelines

### General Principles

1. **Incremental Changes**: Follow the philosophy in `CONTRIBUTING.md` — start small, solve immediate needs, avoid over-design.
2. **Conventional Commits**: All commits must follow [Conventional Commits](https://www.conventionalcommits.org/): `feat:`, `fix:`, `chore:`, etc.
3. **Commit Message Limits**: Max 80 characters for header and body lines (enforced by `commitlint.config.mjs`).

### Go Conventions

- **Package Documentation**: Every package must have a doc comment (see existing files).
- **Function Documentation**: Public functions, especially controller endpoints, must have doc comments (used for Swagger generation).
- **Error Handling**: Always check and handle errors. Use appropriate HTTP status codes.
- **Naming**: Use idiomatic Go naming (camelCase for unexported, PascalCase for exported).
- **Imports**: Group standard library, then third-party, then local packages.

### Swagger/OpenAPI Annotations

Controller functions require Swagger annotations:
```go
// FunctionName does something
//
// @Summary Brief description
// @Tags tagname
// @Accept application/json
// @Produce application/json
// @Param paramname path/query/body type true "Description"
// @Success 200 {object} model.TypeName "Description"
// @Failure 400 "Description"
// @Router /path [method]
```

**Important**: After modifying Swagger annotations, regenerate docs with:
```bash
swag init
```
This updates `docs/docs.go`, `docs/swagger.json`, and `docs/swagger.yaml`.

### Testing Standards

- **Test Structure**: Follow Given-When-Then pattern with AAA (Arrange-Act-Assert).
- **Test Naming**: `Test<Method><Scenario>Response<Expected>` (e.g., `TestRequestPOSTBodyEmptyResponseStatusBadRequest`).
- **Coverage Targets**: 80% for services, controllers, and routes (enforced by `codecov.yml`).
- **Test Database**: Uses in-memory SQLite (`file::memory:?cache=shared`) with test data from `tests/players.json`.

## Development Workflow

### Running Locally

```bash
# Start the server (port 9000)
go run .

# View Swagger documentation
open http://localhost:9000/swagger/index.html
```

The app reads the database from `./storage/players-sqlite3.db` unless `STORAGE_PATH` env var is set (Docker mode).

### Building & Testing

```bash
# Install dependencies
go mod tidy
go get .

# Build
go build -v ./...

# Run tests with coverage
go test -v ./... \
  -coverpkg=github.com/nanotaboada/go-samples-gin-restful/service,github.com/nanotaboada/go-samples-gin-restful/controller,github.com/nanotaboada/go-samples-gin-restful/route \
  -covermode=atomic \
  -coverprofile=coverage.out
```

**Note**: CGO is required for SQLite support. On Alpine/Docker, ensure `gcc`, `musl-dev`, and `sqlite-dev` are installed.

### Docker Workflow

```bash
# Build image
docker compose build

# Start container (first run initializes DB from seed)
docker compose up

# Stop container
docker compose down

# Reset database (removes volume, reinitializes on next `up`)
docker compose down -v
```

**How Docker Persistence Works**:
- Pre-seeded database is copied into image at `/app/hold/players-sqlite3.db`
- On first run, `entrypoint.sh` copies it to the persistent volume at `/storage/players-sqlite3.db`
- On subsequent runs, the volume is reused (data persists)

## Common Patterns & Best Practices

### Adding a New Endpoint

1. **Define the route constant** in `route/path.go`.
2. **Create the controller function** in `controller/` with Swagger annotations.
3. **Implement service logic** in `service/` if needed.
4. **Register the route** in `route/player_route.go`:
  - Use `cache.CachePage()` for GET requests.
  - Use `ClearCache()` wrapper for POST/PUT/DELETE.
5. **Regenerate Swagger docs**: `swag init`.
6. **Write tests** in `tests/main_test.go` following the existing pattern.

### Cache Invalidation

All GET endpoints use in-memory caching (1 hour TTL). POST/PUT/DELETE operations automatically clear relevant cache keys via the `ClearCache()` middleware in `route/player_route.go`.

### Database Operations

Always use GORM methods in the `service/` layer. Do NOT write raw SQL queries. Refer to [GORM documentation](https://gorm.io/docs/) for query patterns.

## CI/CD Pipeline

The GitHub Actions workflow (`.github/workflows/go.yml`) runs on every push/PR:

1. **Build**: Compiles all packages.
2. **Test**:
  - Lints commit messages with `commitlint`.
  - Runs tests with coverage.
  - Uploads coverage to Codecov & Codacy.
3. **Container**: On `master` branch pushes, builds and pushes Docker image to GitHub Container Registry.

**Environment Variables Used**:
- `GO_VERSION`: 1.25.0
- `PKG_SERVICE`, `PKG_CONTROLLER`, `PKG_ROUTE`: Coverage target packages

## Key Files & Configurations

| File | Purpose |
|------|---------|
| `go.mod`, `go.sum` | Go module dependencies (use `go mod tidy` to update) |
| `commitlint.config.mjs` | Commit message format enforcement |
| `codecov.yml` | Code coverage requirements (80% target) |
| `compose.yaml` | Docker Compose service definition |
| `Dockerfile` | Multi-stage build (builder + runtime) |
| `scripts/entrypoint.sh` | Initializes persistent volume on first run |
| `scripts/healthcheck.sh` | Container health check (`curl http://localhost:9000/health`) |
| `storage/players-sqlite3.db` | Pre-seeded SQLite database |
| `tests/players.json` | Test fixture data |

## Common Pitfalls & Solutions

1. **Swagger docs out of sync**: Always run `swag init` after changing controller annotations.
2. **CGO errors in Docker**: Ensure `CGO_ENABLED=1` and build deps (`gcc`, `musl-dev`, `sqlite-dev`) are installed.
3. **Test failures after DB changes**: Update `tests/players.json` fixture or `tests/player_fake.go` accordingly.
4. **Cache not clearing**: Ensure POST/PUT/DELETE routes use the `ClearCache()` wrapper.
5. **Port conflicts**: Default port is 9000. Check for conflicts with `lsof -i :9000`.

## Resources

- **Gin Documentation**: https://gin-gonic.com/docs/
- **GORM Documentation**: https://gorm.io/docs/
- **Swagger Annotations**: https://github.com/swaggo/swag#general-api-info
- **Conventional Commits**: https://www.conventionalcommits.org/

## Quick Command Reference

```bash
# Development
go run .                    # Start server
go test ./...               # Run all tests
go mod tidy                 # Clean dependencies
swag init                   # Regenerate Swagger docs

# Docker
docker compose build        # Build image
docker compose up           # Start container
docker compose down -v      # Stop and reset DB

# Quality Checks
go fmt ./...                # Format code
go vet ./...                # Static analysis
go test -v ./... -coverprofile=coverage.out
```
