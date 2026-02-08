# AGENTS.md

> **⚡ Token Efficiency Note**: This file contains complete operational instructions (~2,500 tokens).
> **Auto-loaded**: NO (load explicitly with `#file:AGENTS.md` when you need detailed procedures)
> **When to load**: Complex workflows, troubleshooting, CI/CD setup, detailed architecture questions
> **Related files**: See `#file:.github/copilot-instructions.md` for quick context (auto-loaded, ~500 tokens)

---

## Quick Start

```bash
# Install dependencies
go mod download

# Run development server
go run main.go
# Server starts on http://localhost:9000

# View API documentation
# Open http://localhost:9000/swagger/index.html in browser
```

## Go Version

This project uses **Go 1.25** (specified in `go.mod`).

## Development Workflow

### Running Tests

```bash
# Run all tests with verbose output (recommended for local development)
go test -v ./...

# Run tests with coverage report (CI uses full-package coverage)
go test -v ./... \
  -coverpkg=github.com/nanotaboada/go-samples-gin-restful/service,github.com/nanotaboada/go-samples-gin-restful/controller,github.com/nanotaboada/go-samples-gin-restful/route \
  -covermode=atomic \
  -coverprofile=coverage.out

# View coverage in browser
go tool cover -html=coverage.out

# Run specific test function
go test -v ./... -run TestGetPlayers
```

**Coverage requirement**: Tests must maintain coverage. The CI pipeline enforces this.

### Code Quality

```bash
# Format code (auto-fix, must run before commit)
go fmt ./...

# Run linter (matches CI validation)
go vet ./...

# Check for common issues (recommended)
golangci-lint run  # if installed
```

**Pre-commit checklist**:

1. Update CHANGELOG.md `[Unreleased]` section with your changes (Added/Changed/Deprecated/Removed/Fixed/Security)
2. Run `go fmt ./...` - formats all Go files
3. Run `go vet ./...` - must pass with no warnings
4. Run `go test -v ./...` - all tests must pass

### Swagger Documentation

```bash
# Regenerate Swagger docs (after changing API annotations)
swag init

# This updates docs/ folder (DO NOT EDIT docs/ manually)
# Swagger comments are in:
# - swagger/swagger.go (API metadata)
# - controller/player_controller.go (endpoint annotations)
```

**Important**: The `docs/` folder is auto-generated. Only edit source files with `// @` annotations.

### Database Management

```bash
# Database auto-initializes on first app startup
# Pre-seeded database ships in storage/players.db

# To reset database to seed state (local development)
rm storage/players.db
# Next app startup will recreate via GORM AutoMigrate + seeding

# Database location: storage/players.db
```

**Important**: SQLite database with GORM ORM. Auto-migrates schema and seeds with football player data on first run.

## Docker Workflow

```bash
# Build container image
docker compose build

# Start application in container
docker compose up

# Start in detached mode (background)
docker compose up -d

# View logs
docker compose logs -f

# Stop application
docker compose down

# Stop and remove database volume (full reset)
docker compose down -v

# Health check (when running)
curl http://localhost:9000/health
```

**First run behavior**: Container initializes SQLite database with seed data. Volume persists data between runs.

## Release Management

### CHANGELOG Maintenance

**Important**: Update CHANGELOG.md continuously as you work, not just before releases.

**For every meaningful commit**:

1. Add your changes to the `[Unreleased]` section in CHANGELOG.md
2. Categorize under the appropriate heading:
   - **Added**: New features
   - **Changed**: Changes in existing functionality
   - **Deprecated**: Soon-to-be removed features
   - **Removed**: Removed features
   - **Fixed**: Bug fixes
   - **Security**: Security vulnerability fixes
3. Use clear, user-facing descriptions (not just commit messages)
4. Include PR/issue numbers when relevant (#123)

**Example**:

```markdown
## [Unreleased]

### Added
- User authentication with JWT tokens (#145)
- Rate limiting middleware for API endpoints

### Deprecated
- Legacy authentication endpoint /api/v1/auth (use /api/v2/auth instead)

### Fixed
- Null reference exception in player service (#147)

### Security
- Fix SQL injection vulnerability in search endpoint (#148)
```

### Creating a Release

When ready to release:

1. **Update CHANGELOG.md**: Move items from `[Unreleased]` to a new versioned section:

   ```markdown
   ## [1.1.0 - bobby] - 2026-02-15
   ```

2. **Commit and push** CHANGELOG changes
3. **Create and push tag**:

   ```bash
   git tag -a v1.1.0-bobby -m "Release 1.1.0 - Bobby"
   git push origin v1.1.0-bobby
   ```

4. **CD workflow runs automatically** to publish Docker images and create GitHub Release

See [CHANGELOG.md](CHANGELOG.md#how-to-release) for complete release instructions and player naming convention.

## CI/CD Pipeline

### Continuous Integration (go.yml)

**Trigger**: Push to `main`/`master` or PR

**Jobs**:

1. **Setup**: Go 1.25 installation, dependency caching
2. **Format Check**: `go fmt` validation
3. **Lint**: `go vet` checks
4. **Build**: `go build -v ./...`
5. **Test**: `go test -v ./...` with `-race`, `-coverpkg` (service/controller/route), and coverage output
6. **Coverage**: Upload to Codecov

**Local validation** (run this before pushing):

```bash
# Quick local validation (recommended for most development)
go fmt ./... && \
go vet ./... && \
go build -v ./... && \
go test -v ./...

# Full CI validation (exact match - use when you need complete verification)
go fmt ./... && \
go vet ./... && \
go build -v ./... && \
go test -v ./... \
  -coverpkg=github.com/nanotaboada/go-samples-gin-restful/service,github.com/nanotaboada/go-samples-gin-restful/controller,github.com/nanotaboada/go-samples-gin-restful/route \
  -covermode=atomic \
  -coverprofile=coverage.out
```

## Project Architecture

**Structure**: Layered architecture (Controller → Service → Data)

```tree
main.go              # Entry point: DB init, router setup, server start

controller/          # HTTP handlers
  └── player_controller.go  # Request/response with Swagger annotations

service/             # Business logic
  └── player_service.go     # GORM operations (CRUD)

route/               # Router configuration
  ├── player_route.go        # Route registration with caching middleware
  └── path.go                # Path constants

model/               # Data structures
  └── player_model.go        # Player struct (GORM model)

data/                # Database layer
  └── player_data.go         # DB connection, GORM setup

swagger/             # Swagger config
  └── swagger.go             # API metadata annotations

docs/                # Auto-generated (DO NOT EDIT)
  ├── docs.go
  ├── swagger.json
  └── swagger.yaml

tests/               # Integration tests
  ├── main_test.go           # Endpoint tests with testify
  ├── player_fake.go         # Test fixtures
  └── players.json           # Seed data
```

**Key patterns**:

- GORM ORM for database operations
- Gin middleware for routing, caching, CORS
- Swagger annotations for API docs (`// @` comments)
- In-memory caching with 1-hour TTL
- Thread-safe GORM operations

## Troubleshooting

### Port already in use

```bash
# Kill process on port 9000
lsof -ti:9000 | xargs kill -9
```

### Module dependency errors

```bash
# Clean module cache and reinstall
go clean -modcache
go mod download
go mod tidy  # Cleanup unused dependencies
```

### Database locked errors

```bash
# Stop all running instances
pkill -f "go run main.go"

# Reset database
rm storage/players.db
```

### Swagger not updating

```bash
# Ensure swag is installed
go install github.com/swaggo/swag/cmd/swag@latest

# Regenerate docs
swag init

# Restart server to load new docs
```

### Build failures

```bash
# Verify Go version
go version  # Should be 1.25

# Clean build cache
go clean -cache

# Rebuild
go build -v ./...
```

### Docker issues

```bash
# Clean slate
docker compose down -v
docker compose build --no-cache
docker compose up
```

## Testing the API

### Using Swagger UI (Recommended)

Open <http://localhost:9000/swagger/index.html> - Interactive documentation with "Try it out"

### Using Postman

Pre-configured collection available in `postman-collections/`

### Using curl

```bash
# Health check
curl http://localhost:9000/health

# Get all players
curl http://localhost:9000/api/v1/players

# Get player by ID
curl http://localhost:9000/api/v1/players/1

# Create player
curl -X POST http://localhost:9000/api/v1/players \
  -H "Content-Type: application/json" \
  -d '{"firstName":"Pele","lastName":"Nascimento","club":"Santos","nationality":"Brazil","dateOfBirth":"1940-10-23","squadNumber":10}'

# Update player
curl -X PUT http://localhost:9000/api/v1/players/1 \
  -H "Content-Type: application/json" \
  -d '{"firstName":"Diego","lastName":"Maradona","club":"Napoli","nationality":"Argentina","dateOfBirth":"1960-10-30","squadNumber":10}'

# Delete player
curl -X DELETE http://localhost:9000/api/v1/players/1
```

## Important Notes

- **CHANGELOG maintenance**: Update CHANGELOG.md `[Unreleased]` section with every meaningful change
- **Never commit secrets**: No API keys, tokens, or credentials in code
- **Test coverage**: Maintain existing coverage levels
- **Commit messages**: Follow conventional commits (enforced by commitlint)
- **Go version**: Must use 1.25 for consistency with CI/CD
- **Swagger annotations**: Required for all new endpoints
- **Database**: SQLite is for demo/development only - not production-ready
- **Module management**: Use `go mod tidy` to keep dependencies clean
