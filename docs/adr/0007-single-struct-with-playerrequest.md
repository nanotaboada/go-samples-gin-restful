# ADR-0007: Use Single Domain Struct with Dedicated Request Binding Type

Date: 2026-03-21

## Status

Accepted

## Context

The `Player` struct currently serves three roles simultaneously: GORM entity, HTTP request body, and HTTP response body. This is viable while all three roles share the same field set, but two specific problems arise:

1. **Mass assignment**: The `id` field is present in the struct bound from POST/PUT request bodies. The controller overwrites it with `uuid.NewString()`, but the protection is enforced by discipline, not structure. A future contributor adding a server-controlled field (e.g., `createdAt`, an audit flag) could inadvertently expose it to client writes.

2. **Swagger schema accuracy**: Swaggo generates a single schema from `model.Player` used for both POST request and GET response documentation. The `id` field appears as client-writable in the POST schema, which is incorrect — it is server-generated and any client-supplied value is silently overwritten.

Three approaches were evaluated:

- **Keep single struct**: Simple, zero mapping code. Does not address the Swagger inaccuracy or the mass assignment surface.
- **Full separation (request / domain / response)**: Three distinct types with explicit mapping functions between them. Correct in theory, but `PlayerResponse` today would be an identity mapping of `model.Player` — all fields identical, no justification.
- **Dedicated request type only**: Introduce `PlayerRequest` for binding, keep `model.Player` for GORM and responses. Closes the specific gaps without premature separation.

## Decision

We will introduce `PlayerRequest` as a dedicated type used exclusively for request binding on POST and PUT handlers. `PlayerRequest` omits the `id` field, making it structurally impossible to bind a client-supplied ID. The controller constructs a `model.Player` from the `PlayerRequest` before calling the service — a direct field assignment, not a mapping helper.

`model.Player` remains the single type for GORM operations and HTTP responses. No `PlayerResponse` type is introduced until response fields genuinely diverge from the domain struct.

## Consequences

**Positive:**

- Client-supplied `id` is structurally rejected at binding time, not silently overwritten in the controller.
- Swagger POST schema accurately excludes `id` — the generated documentation reflects the actual API contract.
- No mapping code for responses; `model.Player` is returned directly from GET handlers.

**Negative:**

- Asymmetry between request and response types may surprise contributors expecting either full separation or full unification.
- When response fields eventually diverge (computed fields, renamed fields, omitted internals), `PlayerResponse` must be introduced at that point — this ADR should be superseded at that time.
