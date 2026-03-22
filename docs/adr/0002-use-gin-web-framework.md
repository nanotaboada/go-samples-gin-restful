# ADR-0002: Use Gin Web Framework

Date: 2026-03-21

## Status

Accepted

## Context

Go's standard library `net/http` is production-capable and requires no external dependencies. Alternatives evaluated:

- **net/http + ServeMux**: Zero dependencies. Requires manual JSON binding, no built-in middleware chaining, verbose path parameter handling.
- **Echo**: API similar to Gin, comparable performance. Smaller middleware ecosystem; `echo/middleware` does not include a page-cache equivalent.
- **Fiber**: Express-inspired, high performance. Does **not** use `net/http` — incompatible with standard `http.Handler` middleware and the broader Go HTTP ecosystem.
- **Chi**: Minimal, `net/http`-compatible, no built-in JSON binding or response helpers.

Two specific integration requirements shaped this decision: response caching via `gin-contrib/cache` and Swagger UI via `swaggo/gin-swagger`. Both are first-party integrations designed around Gin's handler type.

## Decision

We will use Gin as the HTTP framework. `gin-contrib/cache` provides the page-level caching middleware required by the caching strategy (see ADR-0009). `swaggo/gin-swagger` provides the Swagger UI integration. Gin's handler type (`func(*gin.Context)`) is the standard unit of composition throughout the codebase — for routes, middleware, and handlers alike.

## Consequences

**Positive:**

- `gin-contrib/cache` and `swaggo/gin-swagger` are available as first-party integrations built around Gin's types.
- Route groups, path parameters, JSON binding (`ShouldBindJSON`), and response helpers (`IndentedJSON`, `Status`) are built in.
- Widely used in Go tutorials and examples — high familiarity for contributors new to Go.

**Negative:**

- Gin wraps `net/http` but handlers receive `*gin.Context`, not `http.ResponseWriter`/`*http.Request` — standard `http.Handler` middleware requires an adapter.
- All handlers are coupled to `gin.Context`; switching frameworks requires rewriting every handler.
- Fiber's raw performance benchmarks exceed Gin's, but Fiber's incompatibility with `net/http` was disqualifying.
