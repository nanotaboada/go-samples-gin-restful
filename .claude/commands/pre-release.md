Run the pre-release checklist for this project. Work through all three phases
in order, pausing for explicit confirmation at each decision point before
proceeding. Never create a branch, commit, tag, or push without approval.

---

## Phase 1 — Determine next release

1. Run `git status` and confirm the working tree is clean and on `master`.
   If not, stop and report the problem.

2. Run `git tag --sort=-v:refname` to list existing tags. Identify the most
   recent tag matching `v*.*.*-*` and extract its player codename.

3. Read the A–Z player table from `CHANGELOG.md` to find the next player:
   - **No tags yet**: start at `A` (first player in the table).
   - **Normal case**: use the player that follows the last used codename
     alphabetically. If letters were skipped, pick the next after the
     highest existing codename — do not backfill gaps.
   - **Last codename is `Z`** (Zico): the list is finite. Stop and flag
     that the naming convention needs to be revisited before proceeding.

4. Read the `[Unreleased]` section of `CHANGELOG.md` and infer the version
   bump using these rules (applied in order — first match wins):
   - Any entry contains the word **BREAKING** (case-insensitive), a
     `BREAKING CHANGE:` token in a commit footer, or a `!` suffix after
     the commit type/scope (e.g. `feat!:` or `feat(scope)!:`) → **major** bump
   - Any `### Added` entry introduces a **new endpoint, request parameter, or
     client-visible API response** → **minor** bump. Test additions, documentation
     (ADRs, README, Swagger regen), and internal tooling do **not** qualify.
   - Otherwise (only `### Changed`, `### Fixed`, `### Removed`, or non-API
     additions) → **patch** bump

5. Compute the next version by applying the bump to the current latest tag's
   semver (e.g. `v2.0.0-bobby` + minor → `2.1.0`).

6. Present a summary for confirmation before continuing:
   - Last tag and player
   - Next version and player codename
   - Bump type and the reasoning (what triggered it)
   - Proposed tag: `vX.Y.Z-{player}`
   - Proposed branch: `release/vX.Y.Z-{player}`

   **Wait for explicit approval before proceeding to Phase 2.**

---

## Phase 2 — Prepare release branch

1. Create branch `release/vX.Y.Z-{player}` from `master`.

2. Edit `CHANGELOG.md`:
   - Replace `## [Unreleased]` with `## [X.Y.Z - PlayerName] - YYYY-MM-DD`
     (use today's date; use the player's display name from the table, e.g.
     "Cafu", "Di Stéfano").
   - Consolidate duplicate subsection headings (e.g. two `### Added` blocks
     should be merged into one).
   - Add a new empty `## [Unreleased]` section at the top (above the new
     versioned heading) with the standard subsections.
   - Update the compare links at the bottom of the file:
     - `[unreleased]` → `.../compare/vX.Y.Z-{player}...HEAD`
     - Add `[X.Y.Z - PlayerName]` → `.../compare/v{prev-tag}...vX.Y.Z-{player}`

3. Show the full diff of `CHANGELOG.md` and propose this commit message:

   ```bash
   docs(changelog): prepare release notes for vX.Y.Z-{player} (#issue)
   ```

   **Wait for explicit approval before committing.**

4. Run the following checks. Run steps a–c and e in parallel; step d depends
   on the answer to the Swagger question so ask first then run if needed:

   a. `go fmt ./...`
   b. `go vet ./...` — must pass.
   c. `go build -v ./...` — must succeed.
   d. Ask: were any Swagger annotations modified? If yes, run `swag init`.
   e. `go test -v ./... -coverpkg=github.com/nanotaboada/go-samples-gin-restful/service,github.com/nanotaboada/go-samples-gin-restful/controller,github.com/nanotaboada/go-samples-gin-restful/route -covermode=atomic -coverprofile=coverage.out`
      — all tests must pass; target 80%+ coverage for service, controller, route.
   f. Run `golangci-lint run` — must pass with 0 issues before continuing.
   g. If Docker is running, run `docker compose build` — must succeed. Skip with
      a note if Docker Desktop is not running.
   h. If `coderabbit` CLI is installed, run
      `coderabbit review --type uncommitted --prompt-only`:
      - Actionable/serious findings → stop and address before continuing.
      - Nitpick-level findings → report and continue.
      - Not installed → skip with a note.

5. Propose opening a PR from `release/vX.Y.Z-{player}` into `master`.
   **Wait for explicit approval before opening.**

6. Open the PR with:
   - Title: `docs(changelog): prepare release notes for vX.Y.Z-{player}`
   - Body summarising what is included in this release.

---

## Phase 3 — Tag and release

1. Wait — do not proceed until the user confirms:
   - CI is green
   - The PR has been merged into `master`

2. Once confirmed, run:
   ```bash
   git checkout master && git pull origin master
   ```
   and show the resulting `git log --oneline -3`.

3. Propose the annotated tag:
   ```bash
   git tag -a vX.Y.Z-{player} -m "Release X.Y.Z - PlayerName"
   ```

   **Wait for explicit approval before creating the tag.**

4. Create the tag, then propose:
   ```bash
   git push origin vX.Y.Z-{player}
   ```

   **Wait for explicit approval before pushing.** Remind the user that pushing
   the tag triggers the CD workflow which will build, publish the Docker image,
   and create the GitHub Release.
