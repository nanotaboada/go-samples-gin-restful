Run the pre-commit checklist for this project:

1. Remind me to update `CHANGELOG.md` `[Unreleased]` section (Added / Changed / Fixed / Removed) — I must do this manually.
2. Run `go fmt ./...`
3. Run `go vet ./...` — must pass.
4. Run `go build -v ./...` — must succeed.
5. Ask me: were any Swagger annotations modified? If yes, run `swag init`.
6. Run `go test -v ./... -coverpkg=github.com/nanotaboada/go-samples-gin-restful/service,github.com/nanotaboada/go-samples-gin-restful/controller,github.com/nanotaboada/go-samples-gin-restful/route -covermode=atomic -coverprofile=coverage.out` — all tests must pass, target 80%+ coverage for service, controller, route packages.

Run steps 2–4 and 6 (ask about step 5), report the results clearly, then propose a branch name and commit message for my approval using the format `type(scope): description (#issue)` (max 80 chars; types: `feat` `fix` `chore` `docs` `test` `refactor` `ci` `perf`). Do not create the branch or commit until I explicitly confirm.
