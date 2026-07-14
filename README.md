# demo-apis

A collection of small, self-contained demo REST APIs used for teaching students. Each
subproject is an independent Go module — no shared code, no monorepo tooling. Pick a
directory, `cd` in, and follow its own README.

## Projects

| Project | Description |
| --- | --- |
| [`go-bookshop/`](go-bookshop/) | Book webshop API with books, orders, and a customer loyalty programme. MVC-style layout (`internal/models`, `internal/store`, `internal/view`, `internal/handlers`). |
| [`go-library/`](go-library/) | Minimal books CRUD API (list/get/create). Flat two-file layout, includes a `Dockerfile` set up for deployment (e.g. Coolify). |

## Common traits

Both APIs share the same teaching-oriented design:

- **Go standard library only** — `net/http` and `encoding/json`, no web frameworks or
  routers, no ORMs, no external dependencies.
- **Go 1.22+ routing** — uses the stdlib `ServeMux` method/wildcard patterns
  (`mux.HandleFunc("GET /books/{id}", ...)` + `r.PathValue("id")`).
- **In-memory data** — everything resets when the server restarts. There is no
  database; that's intentional, to keep the demos easy to read and run.
- **JSON everywhere** — every response is JSON, and errors follow `{"error": "..."}`.
- **`GET /health`** — both expose a health check endpoint returning `{"status":"ok"}`.

## Running a project

```bash
cd go-bookshop   # or go-library
go run .
```

See each project's own `README.md` for its full endpoint list, data model, and example
`curl` requests.

## Repo layout

```
go-bookshop/   Books + orders + loyalty programme API (MVC layout)
go-library/    Minimal books CRUD API (flat layout, Dockerfile included)
```
