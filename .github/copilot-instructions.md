# GitHub Copilot Instructions

## Overview

REST API for managing football players built with Go and Gin Web Framework. Implements CRUD operations with SQLite + GORM, in-memory caching, and Swagger documentation. Part of a cross-language comparison study (.NET, Java, Python, Rust, TypeScript).

## Tech Stack

- **Language**: Go 1.25+
- **Framework**: Gin Web Framework
- **ORM**: GORM
- **Database**: SQLite
- **Caching**: gin-contrib/cache (in-memory, 1-hour TTL)
- **Testing**: Go testing package + testify/assert
- **Linting**: golangci-lint
- **API Docs**: Swaggo (Swagger generated from comments)
- **Containerization**: Docker

## Structure

```text
main.go         — application entry point: Gin setup, DB init, route registration
go.mod          — module dependencies
/route          — route registration + caching middleware       [HTTP layer]
/controller     — HTTP handlers; request/response logic         [HTTP layer]
/service        — business logic + GORM interactions            [business layer]
/data           — database connection setup                     [data layer]
/model          — Player struct (domain model)
/storage        — SQLite database file (players-sqlite3.db, pre-seeded)
/docs           — auto-generated Swagger docs (DO NOT EDIT manually)
/tests          — integration tests with testify assertions
/scripts        — Docker entrypoint and healthcheck scripts
```

**Layer rule**: `Routes → Controllers → Services → Data`. Never skip a layer. Controllers must not contain business logic.

## Coding Guidelines

- **Naming**: camelCase (unexported), PascalCase (exported), short names in small scopes
- **Files**: snake_case for all file names
- **Errors**: Always check errors immediately after function calls; never discard with `_`
- **Pointers**: Use pointers for structs in function signatures to avoid copying
- **Logging**: Standard `log` package (structured `slog` for new code)
- **Tests**: Table-driven tests for multiple cases; `Test*` naming convention; target 80%+ coverage for service, controller, route packages
- **Avoid**: ignoring errors, `panic` in library code, global mutable state, `interface{}` without type assertions, complex goroutines for simple CRUD

## Commands

### Quick Start

```bash
go mod download
go run .            # starts on port 9000
go build -v ./...
go test ./...       # all tests
go test -v ./... -coverpkg=github.com/nanotaboada/go-samples-gin-restful/service,github.com/nanotaboada/go-samples-gin-restful/controller,github.com/nanotaboada/go-samples-gin-restful/route -covermode=atomic -coverprofile=coverage.out
go tool cover -html=coverage.out
swag init           # regenerate Swagger docs
docker compose up --build
```

### Pre-commit Checks

1. Update `CHANGELOG.md` `[Unreleased]` section (Added / Changed / Fixed / Removed)
2. `go fmt ./...`
3. `go build`
4. `go test ./...` — all tests must pass
5. Full coverage command above — target 80%+ for service, controller, route
6. `golangci-lint run`
7. Verify all errors explicitly checked; JSON struct tags present on model structs
8. Commit message follows Conventional Commits format (enforced by commitlint)

### Commits

Format: `type(scope): description (#issue)` — max 80 chars
Types: `feat` `fix` `chore` `docs` `test` `refactor` `ci` `perf`
Example: `feat(api): add player stats endpoint (#42)`

## Agent Mode

### Proceed freely

- Route handlers and controllers
- Service layer logic and validation
- Unit and integration tests
- Refactoring within controller/service layers
- Documentation updates and bug fixes
- Utility functions and helpers

### Ask before changing

- Database schema (`Player` struct fields)
- Dependencies (`go.mod`)
- CI/CD configuration (`.github/workflows/`)
- Docker setup
- Gin middleware or router configuration
- HTTP status codes or error response formats
- Package organization

### Never modify

- `go.mod` module path
- Port configuration (9000)
- Database type (SQLite)
- Auto-generated Swagger docs in `/docs` (run `swag init` instead)

### Key workflows

**Add an endpoint**: Define model in `/model/` (if needed) → add service method in `/service/` → create controller handler in `/controller/` → register route in `/route/` → add Swagger comments → add tests → run `swag init` → run pre-commit checks.

**Modify schema**: Update `Player` struct → update GORM queries in `/service/` → update controller handlers → update `/tests/players.json` → fix test assertions → run `swag init` → run `go test ./...`.

**After completing work**: Suggest a branch name (e.g. `feat/add-player-stats`) and a commit message following Conventional Commits including co-author line:

```text
feat(scope): description (#issue)

Co-authored-by: Copilot <175728472+Copilot@users.noreply.github.com>
```
