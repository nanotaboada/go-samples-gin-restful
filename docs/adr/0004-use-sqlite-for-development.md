# ADR-0004: Use SQLite for Development

Date: 2026-03-21

## Status

Accepted

## Context

Production Go APIs typically use PostgreSQL or MySQL. Running a database server locally requires Docker, port allocation, and credential management — a meaningful barrier in environments where Docker is unavailable or restricted (legacy CI systems, locked-down corporate machines).

This project is part of a cross-language PoC suite intended to run out of the box with a single command (`go run .`). Minimizing infrastructure prerequisites is an explicit goal.

Alternatives evaluated:

- **PostgreSQL (local)**: Requires a running server; `docker compose up` or a local install.
- **PostgreSQL (Docker only)**: Enforces Docker as a hard dependency — disqualifying for restricted environments.
- **MySQL / MariaDB**: Same infrastructure concerns as PostgreSQL.
- **SQLite**: File-based, zero server process, embeds directly in the Go binary via CGo.

## Decision

We will use SQLite as the database engine. The pre-seeded database file (`storage/players-sqlite3.db`) is committed to the repository so the API is operational immediately after `go run .` with no external dependencies. Tests use in-memory SQLite (`file::memory:?cache=shared`), eliminating test database setup entirely.

PostgreSQL remains the intended production target. GORM's driver abstraction makes the migration primarily a configuration change — swap `gorm.io/driver/sqlite` for `gorm.io/driver/postgres` and update the DSN.

## Consequences

**Positive:**

- Zero infrastructure: no Docker, no server process, no port conflicts, no credentials.
- The pre-seeded file provides a working 23-player dataset immediately after cloning.
- In-memory mode enables fast, isolated integration tests without test containers or teardown logic.
- GORM's driver abstraction means the service and controller layers are unaware of the underlying database engine.

**Negative:**

- `isUniqueConstraintError()` matches the SQLite-specific string `"UNIQUE constraint failed"` — this function must be updated when switching to PostgreSQL (PostgreSQL uses a different error format).
- SQLite's single-writer concurrency limitation is invisible during development; workloads that would deadlock under concurrent writes will not surface locally.
- Some PostgreSQL-specific features (advisory locks, `LISTEN`/`NOTIFY`, advanced types, full-text search) are unavailable and cannot be prototyped against SQLite.
