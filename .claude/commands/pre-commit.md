Before running the checklist, run `git fetch origin`. If the current branch is behind `origin/master`, stop and rebase before proceeding.

Run the pre-commit checklist for this project:

1. Update `CHANGELOG.md` `[Unreleased]` section — read the current
   `CHANGELOG.md`, inspect `git diff` to understand what changed, then write
   the appropriate entry under the correct subsection (Added / Changed / Fixed /
   Removed), referencing the issue number. If the `[Unreleased]` section
   already contains an entry that covers these changes (e.g. added during
   release branch preparation via `/pre-release`), skip this step. Also skip
   this step for tooling/maintenance changes with no backing GitHub issue —
   feature and bug work should already have an issue per this project's SDD
   workflow, so this exception does not apply to them.
2. Run `go fmt ./...`
3. Note whether `go.mod`/`go.sum` already have uncommitted changes before
   running `go mod tidy` (`git diff --stat -- go.mod go.sum`). Run
   `go mod tidy`, then compare again — if tidy introduced additional changes
   beyond that baseline, stop and report it — propose committing the tidy
   result first and wait for my explicit confirmation before doing so.
4. Run `go vet ./...` — must pass.
5. Run `go build -v ./...` — must succeed.
6. Ask me: were any Swagger annotations modified? If yes, run `swag init`.
7. Run `go test -v ./... -coverpkg=github.com/nanotaboada/go-samples-gin-restful/service,github.com/nanotaboada/go-samples-gin-restful/controller,github.com/nanotaboada/go-samples-gin-restful/route -covermode=atomic -coverprofile=coverage.out` — all tests must pass, target 80%+ coverage for service, controller, route packages.

8. If Docker is running, run `docker compose build` — must succeed with no
   errors. Skip this step with a note if Docker Desktop is not running.
9. If `coderabbit` CLI is installed and `coderabbit review --help` reports
   `--agent` support, run `coderabbit review --type uncommitted --agent`:
   - If actionable/serious findings are reported, stop and address them before proposing the commit.
   - If only nitpick-level findings, report them and continue to the commit proposal.
   - If `coderabbit` is not installed, or the installed version does not
     support `--agent`, skip this step with a note.

Run step 1 (CHANGELOG update), then run steps 2–3 sequentially (stop at step 3 if go.mod or go.sum changed), then ask about step 6 and run `swag init` first if needed, then run steps 4, 5, and 7 in parallel, run step 8 (docker build), then run step 9 (CodeRabbit review) if available, report the results clearly, then propose a branch name and commit message for my approval using the format `type(scope): description (#issue)` (max 80 chars; types: `feat` `fix` `chore` `docs` `test` `refactor` `ci` `perf`) with the required `Co-Authored-By: Claude Sonnet 5 <noreply@anthropic.com>` line. Do not create the branch or commit until I explicitly confirm.
