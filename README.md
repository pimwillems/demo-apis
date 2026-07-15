# demo-apis

A collection of small, self-contained demo REST APIs used for teaching students. Each
subproject is an independent Go module — no shared code, no monorepo tooling. Pick a
directory, `cd` in, and follow its own README.

## Projects

| Project | Description |
| --- | --- |
| [`go-bookshop/`](go-bookshop/) | Book webshop API with books, orders, and a customer loyalty programme. MVC-style layout (`internal/models`, `internal/store`, `internal/view`, `internal/handlers`). |
| [`go-library/`](go-library/) | Minimal books CRUD API (list/get/create). Flat two-file layout. |

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
`curl` requests. Standalone, each app always listens at its own root paths
(`/books`, `/health`, ...) on `:8080`.

## Repo layout

```
go-bookshop/   Books + orders + loyalty programme API (MVC layout)
go-library/    Minimal books CRUD API (flat layout)
Dockerfile     Builds every app above + Caddy into one deployable image
Caddyfile      Path-routes /go-<name>/* to each app inside the image
entrypoint.sh  Starts every app + Caddy when the container runs
```

## Deployment

All demo APIs in this repo are deployed together as **one container, one URL**. The root
`Dockerfile` builds every app's binary plus [Caddy](https://caddyserver.com/), and
`entrypoint.sh` starts each app on its own internal port while Caddy listens on `:8080`
and path-routes requests to the right app, stripping the prefix:

| Public path | Routed to |
| --- | --- |
| `<url>/go-library/*` | `go-library` (internal port `8081`) |
| `<url>/go-bookshop/*` | `go-bookshop` (internal port `8082`) |
| `<url>/health` | Container-level health check (static `ok`) |

So `<url>/go-library/books` hits `go-library`'s `GET /books` handler, `<url>/go-bookshop/orders`
hits `go-bookshop`'s `GET /orders` handler, and so on — every endpoint documented in each
project's own README is reachable the same way, just under its app's prefix.

Build and run locally:

```bash
docker build -t demo-apis .
docker run --rm -p 8080:8080 demo-apis
curl localhost:8080/go-library/books
curl localhost:8080/go-bookshop/books
```

### Deploying on Coolify

Create **one** application resource pointing at this repo:

- Build pack: **Dockerfile** (uses the root `Dockerfile`)
- Port: `8080`
- Health check path: `/health`

No "Base Directory" or per-app resource is needed — one deploy serves every demo API in
the repo under its own path prefix.

### Adding a new demo API

See `AGENTS.md` for the exact steps — in short: drop a new self-contained Go app in its own
folder, add a build stage for it in `Dockerfile`, add a `handle_path /go-<name>/*` block to
`Caddyfile` on the next free internal port, and start it in `entrypoint.sh`. No changes to
any existing app are needed.
