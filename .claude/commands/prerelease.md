Run the pre-release checklist for this project. Work through all three phases
in order, pausing for explicit confirmation at each decision point before
proceeding. Never create a branch, commit, tag, or push without approval.

---

## Phase 1 — Determine next release

1. Run `git status` and confirm the working tree is clean and on `master`.
   If not, stop and report the problem.

2. Run `git tag --sort=-v:refname` to list existing tags. Identify the most
   recent tag matching `v*.*.*-*` and extract its player codename.

3. Read the A–Z player table from `CHANGELOG.md` to find the player that
   follows the last used codename alphabetically. That is the next player.

4. Read the `[Unreleased]` section of `CHANGELOG.md` and infer the version
   bump using these rules (applied in order — first match wins):
   - Any entry contains the word **BREAKING** → **major** bump
   - Any `### Added` subsection has entries → **minor** bump
   - Otherwise (only `### Changed`, `### Fixed`, `### Removed`) → **patch** bump

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

4. Run `/precommit`, explicitly skipping Step 1 (CHANGELOG update — already
   completed in step 2 above). Tell `/precommit` to skip Step 1 by opening
   with: "Skip Step 1 — CHANGELOG was already updated as part of this release
   branch." `/precommit` Step 1 is optional when the CHANGELOG has been
   updated by the release preparation step immediately prior.

5. Propose opening a PR from `release/vX.Y.Z-{player}` into `master`.
   **Wait for explicit approval before opening.**

6. Open the PR with:
   - Title: `docs(changelog): prepare release notes for vX.Y.Z-{player}`
   - Body summarising what is included in this release.

---

## Phase 3 — Tag and release

1. Wait — do not proceed until the user confirms:
   - CI is green
   - CodeRabbit has reviewed (if applicable)
   - The PR has been merged into `master`

2. Once confirmed, run:
   ```bash
   git checkout master && git pull origin master
   ```
   and show the resulting `git log --oneline -3`.

3. Propose the annotated tag:
   ```bash
   git tag -a vX.Y.Z-{player} -m "Release vX.Y.Z - PlayerName"
   ```

   **Wait for explicit approval before creating the tag.**

4. Create the tag, then propose:
   ```bash
   git push origin vX.Y.Z-{player}
   ```

   **Wait for explicit approval before pushing.** Remind the user that pushing
   the tag triggers the CD workflow which will build, publish the Docker image,
   and create the GitHub Release.
