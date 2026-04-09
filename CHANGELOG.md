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

- `tests/player_fake.go`: added `MakeUnknownPlayer()` factory — valid UUID absent from the database, for 404-by-lookup scenarios (#244)
- `tests/player_fake.go`: added package-level doc comment defining the three-term data-state vocabulary (`existing`, `nonexistent`, `unknown`) (#244)
- `migrations/00001_create_players_table.sql`: schema migration — creates `players` table with `id TEXT PRIMARY KEY` (UUID string) and unique index on `squadNumber` (#194)
- `migrations/00002_seed_starting11.sql`: seed migration — inserts 11 Starting XI players with deterministic UUID v5 values (#194)
- `migrations/00003_seed_substitutes.sql`: seed migration — inserts 14 Substitute players with deterministic UUID v5 values (#194)
- `migrations/embed.go`: embeds SQL migration files into the binary via `embed.FS` for path-independent deployment (#194)

### Changed

- `tests/player_fake.go`: renamed `MakeNonExistingPlayer` → `MakeNonexistentPlayer` to align with canonical cross-repo spelling (#244)
- `tests/main_test.go`: updated all `MakeNonExistingPlayer()` call sites to `MakeNonexistentPlayer()` (#244)
- `tests/main_test.go`: renamed `TestRequestGETPlayerByIDNonExistingResponseStatusNotFound` → `TestRequestGETPlayerByIDUnknownResponseStatusNotFound`; now uses `MakeUnknownPlayer()` (#244)
- `tests/main_test.go`: renamed `TestRequestPUTPlayerBySquadNumberNonExistingResponseStatusNotFound` → `TestRequestPUTPlayerBySquadNumberUnknownResponseStatusNotFound`; now uses `MakeUnknownPlayer()` (#244)
- `tests/main_test.go`: updated `NonExisting` table sub-cases to `Unknown` for GET/DELETE by squad number, driven by `MakeUnknownPlayer().SquadNumber` (#244)
- `data/player_data.go`: replaced `AutoMigrate` with `goose.SetBaseFS` + `goose.Up` for versioned, reproducible schema and seed management (#194)
- `data/player_data.go`: set `SetMaxOpenConns(1)` and `SetMaxIdleConns(1)` to prevent SQLite write contention under concurrent load (#194)
- `migrations/00001_create_players_table.sql`: changed `dateOfBirth` from `VARCHAR(20)` to `TEXT` (values are 24 chars; PostgreSQL enforces length) and `starting11` from `INTEGER` to `BOOLEAN` for cross-dialect correctness (#194)
- `migrations/00003_seed_substitutes.sql`: corrected team name typos — `Villareal` → `Villarreal`, `Nottingham Forrest` → `Nottingham Forest` (#194)
- `migrations/00002_seed_starting11.sql`, `migrations/00003_seed_substitutes.sql`: tightened Down migrations to delete by explicit UUID instead of broad `WHERE starting11 = 0/1` (#194)
- `tests/main_test.go`: removed `AutoMigrate` + JSON seeding from `TestMain`; the in-memory SQLite DB is now initialized via goose migrations, consistent with production (#194)
- `scripts/entrypoint.sh`: removed pre-seeded DB copy logic; migrations run automatically at app startup (#194)
- `Dockerfile`: removed `storage/hold/` copy; added `migrations/` to builder stage for embed compilation (#194)
- `README.md`: replaced pre-seeded database notes with goose migration commands and updated tech stack table (#194)
- `.gitignore`: added `storage/*.db` and SQLite sidecar patterns (`*.db-wal`, `*.db-shm`, `*.db-journal`) to exclude runtime-generated files (#194)

### Fixed

### Removed

- `storage/players-sqlite3.db`: removed pre-seeded binary database file; replaced by versioned SQL migrations (#194)
- `tests/player_fake.go`: removed `MakePlayersFromJSON()` — superseded by goose migrations as the seeding mechanism (#194)
- `tests/players.json`: removed JSON fixture file — superseded by SQL seed migrations as the authoritative source of player data (#194)

---

## [2.1.0 - Cafu] - 2026-04-04

### Added

- `.sonarcloud.properties`: SonarCloud Automatic Analysis configuration —
  sources, tests, coverage exclusions aligned with `codecov.yml` (#239)
- `.dockerignore`: added `.claude/`, `CLAUDE.md`, `.coderabbit.yaml`,
  `.sonarcloud.properties`, `.sonarlint/`, `CHANGELOG.md`, `README.md`
  (#239)
- `docs/adr/`: 10 Architecture Decision Records documenting key design choices — layered architecture, Gin, GORM, SQLite, UUID v4 primary key, squad number as mutation identifier, single domain struct with dedicated request binding type, full update semantics (PUT / PATCH deferred to #172), in-memory cache strategy, and mixed test strategy (closes #162)
- CD workflow (`go-cd.yml`): added "Verify tag commit is reachable from master" step to `deploy` job — fails early with a descriptive error if the tag's commit is not an ancestor of `origin/master` (closes #231)
- `.claude/commands/prerelease.md`: `/prerelease` slash command implementing a three-phase pre-release checklist — determine next version and player codename, prepare release branch and CHANGELOG, tag and push (closes #233)
- `.claude/commands/precommit.md`: step 1 marked as skippable when CHANGELOG has already been updated as part of release branch preparation (closes #233)

### Changed

- Player dataset normalised to November 2022 World Cup snapshot: Di María `abbrPosition` → `RW`, Mac Allister `team` → Brighton & Hove Albion, Fernández team/league → SL Benfica / Liga Portugal, Messi team/league → Paris Saint-Germain / Ligue 1 (closes #227)
- All player UUIDs migrated to canonical UUID v5 (namespace `f201b13e-c670-473d-885d-e2be219f74c8`, formula `{firstName}-{lastName}`) (closes #227)
- Test fixture for Create/Delete replaced: Paredes (squad 5) → Lo Celso (squad 27) (closes #227)
- `MakeUpdatePlayer()` added to `player_fake.go`: Emiliano Martínez (squad 23) (closes #227)
- DELETE test restructured: Armani (squad 1) → Lo Celso (squad 27) via POST+DELETE pattern (closes #227)
- GET by squad number body assertion retargeted to Messi (squad 10) (closes #227)
- `rest/players.rest` updated: `@newSquadNumber = 27`, `@existingSquadNumber = 23` (closes #227)

### Removed

- `sonar-project.properties`: replaced by `.sonarcloud.properties` for
  SonarCloud Automatic Analysis compatibility (#239)

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

[unreleased]: https://github.com/nanotaboada/go-samples-gin-restful/compare/v2.1.0-cafu...HEAD
[2.1.0 - Cafu]: https://github.com/nanotaboada/go-samples-gin-restful/compare/v2.0.0-bobby...v2.1.0-cafu
[2.0.0 - Bobby]: https://github.com/nanotaboada/go-samples-gin-restful/compare/v1.0.0-ademir...v2.0.0-bobby
[1.0.0 - Ademir]: https://github.com/nanotaboada/go-samples-gin-restful/releases/tag/v1.0.0-ademir