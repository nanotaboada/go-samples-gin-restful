# ADR-0011: Docker and Compose Strategy

Date: 2026-04-02

## Status

Accepted

## Context

The project needs to run in a self-contained environment for demos, CI, and
as a reference point in the cross-language comparison. Two concerns apply:

1. **Image size and security**: a single-stage build would include the Go
   toolchain, compiler cache, and development headers in the final image.
2. **Local orchestration**: contributors should be able to start the
   application with a single command, without installing Go, configuring
   environment variables, or managing a database file manually.

An additional constraint is CGO: GORM's SQLite driver (`mattn/go-sqlite3`)
is a CGO library. Pure-Go SQLite drivers exist (e.g. `modernc.org/sqlite`)
but were not evaluated as an alternative for this project — the CGO
requirement shapes several build decisions below.

Available approaches:

- **Single-stage build (golang base image)**: Simplest Dockerfile, but the
  final image carries the full Go toolchain (~300 MB for alpine) and
  exposes build tooling unnecessarily.
- **Multi-stage build**: Builder stage compiles the binary; runtime stage
  is a minimal base with only what the application needs at runtime.
- **Scratch or distroless runtime**: Smallest possible image, but CGO
  introduces dynamic linking against musl libc; a fully static binary
  requires extra linker flags and produces a larger binary than
  dynamically linked equivalents. Both options also complicate the
  health check, which relies on `curl`.

## Decision

We will use a multi-stage Docker build with an Alpine-based runtime, and
Docker Compose to orchestrate the application locally.

- **Builder stage** (`golang:1.25-alpine3.23`): installs `gcc`,
  `musl-dev`, and `sqlite-dev` to satisfy CGO dependencies; builds the
  binary with `-trimpath -ldflags="-s -w"` to strip debug symbols and
  reduce size; uses `--mount=type=cache` for the Go module and build caches
  to speed up successive builds.
- **Runtime stage** (`alpine:3.23`): installs only `curl` (for the health
  check); copies the compiled binary, pre-seeded database, entrypoint
  script, and healthcheck script; creates a non-root `gin` user following
  the principle of least privilege.
- **Entrypoint script**: on first start, copies the pre-seeded database
  from the image's read-only `hold/` directory to the writable named
  volume at `/storage/`. On subsequent starts, the existing volume file
  is used unchanged.
- **Compose (`compose.yaml`)**: defines a single service with port mapping
  (`9000`), a named volume (`storage`), and environment variables
  (`STORAGE_PATH`, `GIN_MODE=release`). Health checks are declared in the
  Dockerfile (`GET /health`); Compose relies on that declaration.

## Consequences

**Positive:**

- The runtime image excludes the Go toolchain, C compiler, and SQLite
  headers — smaller attack surface and faster pulls.
- `docker compose up` is a complete local setup with no prerequisites
  beyond Docker.
- The named volume preserves data across restarts; `docker compose down -v`
  resets it cleanly.
- Bundling the seed database means the image is self-contained and
  demo-ready without a separate seeding step.

**Negative:**

- CGO requires `gcc` and `musl-dev` in the builder stage, making the
  builder heavier than a pure-Go project and coupling the build to a
  libc implementation (musl).
- Multi-stage builds are more complex to read and maintain than
  single-stage builds.
- The SQLite database file is versioned and bundled, meaning schema changes
  require a Docker image rebuild.
- The copy-on-first-run pattern in the entrypoint adds startup logic that
  must be understood when debugging volume state.

**When to revisit:**

- If the SQLite driver is replaced with a pure-Go alternative
  (`modernc.org/sqlite`), CGO can be disabled and the builder stage
  simplified significantly.
- If a second service (e.g. PostgreSQL) is added, Compose will need a
  dedicated network and dependency ordering.
