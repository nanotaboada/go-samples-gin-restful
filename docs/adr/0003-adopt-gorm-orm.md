# ADR-0003: Adopt GORM as ORM

Date: 2026-03-21

## Status

Accepted

## Context

Data access options in Go range from raw SQL to full ORM:

- **database/sql**: Standard library, full control, maximum boilerplate — every query written and scanned manually.
- **sqlx**: Thin wrapper over `database/sql`, adds struct scanning, still requires manual SQL.
- **ent**: Code-generation-based, strongly typed query builder — higher setup cost, generated files in version control.
- **bun**: Modern SQL-first ORM, less mature ecosystem than GORM.
- **GORM**: Convention-over-configuration, struct tag mapping, `AutoMigrate`, implicit soft-delete support.

For a learning-focused PoC, reducing schema management friction and making data flow visible in struct tags were primary goals.

## Decision

We will use GORM v2 for all database interactions. The `Player` struct carries both GORM column mappings (`gorm:` tags) and JSON serialization tags (`json:` tags), making the database mapping explicit alongside the domain model. `AutoMigrate` is called at startup to keep the schema synchronized with the struct without manual DDL.

The service layer exposes a `PlayerService` interface, not the concrete GORM-backed struct, allowing tests to inject a mock without a GORM test double.

## Consequences

**Positive:**

- `AutoMigrate` eliminates manual `CREATE TABLE` DDL for a single-entity project.
- Struct tags make column mapping readable alongside the model — the full data contract is visible in one place.
- `gorm.ErrRecordNotFound` provides a typed sentinel for 404 detection via `errors.Is`.
- The interface-based service layer decouples the HTTP layer from GORM entirely.

**Negative:**

- `AutoMigrate` cannot handle breaking schema changes (e.g., changing `id` from integer to UUID string required manual re-seeding). Documented in `data/player_data.go`.
- GORM's implicit behaviors (soft delete via `gorm.DeletedAt`, auto-timestamps via `gorm.Model`) require explicit opt-in/opt-out awareness per struct.
- Switching to PostgreSQL requires updating the driver import and verifying SQLite-specific error handling — `isUniqueConstraintError()` matches the string `"UNIQUE constraint failed"`, which is SQLite-specific (see ADR-0004).
