# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Famous Football Player Names 🎖️

This project uses famous football player names (A-Z) as release codenames:

| Letter | Player Name | Country/Era | Tag Name |
| ------ | ----------- | ----------- | -------- |
| A | Ademir | Brazil | `ademir` |
| B | Bobby (Moore) | England | `bobby` |
| C | Cafu | Brazil | `cafu` |
| D | (Alfredo) Di Stéfano | Argentina/Spain | `distefano` |
| E | Eusébio | Portugal | `eusebio` |
| F | Franz (Beckenbauer) | Germany | `franz` |
| G | Garrincha | Brazil | `garrincha` |
| H | (Thierry) Henry | France | `henry` |
| I | (Filippo) Inzaghi | Italy | `inzaghi` |
| J | Jairzinho | Brazil | `jairzinho` |
| K | (Roy) Keane | Ireland | `keane` |
| L | (Frank) Lampard | England | `lampard` |
| M | (Diego) Maradona | Argentina | `maradona` |
| N | Nilton (Santos) | Brazil | `nilton` |
| O | (Jay-Jay) Okocha | Nigeria | `okocha` |
| P | Pelé | Brazil | `pele` |
| Q | (Fabio) Quagliarella | Italy | `quagliarella` |
| R | Romário | Brazil | `romario` |
| S | (Paul) Scholes | England | `scholes` |
| T | (Francesco) Totti | Italy | `totti` |
| U | Uwe (Seeler) | Germany | `uwe` |
| V | (Patrick) Vieira | France | `vieira` |
| W | (Ian) Wright | England | `wright` |
| X | Xabi (Alonso) | Spain | `xabi` |
| Y | (Lev) Yashin | USSR | `yashin` |
| Z | Zico | Brazil | `zico` |

---

## [Unreleased]

### Added

- `docs/adr/`: 10 Architecture Decision Records documenting key design choices — layered architecture, Gin, GORM, SQLite, UUID v4 primary key, squad number as mutation identifier, single domain struct with dedicated request binding type, full update semantics (PUT / PATCH deferred to #172), in-memory cache strategy, and mixed test strategy (closes #162)

### Changed

- Player dataset normalised to November 2022 World Cup snapshot: Di María `abbrPosition` → `RW`, Mac Allister `team` → Brighton & Hove Albion, Fernández team/league → SL Benfica / Liga Portugal, Messi team/league → Paris Saint-Germain / Ligue 1 (closes #227)
- All player UUIDs migrated to canonical UUID v5 (namespace `f201b13e-c670-473d-885d-e2be219f74c8`, formula `{firstName}-{lastName}`) (closes #227)
- Test fixture for Create/Delete replaced: Paredes (squad 5) → Lo Celso (squad 27) (closes #227)
- `MakeUpdatePlayer()` added to `player_fake.go`: Emiliano Martínez (squad 23) (closes #227)
- DELETE test restructured: Armani (squad 1) → Lo Celso (squad 27) via POST+DELETE pattern (closes #227)
- GET by squad number body assertion retargeted to Messi (squad 10) (closes #227)
- `rest/players.rest` updated: `@newSquadNumber = 27`, `@existingSquadNumber = 23` (closes #227)

---

## [2.0.0 - Bobby] - 2026-03-19

### Added

- `GET /players/{id}` now accepts a UUID string (surrogate key) instead of an integer ID
- `tools/seed_001_starting_eleven.go`, `tools/seed_002_substitutes.go`: `//go:build ignore` standalone Go programs to drop, recreate, and reseed the SQLite database using GORM; numbered to mirror the future Goose migration sequence
- `rest/players.rest`: HTTP request file covering health check, POST, GET all, GET by ID, GET by squad number, PUT, and DELETE — compatible with VS Code REST Client (`humao.rest-client`) and JetBrains built-in HTTP Client
- `humao.rest-client` listed in `.vscode/extensions.json` recommended extensions

### Changed

- **BREAKING** `Player.ID` field changed from `int` to `string` to hold a UUID v4; the server always generates the ID on POST — any client-provided value is overwritten
- **BREAKING** `PUT /players/:squadnumber` and `DELETE /players/:squadnumber` now identify players by `squadNumber` (user-facing unique identifier) instead of internal ID; clients must update URLs from `/players/:id` to `/players/:squadnumber`
- **BREAKING** `GET /players/{id}` parameter type changed from integer to UUID string
- `uniqueIndex` GORM tag added to `Player.SquadNumber` — uniqueness is now enforced at DB level, not only in application logic
- Test fixtures (`tests/players.json`) migrated from integer IDs to deterministic UUID v5 strings derived from `squadNumber` using a project-specific namespace
- SQLite database (`storage/players-sqlite3.db`) re-seeded with UUID v4 primary keys to match the new schema
- Updated `codecov.yml` ignore list: replaced `postman-collections/**/*` with `rest/**/*`
- Updated `README.md`: replaced Postman Collection section with HTTP Requests section referencing `rest/players.rest`

### Removed

- Integer auto-increment `id` — the `id` field is now a server-assigned opaque UUID v4 string
- `postman-collections/` directory and Postman collection JSON file

### Fixed

- Eliminated variable shadowing in `Put` and `Delete` handlers: inner `err :=` assignments replaced with `err =` to reuse the outer declaration

### Migration Notes

This release changes the database schema in a **backwards-incompatible** way.
An existing `players-sqlite3.db` with integer primary keys will fail at startup because `AutoMigrate` cannot alter a column type in SQLite.
Re-create the database from scratch using the provided seed scripts:

```bash
# Drop and reseed with starting-eleven players (squad numbers 1–11)
go run tools/seed_001_starting_eleven.go

# Append substitute players (squad numbers 12–23)
go run tools/seed_002_substitutes.go
```

Both scripts use `//go:build ignore` so they are excluded from normal `go build ./...` and `go test ./...` runs.
Run them only when you need to recreate the local SQLite database.

---

## [1.0.0 - Ademir] - 2026-02-06

Initial release. See [README.md](README.md) for complete feature list and documentation.

---

## How to Release

The full release procedure — branch, PR, tag, and CD workflow — is documented in
[README.md § Create a Release](README.md#create-a-release).

In summary: move items from `[Unreleased]` to a new `[X.Y.Z - Player]` section
(see template below), open a `release/vX.Y.Z-player` PR, merge it into `master`,
then push the annotated tag to trigger the CD workflow.

---

<!-- Template for new releases:

## [X.Y.Z - PLAYER_NAME] - YYYY-MM-DD

### Added
- New features

### Changed
- Changes in existing functionality

### Deprecated
- Soon-to-be removed features

### Removed
- Removed features

### Fixed
- Bug fixes

### Security
- Security vulnerability fixes

-->

---

[unreleased]: https://github.com/nanotaboada/go-samples-gin-restful/compare/v2.0.0-bobby...HEAD
[2.0.0 - Bobby]: https://github.com/nanotaboada/go-samples-gin-restful/compare/v1.0.0-ademir...v2.0.0-bobby
[1.0.0 - Ademir]: https://github.com/nanotaboada/go-samples-gin-restful/releases/tag/v1.0.0-ademir
