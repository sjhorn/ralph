# Ralph

AI-powered feature planner and builder CLI.

## Build & Run

```bash
make build            # or: go build -o ralph ./cmd/ralph/
./ralph               # prints usage
```

## Lint

```bash
go vet ./...
```

## Project Structure

- `cmd/ralph/main.go` — entire CLI (single file, intentional)
- `example/` — separate Flutter test project (not part of the Go build)
- `.ralph/` — runtime directory created by `ralph init` in target projects

The single-file architecture is intentional — ralph is a simple CLI with no internal packages and stdlib-only dependencies.

## Origin & Design Reference

Ralph is a Go reimplementation of Adam Tuttle's RALPH workflow. See `docs/adam-tuttle-ralph-workflow.md` for the original blog post.

Key design principles from the original:
- **`plan` reads existing PRD**: The `plan` command must study the existing `.ralph/prd.json` (if it exists) before planning new work. This allows iterative planning — adding to or refining an existing PRD across multiple sessions.
- **`build` updates PRD inline**: Each build iteration marks items as `passes: true` in `prd.json` and appends to `progress.md`.
- **Git commit per iteration**: Each completed item gets its own git commit.
- **Progress tracking**: `progress.md` is an append-only log that serves as institutional memory for the next iteration.
