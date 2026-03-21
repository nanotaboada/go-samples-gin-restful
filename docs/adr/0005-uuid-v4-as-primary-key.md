# ADR-0005: Use UUID v4 as Primary Key

Date: 2026-03-21

## Status

Accepted

## Context

Primary key options for the `Player` entity:

- **Auto-increment integer**: DB-generated, simple, minimal storage. Leaks record count and insertion order to clients via the API — a privacy and security concern even in a PoC.
- **UUID v4**: Random 128-bit value, server-generated, opaque. Provides no ordering information.
- **UUID v7**: Time-ordered variant of UUID, server-generated, sortable by creation time — a superset of v4 benefits with added ordering.
- **ULID**: Time-ordered, URL-safe, lexicographically sortable — similar trade-offs to UUID v7, less standardized.

The `id` field is exposed in the API response and used as a URL parameter for `GET /players/:id`. Revealing insertion order or record count through sequential IDs is undesirable even in a non-sensitive domain.

## Decision

We will use UUID v4 strings as primary keys, generated server-side at creation time using `github.com/google/uuid`. The database does not generate IDs — GORM receives a fully populated struct on `Create`. Any `id` value supplied in a POST or PUT request body is overwritten by the controller before the service is called (see ADR-0007 for the structural enforcement of this).

UUID v7 was considered but not adopted — sortability by creation time provides no current benefit and would introduce an implicit ordering expectation into the API contract.

## Consequences

**Positive:**

- IDs are opaque: clients cannot infer record count or insertion order from the key.
- Server-side generation removes the need for a DB round-trip to retrieve the assigned ID after insert.
- IDs are stable and portable across environments — seeded data carries its UUIDs into every environment without conflict.

**Negative:**

- Larger storage footprint than an integer (36-byte string vs. 4–8 bytes for `int64`).
- Not sortable by insertion order; queries ordered by creation time require an explicit `created_at` column.
- Test fixtures use deterministic UUID v5 (derived from a fixed namespace + squad number) to avoid randomness in assertions — contributors must understand that test IDs are v5, not v4.
