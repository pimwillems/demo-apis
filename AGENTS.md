# AGENTS.md

Guidance for LLM/agentic tools working in this repository.

## What this repository is

A collection of small, independent demo REST APIs written in Go, used to teach
students. This is **not** a monorepo with shared code — it is a flat container of
unrelated Go modules, each with its own `go.mod`, `README.md`, and `AGENTS.md`.

There is currently no root-level Go module, build system, or CI. Work happens inside
one project directory at a time.

## Projects

- **`go-bookshop/`** — books, orders, and a customer loyalty programme. MVC-style
  layering (`internal/models`, `internal/store`, `internal/view`,
  `internal/handlers`). See `go-bookshop/AGENTS.md` for its specifics.
- **`go-library/`** — minimal books CRUD API in a flat two-file layout
  (`main.go`, `books.go`), with a `Dockerfile` for deployment. See
  `go-library/AGENTS.md` for its specifics.

**Always read the target project's own `AGENTS.md` before making changes there** — it
has the authoritative, up-to-date constraints and conventions for that module.

## Shared conventions across all projects

These demos are intentionally kept simple and consistent so students can compare them:

- **Go standard library only.** No third-party routers, frameworks, or ORMs (no gin,
  echo, chi, sqlx, gorm, etc.) unless a project's own docs or the user explicitly ask
  for one.
- **Go 1.22+ `ServeMux` routing** — method + wildcard patterns, e.g.
  `mux.HandleFunc("GET /books/{id}", ...)`, params via `r.PathValue("id")`.
- **In-memory storage only.** Data resets on restart. Do not introduce a database or
  persistence layer unless explicitly requested — it would defeat the purpose of a
  quick-start teaching demo.
- **All responses are JSON**; errors use the shape `{"error": "..."}` with an
  appropriate 4xx/5xx status.
- **`GET /health`** returns `{"status":"ok"}` in every project — keep it if you touch
  routing.

## Adding a new demo project

If asked to add a new demo API:

1. Create a new top-level directory (e.g. `go-<name>/`) with its own `go.mod` —
   don't fold it into an existing module.
2. Follow the shared conventions above unless told otherwise.
3. Give it its own `README.md` (endpoints, run instructions, examples) and
   `AGENTS.md` (constraints, architecture, how to verify), matching the style of the
   existing projects.
4. Add a row for it to the root `README.md` project table.

## Verifying changes

Each project builds and runs independently:

```bash
cd go-bookshop   # or go-library
go build ./... && go vet ./...
go run .
```

Check the project's own `AGENTS.md` for its exact verification steps (test ports,
curl examples, etc.) before considering a change done.
