Before running the checklist, run `git fetch origin`. If the current branch is behind `origin/master`, stop and rebase before proceeding.

Run the pre-commit checklist for this project:

1. Update `CHANGELOG.md` `[Unreleased]` section — read the current
   `CHANGELOG.md`, inspect `git diff` to understand what changed, then write
   the appropriate entry under the correct subsection (Added / Changed / Fixed /
   Removed), referencing the issue number. If the `[Unreleased]` section
   already contains an entry that covers these changes (e.g. added during
   release branch preparation via `/pre-release`), skip this step.
2. Run `go fmt ./...`
3. Run `go vet ./...` — must pass.
4. Run `go build -v ./...` — must succeed.
5. Ask me: were any Swagger annotations modified? If yes, run `swag init`.
6. Run `go test -v ./... -coverpkg=github.com/nanotaboada/go-samples-gin-restful/service,github.com/nanotaboada/go-samples-gin-restful/controller,github.com/nanotaboada/go-samples-gin-restful/route -covermode=atomic -coverprofile=coverage.out` — all tests must pass, target 80%+ coverage for service, controller, route packages.

7. If Docker is running, run `docker compose build` — must succeed with no
   errors. Skip this step with a note if Docker Desktop is not running.
8. If `coderabbit` CLI is installed, run `coderabbit review --type uncommitted --prompt-only`:
   - If actionable/serious findings are reported, stop and address them before proposing the commit.
   - If only nitpick-level findings, report them and continue to the commit proposal.
   - If `coderabbit` is not installed, skip this step with a note.

Run step 1 (CHANGELOG update), then run steps 2–4 and 6 in parallel (ask about step 5), run step 7 (docker build), then run step 8 (CodeRabbit review) if available, report the results clearly, then propose a branch name and commit message for my approval using the format `type(scope): description (#issue)` (max 80 chars; types: `feat` `fix` `chore` `docs` `test` `refactor` `ci` `perf`). Do not create the branch or commit until I explicitly confirm.
