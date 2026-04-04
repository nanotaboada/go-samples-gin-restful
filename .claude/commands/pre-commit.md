Before running the checklist, run `git fetch origin`. If the current branch is behind `origin/master`, stop and rebase before proceeding.

Run the pre-commit checklist for this project:

1. *(Skippable)* Update `CHANGELOG.md` `[Unreleased]` section — add an entry
   under the appropriate subsection (Added / Changed / Fixed / Removed)
   describing the changes made, referencing the issue number. Skip this step
   if the CHANGELOG was already updated immediately before invoking
   `/pre-commit` (e.g. during release branch preparation via `/pre-release`).
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

Run steps 2–4 and 6 (ask about step 5), run step 7 (docker build), then run step 8 (CodeRabbit review) if available, report the results clearly, then propose a branch name and commit message for my approval using the format `type(scope): description (#issue)` (max 80 chars; types: `feat` `fix` `chore` `docs` `test` `refactor` `ci` `perf`). Do not create the branch or commit until I explicitly confirm.
