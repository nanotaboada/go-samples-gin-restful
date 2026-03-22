# ADR-0008: Full Update Semantics for PUT, PATCH Deferred

Date: 2026-03-21

## Status

Accepted (PUT — full update). Proposed (PATCH — tracked in #172).

## Context

HTTP defines two update semantics:

- **PUT**: Full replacement. The request body contains the complete new state of the resource. Fields absent from the body are zeroed.
- **PATCH**: Partial update. The request body contains only the fields to change. Unchanged fields are left as-is.

GORM provides two corresponding methods:

- `Save()`: Issues `UPDATE` covering all columns — PUT semantics.
- `Updates()`: Issues `UPDATE` for non-zero fields or an explicit map — PATCH semantics.

Implementing PATCH correctly in Go requires handling the zero-value ambiguity: a missing field and a field explicitly set to its zero value are indistinguishable in a decoded struct. Correct PATCH implementations typically use a map (`map[string]interface{}`) or a pointer-field struct, both of which sacrifice type safety or introduce a parallel type.

For a CRUD PoC scoped to basic operations, PATCH adds implementation complexity without introducing new architectural patterns relevant to the comparison study.

## Decision

We will implement PUT as a full replacement using GORM's `Save()`. The request body must contain all player fields; omitted fields are zeroed in the database. PATCH support is deferred and tracked in issue #172, where the zero-value handling strategy will be decided.

## Consequences

**Positive:**

- Simple, predictable behavior: the database state after a successful PUT exactly matches the request body — no partial-update edge cases.
- `Save()` semantics are straightforward to test: assert the full struct, not a subset.
- No zero-value ambiguity: all fields are always present.

**Negative:**

- Clients must fetch the current state before updating a single field — two round-trips for small changes.
- Higher payload size for updates that touch only one or two fields.
- Absence of PATCH is a REST completeness gap; acknowledged and tracked in #172.
