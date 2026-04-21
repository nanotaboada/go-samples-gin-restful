# Architecture Decision Records

This directory contains Architecture Decision Records (ADRs) for this project.

An ADR documents a significant architectural decision: the context that forced it,
the decision made, and the consequences — positive and negative. ADRs are
append-only. When a decision changes, a new ADR is added with status
`SUPERSEDED by ADR-NNNN`; the original is never deleted or rewritten.

## When to write an ADR

A decision warrants an ADR when all three conditions are true:

1. **A real fork existed** — a genuine alternative was considered and rejected.
2. **The code doesn't explain the why** — a contributor reading the source cannot infer the reasoning.
3. **Revisiting it would be costly** — changing it later requires significant rework.

If a decision is well-explained by inline comments or belongs in `CONTRIBUTING.md`
(process conventions, style guidelines), it does not need an ADR.

## Template

See [template.md](template.md) for the standard format (Michael Nygard).

## Index

| # | Title | Status | Date |
|---|---|---|---|
| [0001](0001-layered-architecture.md) | Layered Architecture | Accepted | 2026-03-21 |
| [0002](0002-use-gin-web-framework.md) | Use Gin Web Framework | Accepted | 2026-03-21 |
| [0003](0003-adopt-gorm-orm.md) | Adopt GORM as ORM | Accepted | 2026-03-21 |
| [0004](0004-use-sqlite-for-development.md) | Use SQLite for Development | Accepted | 2026-03-21 |
| [0005](0005-uuid-v4-as-primary-key.md) | Use UUID v4 as Primary Key | Accepted | 2026-03-21 |
| [0006](0006-squad-number-as-mutation-identifier.md) | Use Squad Number as Mutation Identifier | Accepted | 2026-03-21 |
| [0007](0007-single-struct-with-playerrequest.md) | Single Domain Struct with Dedicated Request Binding Type | Accepted | 2026-03-21 |
| [0008](0008-full-update-no-patch.md) | Full Update Semantics for PUT, PATCH Deferred | Accepted | 2026-03-21 |
| [0009](0009-in-memory-cache-strategy.md) | In-Memory Page Cache for GET Endpoints | Accepted | 2026-03-21 |
| [0010](0010-mixed-test-strategy.md) | Mixed Test Strategy | Accepted | 2026-03-21 |
| [0011](0011-docker-and-compose-strategy.md) | Docker and Compose Strategy | Accepted | 2026-04-02 |
| [0012](0012-pure-go-sqlite-driver.md) | Pure-Go SQLite Driver | Accepted | 2026-04-20 |
