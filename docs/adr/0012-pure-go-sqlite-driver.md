# ADR-0012: Pure-Go SQLite Driver

Date: 2026-04-20

## Status

Accepted

## Context

The project originally used `gorm.io/driver/sqlite`, which hard-codes a blank
import of `mattn/go-sqlite3` — a CGO library. This introduced three friction
points:

1. The Docker builder stage required `gcc`, `musl-dev`, and `sqlite-dev` to
   compile the CGO bindings, making it heavier than a pure-Go project.
2. Cross-compilation was constrained by the host toolchain because CGO cannot
   cross-compile without a matching sysroot.
3. Local builds required a C toolchain, raising the contributor setup bar.

ADR-0011 identified this constraint explicitly under "When to revisit" and
flagged `modernc.org/sqlite` as the natural resolution path.

Two options were evaluated:

- **Configure `gorm.io/driver/sqlite` via build tags** — not viable; the
  package hard-codes the blank import of `mattn/go-sqlite3` in its source, so
  there is no build-tag mechanism to swap the underlying driver.
- **Replace with `github.com/glebarez/sqlite`** — a GORM v2 wrapper purpose-
  built around `modernc.org/sqlite`, the pure-Go SQLite implementation produced
  by transpiling the upstream C source. It exposes the same `sqlite.Open()`
  function signature as `gorm.io/driver/sqlite`, making the swap a one-line
  import change. It supports both file-path and `file::memory:?cache=shared`
  DSNs, covering all existing test and production usage.

## Decision

We will replace `gorm.io/driver/sqlite` (backed by `mattn/go-sqlite3`) with
`github.com/glebarez/sqlite` (backed by `modernc.org/sqlite`).

The change is confined to a single import in `data/player_data.go`. The
`gorm.Open(sqlite.Open(dataSourceName), ...)` call, DSN formats, and the
`Connect()` function signature are all unchanged. `mattn/go-sqlite3` is removed
from the dependency graph entirely. The Docker builder stage drops `gcc`,
`musl-dev`, and `sqlite-dev` and sets `CGO_ENABLED=0`.

## Consequences

**Positive:**

- The Docker builder stage no longer requires a C toolchain; the image is
  simpler to maintain and faster to build.
- Setting `CGO_ENABLED=0` produces a fully statically linked binary, improving
  portability across Linux variants and enabling scratch-based images if desired
  in future.
- Cross-compilation becomes straightforward — no sysroot or cross-compiler
  required.
- Local development no longer requires `gcc` or SQLite development headers.

**Negative:**

- `modernc.org/sqlite` is a transpiled C codebase; its binary footprint is
  larger than the dynamically linked CGO equivalent.
- Performance characteristics differ slightly from the native CGO driver under
  write-heavy workloads; this is acceptable for this project's CRUD workload.
- The `modernc.org/sqlite` dependency graph is deeper than `mattn/go-sqlite3`,
  adding indirect dependencies to `go.mod`.
