# ADR-0010: Mixed Test Strategy (Integration Tests + Mock Injection)

Date: 2026-03-21

## Status

Accepted

## Context

Testing strategies for a layered HTTP API:

- **Fully mocked unit tests**: Fast and isolated. Mocks can diverge from real behavior; do not catch database-specific issues (constraint violations, GORM edge cases).
- **Fully integrated tests (real DB)**: High confidence. Require a running database; slower. With SQLite in-memory mode, setup cost is near zero.
- **Test containers**: Real database in Docker per test run. High confidence with production-equivalent behavior, but introduces infrastructure dependency.
- **Contract tests**: Verify interface boundaries between layers. Add complexity without clear benefit for a single-service project.

Two distinct categories of test scenario exist in this codebase:

1. **Happy paths and expected branches** (200, 201, 204, 400, 404, 409): These can be exercised reliably with a real in-memory SQLite database seeded via goose migrations.
2. **Unreachable error branches** (500s from database failure): A healthy in-memory SQLite database cannot simulate a connection failure or a general query error. These branches require a controlled failure injection.

## Decision

We will use two complementary approaches:

1. **Integration tests with real in-memory SQLite** for all happy paths, validation errors (400), not-found cases (404), conflict detection (409), and business logic branches. `TestMain` calls `data.Connect` which applies all goose migrations (schema + 25 seed players) to a shared in-memory database before the test suite runs.

2. **Mock injection via `MockPlayerService`** exclusively for error branches unreachable with a healthy database (e.g., simulating a 500 when `RetrieveAll` fails). `MockPlayerService` uses an opt-in function field pattern — each test overrides only the method relevant to the scenario; unset methods return safe zero-value defaults.

The boundary is explicit: if a scenario can be triggered with a real database, it must use the real database.

## Consequences

**Positive:**

- High confidence without test containers: real SQLite behavior — constraint enforcement, GORM error types, `ErrRecordNotFound` propagation — is exercised directly.
- Fast execution: in-memory SQLite requires no setup, teardown, or Docker.
- Error branches are explicitly covered without relying on fault injection at the infrastructure level.
- `MockPlayerService` is surgically scoped — tests are readable, overriding only what they test.

**Negative:**

- The test suite is coupled to SQLite behavior; switching to PostgreSQL in tests would require test containers or a separate PostgreSQL test instance.
- SQLite-specific constraint error strings are implicitly tested and would need updating for PostgreSQL.
- `TestMain` uses a shared in-memory database; tests that modify data must restore state or account for test execution order.
