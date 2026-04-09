# Custom Instructions

## Overview

REST API for managing football players built with Go and Gin Web Framework. Implements CRUD operations with SQLite + GORM, in-memory caching, and Swagger documentation. Part of a cross-language comparison study (.NET, Java, Python, Rust, TypeScript). Architectural decisions are documented as ADRs in `docs/adr/` — check them before proposing structural changes.

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
/docs/adr       — Architecture Decision Records (read before proposing structural changes)
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
- **Tests**: Table-driven tests for multiple cases; target 80%+ coverage for service, controller, route packages
- **Test strategy**: Integration tests with real in-memory SQLite for all happy paths and expected branches. Use `MockPlayerService` only for error branches that cannot be triggered with a healthy database (e.g. simulated connection failures). If a scenario can be exercised with a real database, it must use a real database.
- **Mock pattern**: `MockPlayerService` uses opt-in function fields — only set the `Func` relevant to the test scenario; unset methods return safe zero-value defaults. Never create a new mock type per test.
- **Test naming**: `TestRequest{METHOD}{Resource}{Condition}Response{Outcome}`:
  - **Resource**: explicit endpoint target — `Players`, `PlayerByID`, `PlayerBySquadNumber`
  - **Condition**: `Existing`, `NonExisting`, `InvalidParam`, `EmptyBody`, `TrailingSlash`, `RetrieveError`, `CreateError`, `UpdateError`, `DeleteError`
  - **Outcome**: `StatusOK`, `StatusCreated`, `StatusNoContent`, `StatusBadRequest`, `StatusNotFound`, `StatusConflict`, `StatusInternalServerError`, or `Players` / `Player` for body assertions
  - Examples: `TestRequestGETPlayerByIDExistingResponseStatusOK`, `TestRequestPOSTPlayersEmptyBodyResponseStatusBadRequest`, `TestRequestDELETEPlayerByIDDeleteErrorResponseStatusInternalServerError`
- **Test godoc**: Each `Test*` function must open with: `// TestFuncName tests that a\n// {METHOD} request to {/path} {condition}\n// returns a {outcome}.`
- **Avoid**: ignoring errors, `panic` in library code, global mutable state, `interface{}` without type assertions, complex goroutines for simple CRUD

## Commands

### Quick Start

```bash
go mod download
go run .            # starts on port 9000 (set STORAGE_PATH to override DB location)
go build -v ./...
go test ./...       # all tests
go test -v ./... -coverpkg=github.com/nanotaboada/go-samples-gin-restful/service,github.com/nanotaboada/go-samples-gin-restful/controller,github.com/nanotaboada/go-samples-gin-restful/route -covermode=atomic -coverprofile=coverage.out
go tool cover -html=coverage.out
swag init           # regenerate Swagger docs
docker compose up --build
```

**Environment variables:**
- `STORAGE_PATH` — path to the SQLite database file. Defaults to `./storage/players-sqlite3.db` when unset (local development). Set by Docker Compose to `/storage/players-sqlite3.db` (persistent volume).
- `GIN_MODE` — `debug` (default locally) or `release` (set by Docker Compose).

### Pre-commit Checks

1. Update `CHANGELOG.md` `[Unreleased]` section (Added / Changed / Fixed / Removed)
2. `go fmt ./...`
3. `go build -v ./...`
4. If Swagger annotations were modified: `swag init`
5. `go test ./...` — all tests must pass
6. Full coverage command above — target 80%+ for service, controller, route
7. `golangci-lint run`
8. Verify all errors explicitly checked; JSON struct tags present on model structs
9. Commit message follows Conventional Commits format (enforced by commitlint)

### Commits

Format: `type(scope): description (#issue)` — max 80 chars
Types: `feat` `fix` `chore` `docs` `test` `refactor` `ci` `perf`
Example: `feat(api): add player stats endpoint (#42)`

### Releases

Tags follow the format `v{SEMVER}-{PLAYER}` (e.g. `v2.0.0-bobby`). The CD workflow validates the player name against the 26-name list in `CHANGELOG.md` and rejects unknown names. Never suggest a release tag without confirming the player name is in that list.

## Agent Mode

### Proceed freely

- Route handlers and controllers
- Service layer logic and validation
- Unit and integration tests
- Refactoring within controller/service layers
- Documentation updates and bug fixes
- Utility functions and helpers

### Ask before changing

- Database schema (`Player` struct fields) — schema changes require a new goose migration file; breaking changes (column type or column removal) also require updating the seed migrations and must be flagged explicitly
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

### Creating Issues

This project uses Spec-Driven Development (SDD): discuss in Plan mode first, create a GitHub Issue as the spec artifact, then implement. Always offer to draft an issue before writing code.

**Feature request** (`enhancement` label):
- **Problem**: the pain point being solved
- **Proposed Solution**: expected behavior and functionality
- **Suggested Approach** *(optional)*: implementation plan if known
- **Acceptance Criteria**: at minimum — behaves as proposed, tests added/updated, no regressions
- **References**: related issues, docs, or examples

**Bug report** (`bug` label):
- **Description**: clear summary of the bug
- **Steps to Reproduce**: numbered, minimal steps
- **Expected / Actual Behavior**: one section each
- **Environment**: runtime versions + OS
- **Additional Context**: logs, screenshots, stack traces
- **Possible Solution** *(optional)*: suggested fix or workaround

### Key workflows

**Add an endpoint**: Define model in `/model/` (if needed) → add service method in `/service/` → create controller handler in `/controller/` → register route in `/route/` → add Swagger comments → add tests → run `swag init` → run pre-commit checks.

**Modify schema**: Update `Player` struct → add a new goose migration in `/migrations/` → update GORM queries in `/service/` → update controller handlers → fix test assertions → run `swag init` → run `go test ./...`. If the change is breaking (column type or column removal), also update the seed migrations and warn that the existing database must be re-created (`goose down` then `goose up`, or `docker compose down -v`).

**After completing work**: Suggest a branch name (e.g. `feat/add-player-stats`) and a commit message following Conventional Commits including co-author line:

```text
feat(scope): description (#issue)

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>
```
