# AGENTS.md

Guidance for LLM/agentic tools working in this repository.

## What this repository is

A collection of small, independent demo REST APIs written in Go, used to teach
students. It is **not** a monorepo with shared Go code — each subdirectory is an
unrelated Go module, each with its own `go.mod`, `README.md`, and `AGENTS.md`, and each
still runs standalone via `go run .` on `:8080`.

The one thing that *is* shared at the root is deployment: a single `Dockerfile`,
`Caddyfile`, and `entrypoint.sh` package every app into one container so they can all be
deployed together (see "Deployment architecture" below). There is no root-level Go
module, workspace, or CI — apps never import each other's code.

## Projects

- **`go-bookshop/`** — books, orders, and a customer loyalty programme. MVC-style
  layering (`internal/models`, `internal/store`, `internal/view`,
  `internal/handlers`). See `go-bookshop/AGENTS.md` for its specifics.
- **`go-library/`** — minimal books CRUD API in a flat two-file layout
  (`main.go`, `books.go`). See `go-library/AGENTS.md` for its specifics.

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
- **Reads `PORT` env, defaults to `8080`.** Both apps start on `:8080` when run
  standalone (`go run .`), but must also honor a `PORT` env var so they can be run on a
  different internal port inside the shared deployment container. Copy the existing
  `port()` helper (`go-library/main.go` or `go-bookshop/main.go`) into any new app.

## Deployment architecture

All demo apps in this repo are built into **one Docker image** and run as separate
processes in **one container**, fronted by [Caddy](https://caddyserver.com/) for
path-based routing. This lets the whole repo deploy as a single Coolify application
(see root `README.md` "Deployment" section) while every app's own code stays a
completely standalone `package main` — no shared Go module, no importing between apps.

- **`Dockerfile`** (repo root) — multi-stage build. One build stage per app compiles its
  binary; the final stage is based on `caddy:2-alpine` and copies in every app binary
  plus `Caddyfile` and `entrypoint.sh`.
- **`Caddyfile`** (repo root) — listens on `:8080`, and for each app has a
  `handle_path /go-<name>/*` block that strips the prefix and reverse-proxies to that
  app's internal port (`go-library` → `8081`, `go-bookshop` → `8082`, ...). `/health` is
  a static container-level check.
- **`entrypoint.sh`** (repo root) — starts every app binary in the background on its
  assigned `PORT`, then runs Caddy in the foreground. Uses `wait -n` so the container
  exits (and gets restarted by the orchestrator) if any one app process dies, rather than
  silently serving a degraded container.

## Adding a new demo project

If asked to add a new demo API:

1. Create a new top-level directory (e.g. `go-<name>/`) with its own `go.mod` —
   don't fold it into an existing module.
2. Follow the shared conventions above unless told otherwise, including reading `PORT`
   from the environment (default `8080`).
3. Give it its own `README.md` (endpoints, run instructions, examples) and
   `AGENTS.md` (constraints, architecture, how to verify), matching the style of the
   existing projects.
4. Add a row for it to the root `README.md` project table.
5. Wire it into the shared deployment (all in the repo root):
   - Add a build stage for it in `Dockerfile` (copy an existing app's stage as a
     template) and `COPY` its binary into the final stage.
   - Pick the next free internal port and add a `handle_path /go-<name>/*` block to
     `Caddyfile` reverse-proxying to it.
   - Start it in `entrypoint.sh` with `PORT=<port> /app/bin/go-<name> &`.
   - Add its route to the table in the root `README.md` "Deployment" section.

No existing app's code needs to change when a new one is added.

## Verifying changes

Each project builds and runs independently:

```bash
cd go-bookshop   # or go-library
go build ./... && go vet ./...
go run .
```

Check the project's own `AGENTS.md` for its exact verification steps (test ports,
curl examples, etc.) before considering a change done.

To verify the combined deployment image after touching `Dockerfile`, `Caddyfile`, or
`entrypoint.sh`:

```bash
docker build -t demo-apis .
docker run --rm -p 8080:8080 demo-apis
curl localhost:8080/health
curl localhost:8080/go-library/books
curl localhost:8080/go-bookshop/books
```
