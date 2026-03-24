# AGENTS.md

## Cursor Cloud specific instructions

This is **labctl**, a Go CLI tool for [iximiuz Labs](https://labs.iximiuz.com). It is a single-binary Go project (no monorepo, no services to run).

### Build & Run

- `make build-dev` — builds `./labctl` binary in the repo root.
- Go 1.24.2 is required (matches `go.mod`); the environment already has it.

### Testing

- `go test -vet=off ./...` — runs all unit tests. The `-vet=off` flag is needed because the repo has pre-existing `go vet` warnings (non-constant format strings in `internal/labcli/cliutil.go` and `cmd/playground/tasks.go`). These are not test failures.
- `go vet ./...` — static analysis; currently reports 3 pre-existing warnings (see above).
- E2E tests exist at `e2e/exec` (`make test-e2e`) but require a live iximiuz Labs session/auth.

### Lint

No `.golangci.yml` or dedicated linter config exists in the repo. Use `go vet ./...` for static analysis.

### CLI demo (without auth)

After building, `./labctl version` and `./labctl --help` work without authentication. Commands that hit the iximiuz Labs API (e.g. `playground list`) require `labctl auth login` first.
