# Agent Instructions

This file provides autonomous coding agent guidance for workflows, decision-making authority, and validation steps. For stack details, build commands, and conventions, see `.github/copilot-instructions.md` (auto-loaded).

## Common Workflows

### Adding a new endpoint

1. Define model (if needed) in `/model/player_model.go`
2. Add service method in `/service/player_service.go` (business logic with GORM)
3. Create controller handler in `/controller/player_controller.go` (HTTP request/response)
4. Register route in `/route/player_route.go` (with caching if needed)
5. Add Swagger comments to controller handler
6. Add tests in `/tests/main_test.go`
7. Run `swag init` to regenerate Swagger docs
8. Validate: `go test ./...` passes

### Modifying database schema

1. Update `Player` struct in `/model/player_model.go` (add/modify fields)
2. Update GORM queries in `/service/player_service.go` (adjust to new schema)
3. Update controller handlers in `/controller/player_controller.go` (handle new fields)
4. Regenerate database: Either manually update `/storage/players-sqlite3.db` or let GORM auto-migrate
5. Update test data in `/tests/players.json` (add new fields)
6. Fix test assertions in `/tests/main_test.go`
7. Update Swagger comments and run `swag init`
8. Validate: `go test ./...` passes with updated assertions

### Adding tests

1. Add test cases to `/tests/main_test.go` (uses testify assertions)
2. Use fake data from `/tests/player_fake.go` or `/tests/players.json`
3. Mock service methods in `/tests/mock_service.go` if needed
4. Run tests: `go test -v ./...`
5. Check coverage: `go test -v ./... -coverpkg=github.com/nanotaboada/go-samples-gin-restful/service,github.com/nanotaboada/go-samples-gin-restful/controller,github.com/nanotaboada/go-samples-gin-restful/route -covermode=atomic -coverprofile=coverage.out`
6. View report: `go tool cover -html=coverage.out`
7. Target: 80%+ for service, controller, route packages

## Autonomy Levels

### Proceed freely

- Route handlers and controllers (following existing patterns)
- Service layer logic and validation
- Unit and integration tests
- Refactoring within controller/service layers
- Documentation updates (comments, README, Swagger)
- Bug fixes in business logic
- Utility functions and helpers

### Ask before changing

- Database schemas (Player struct fields)
- Dependencies (`go.mod` - adding new modules)
- CI/CD configuration (`.github/workflows/`)
- Docker setup (`Dockerfile`, `compose.yaml`)
- Environment variable requirements
- Gin middleware or router configuration
- HTTP status codes or error response formats
- Adding goroutines or concurrency patterns
- Package organization restructuring

### Never modify

- Production configurations
- Deployment secrets
- API contracts (request/response structures without discussion)
- Existing endpoints (deletion)
- `go.mod` module path
- Port configuration (9000)
- Database type (SQLite)
- Error handling patterns (always explicit)

## Pre-commit Checks

1. `go fmt ./...` - Code formatting (should be automatic in editor)
2. `go build` - Compiles without errors
3. `go test ./...` - All tests pass
4. `go test -v ./... -coverpkg=github.com/nanotaboada/go-samples-gin-restful/service,github.com/nanotaboada/go-samples-gin-restful/controller,github.com/nanotaboada/go-samples-gin-restful/route -covermode=atomic -coverprofile=coverage.out` - Coverage 80%+ (service, controller, route packages)
5. `golangci-lint run` - Linting passes (if configured)
6. All errors explicitly checked (no ignored error returns)
7. JSON struct tags present on model structs
8. Only necessary identifiers exported (PascalCase for public, camelCase for private)
