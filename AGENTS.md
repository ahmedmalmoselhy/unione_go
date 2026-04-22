# Repository Guidelines

## Project Structure & Module Organization

This repository is a Go backend for UniOne. Entry code lives in `cmd/api/main.go`. Core application code is under `internal/`:

- `internal/handlers`: HTTP handlers
- `internal/services`: business logic
- `internal/repository`: GORM data access
- `internal/models`: persistence models
- `internal/routes`: Gin route wiring
- `internal/middlewares`: auth, logging, request context
- `internal/apiutil`: shared request/response helpers
- `internal/config`, `internal/db`: startup configuration and database bootstrap

Operational docs live at the repo root: `README.md`, `CURRENT_STATUS.md`, `API_ENDPOINTS.md`, `Enhancements.md`, and `IMPLEMENTATION_PLAN.md`.

## Build, Test, and Development Commands

- `go run ./cmd/api`: run the API locally
- `go test ./...`: run all tests and compile-check all packages
- `gofmt -w ./cmd ./internal`: format the codebase
- `go mod tidy`: clean up module dependencies

Use a valid `.env` with `DB_*`, `JWT_SECRET`, and SMTP settings before running locally.

## Coding Style & Naming Conventions

Use standard Go formatting; run `gofmt` before committing. Prefer short, focused packages and keep the existing layering: handlers should call services, services should call repositories.

- Exported names: `PascalCase`
- Unexported names: `camelCase`
- Route handlers: `CreateStudent`, `GetSection`, `UpdateCourse`
- Files: lowercase with underscores only when needed, e.g. `student_service.go`

Keep API responses consistent by using helpers from `internal/apiutil`.

## Testing Guidelines

The repository currently uses Go’s built-in `testing` tooling. Add tests next to the package they cover using `*_test.go` files, for example `internal/services/student_service_test.go`.

At minimum, new work should pass:

- `go test ./...`

When adding features, prefer service-level tests first, then handler/route coverage for externally visible behavior.

## Commit & Pull Request Guidelines

Follow the existing commit style from history:

- `phase2: add student CRUD APIs`
- `feat: expand implementation plan ...`

Use concise, imperative messages with a scope prefix when helpful. Keep one logical step per commit.

PRs should include:

- a short summary of the change
- linked phase/checklist item when applicable
- affected endpoints or modules
- confirmation that `go test ./...` passed
- doc updates when behavior or status changed

## Security & Configuration Tips

Do not commit real secrets in `.env`. Treat `JWT_SECRET`, database credentials, and SMTP credentials as local-only values. Validate config through `internal/config` rather than hardcoding fallbacks in handlers or services.
