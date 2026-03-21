# ADR-0001: Implement Layered Architecture

Date: 2026-03-21

## Status

Accepted

## Context

Go does not prescribe a project layout. Common alternatives include:

- **Flat package structure**: Everything in root or a single `internal/` package — minimal structure, suitable for very small programs.
- **Feature-based packages**: Organized by domain entity rather than technical concern.
- **golang-standards/project-layout**: Uses `cmd/`, `internal/`, `pkg/` — designed for multi-binary projects or libraries with public APIs.
- **Hexagonal / clean architecture**: Ports and adapters pattern — maximizes decoupling at the cost of significant boilerplate.

Gin is a minimal framework that imposes no structure. A single-domain, learning-focused PoC with one binary and no public library surface does not benefit from the golang-standards layout, which targets a different project shape.

## Decision

We will use a four-layer architecture with explicit package boundaries:

```
route → controller → service → data
```

Each layer depends only on the abstraction of the layer below it. Controllers depend on the `PlayerService` interface (not the concrete `playerService` struct), enabling mock injection in tests without modifying production code. No layer may skip another — routes call controllers, controllers call services, services call data.

## Consequences

**Positive:**

- Familiar to contributors from other languages — the pattern closely matches MVC-adjacent conventions in Spring Boot, ASP.NET Core, and Express.
- Controllers are independently testable via interface injection without a database.
- Each concern has a single, predictable home: HTTP parsing in controllers, business logic in services, GORM queries in services/data.

**Negative:**

- More packages and files than a flat structure for a single-domain project.
- The golang-standards layout (`cmd/`, `internal/`) may be expected by experienced Go contributors; the deviation is intentional but should be explained to new contributors.

**When to revisit:**

- If a second domain entity is added and the flat layer packages become difficult to navigate.
- If a second binary (e.g., a migration CLI) is introduced, warranting `cmd/` separation.
