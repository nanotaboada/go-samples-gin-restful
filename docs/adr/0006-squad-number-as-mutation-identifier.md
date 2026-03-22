# ADR-0006: Use Squad Number as Resource Identifier for Mutations

Date: 2026-03-21

## Status

Accepted

## Context

The `Player` entity has two unique identifiers:

- `id`: UUID v4, internal, server-generated, opaque to clients.
- `squadNumber`: integer, user-facing, domain-meaningful, unique within the squad.

REST mutations (PUT, DELETE) require a stable URL identifier. The question is which identifier to expose:

- **UUID as the sole identifier**: Consistent — one identifier for all operations. Requires clients to know the UUID before mutating, typically obtained via a prior GET. UUIDs are impractical for manual API calls or documentation examples.
- **Squad number as the mutation identifier**: Human-readable and domain-meaningful. Squad numbers are publicly known (squad lists are published); no prior lookup required.
- **Both identifiers for all operations**: Maximum flexibility, doubles the route surface, introduces ambiguity about which is canonical.

## Decision

We will use squad number as the URL parameter for all mutation endpoints:

- `PUT /players/squadnumber/:squadnumber`
- `DELETE /players/squadnumber/:squadnumber`

GET supports both identifiers:

- `GET /players/:id` — by internal UUID
- `GET /players/squadnumber/:squadnumber` — by squad number

The routing asymmetry is intentional: the UUID is available for programmatic lookup but clients are not expected to use it for day-to-day mutations. Gin's static-segment-priority trie routing resolves the `/players/squadnumber/` vs `/players/:id` conflict without ambiguity.

## Consequences

**Positive:**

- Mutation URLs are human-readable and meaningful in the domain (`/players/squadnumber/10` vs `/players/6ba7b810-...`).
- No prior GET required to update or delete a known player.
- Squad numbers are stable for the duration of a tournament.

**Negative:**

- Routing is asymmetric: GET exposes two identifiers; mutations expose only one. New contributors may find this surprising.
- PUT must perform two DB calls: a `RetrieveBySquadNumber` to fetch the existing UUID, then `Save` — the UUID must be preserved on the struct or GORM will attempt to zero the primary key.
- A squad number change would require special handling: the URL key and the body field are validated to match (`player.SquadNumber != squadNumber` → 400), preventing silent key drift.
