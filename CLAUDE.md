# CLAUDE.md

## Overview

REST API for managing football players built with Go and Gin Web Framework. Implements CRUD operations with SQLite + GORM, in-memory caching, and Swagger documentation. Architectural decisions are documented as ADRs in `docs/adr/` — check them before proposing structural changes.

## Structure

**Layer rule**: `Routes → Controllers → Services → Data`. Never skip a layer. Controllers must not contain business logic.

## Coding Guidelines

- **Pointers**: Use pointers for structs in function signatures to avoid copying
- **Logging**: Standard `log` package (structured `slog` for new code)
- **Migrations**: `migrations/embed.go` embeds all `.sql` files into the binary at compile time (`//go:embed *.sql`). No migration files are needed on the filesystem at runtime. Migration files use 5-digit zero-padded names (`00001_`, `00002_`).
- **Seed tools**: scripts in `tools/` use `//go:build ignore` and are excluded from normal builds. Run individually with `go run ./tools/seed_001_starting_eleven.go` (recreates the DB from scratch).
- **Tests**: Table-driven tests for multiple cases; target 80%+ coverage for service, controller, route packages
- **Test strategy**: Integration tests with real in-memory SQLite for all happy paths and expected branches. Use `MockPlayerService` only for error branches that cannot be triggered with a healthy database (e.g. simulated connection failures). If a scenario can be exercised with a real database, it must use a real database.
- **Mock pattern**: `MockPlayerService` uses opt-in function fields — only set the `Func` relevant to the test scenario; unset methods return safe zero-value defaults. Never create a new mock type per test.
- **Test naming**: `TestRequest{METHOD}{Resource}{Condition}Response{Outcome}`:
  - **Resource**: explicit endpoint target — `Players`, `PlayerByID`, `PlayerBySquadNumber`
  - **Condition**: `Existing`, `NonExisting`, `Unknown`, `InvalidParam`, `Mismatch`, `EmptyBody`, `TrailingSlash`, `Validation`, `RetrieveError`, `CreateError`, `UpdateError`, `DeleteError`
  - **Outcome**: `StatusOK`, `StatusCreated`, `StatusNoContent`, `StatusBadRequest`, `StatusNotFound`, `StatusConflict`, `StatusUnprocessableEntity`, `StatusInternalServerError`, or `Players` / `Player` for body assertions
  - Examples: `TestRequestGETPlayerByIDExistingResponseStatusOK`, `TestRequestPOSTPlayersEmptyBodyResponseStatusBadRequest`, `TestRequestDELETEPlayerByIDDeleteErrorResponseStatusInternalServerError`
- **Test godoc**: Each `Test*` function must open with: `// TestFuncName tests that a\n// {METHOD} request to {/path} {condition}\n// returns a {outcome}.`
- **Avoid**: ignoring errors, `panic` in library code, global mutable state, `interface{}` without type assertions, complex goroutines for simple CRUD

## Commands

### Quick Start

```bash
go run .            # starts on port 9000 (set STORAGE_PATH to override DB location)
go test -v ./... -coverpkg=github.com/nanotaboada/go-samples-gin-restful/service,github.com/nanotaboada/go-samples-gin-restful/controller,github.com/nanotaboada/go-samples-gin-restful/route -covermode=atomic -coverprofile=coverage.out
go tool cover -html=coverage.out
```

**Environment variables:**
- `STORAGE_PATH` — path to the SQLite database file. Defaults to `./storage/players-sqlite3.db` when unset (local development). Set by Docker Compose to `/storage/players-sqlite3.db` (persistent volume).
- `GIN_MODE` — `debug` (default locally) or `release` (set by Docker Compose).

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

This project uses Spec-Driven Development (SDD): discuss in Plan mode first, create a GitHub Issue as the spec artifact, then implement. Always offer to draft an issue before writing code. See `.claude/skills/create-issue/SKILL.md` for the feature/bug issue templates.

### Key workflows

**Add an endpoint**: Define model in `/model/` (if needed) → add service method in `/service/` → create controller handler in `/controller/` → register route in `/route/` → add Swagger comments → add tests → run `swag init` → run pre-commit checks.

**Modify schema**: Update `Player` struct → add a new goose migration in `/migrations/` → update GORM queries in `/service/` → update controller handlers → fix test assertions → run `swag init` → run `go test ./...`. If the change is breaking (column type or column removal), also update the seed migrations and warn that the existing database must be re-created (`goose down` then `goose up`, or `docker compose down -v`).

**After completing work**: Suggest a branch name (e.g. `feat/add-player-stats`) and a commit message following Conventional Commits including co-author line:

```text
feat(scope): description (#issue)

Co-Authored-By: Claude Sonnet 5 <noreply@anthropic.com>
```

## Invariants (never change without explicit discussion)

- **Port**: 9000
- **API contract**: endpoints, HTTP status codes, and response shapes are fixed; do not change them without explicit discussion
- **Commit format**: `type(scope): description (#issue)` — max 80 chars
- **Conventional Commits types**: `feat` `fix` `chore` `docs` `test` `refactor` `ci` `perf`
- **CHANGELOG.md** `[Unreleased]` section must be updated before every commit that has a backing GitHub issue; changes with no issue behind them don't need an entry

## Architecture Decision Records

Significant architectural decisions are documented in `docs/adr/` (ADR-0001–ADR-0015). Load these before proposing structural changes. When a proposal would change an accepted decision, create a new ADR rather than editing the existing one.

## Claude Code

- Run `/pre-commit` to execute the full pre-commit checklist for this project.
- CLAUDE.md is maintained with the [CLAUDE.md Management plugin](https://claude.com/plugins/claude-md-management).
