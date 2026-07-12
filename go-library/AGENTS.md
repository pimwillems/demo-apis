# AGENTS.md

Instructions for AI agents (and humans) working on this repository, plus a record of how it was built.

## What this project is

A minimal JSON REST API for books in Go. Standard library only — no web frameworks, no external dependencies. Data lives in memory and resets on restart. See `README.md` for the full endpoint documentation.

## How this project was built

This project was created in a single Claude Code session (2026-07-11):

1. **Module setup** — `go mod init github.com/pimwillems/go-booksapi` in an empty directory (Go 1.26).
2. **Implementation** — split across two files in `package main`:
   - `main.go`: route registration on a plain `http.NewServeMux`, `PORT` env handling (default `8080`), seed data, `/health` handler.
   - `books.go`: the `Book` struct, a mutex-guarded in-memory `Store`, the three book handlers, and shared `writeJSON`/`writeError` helpers.
3. **Verification** — `go vet`, `go build`, then the server was run on a test port and every endpoint (including all error paths: bad id, unknown id, invalid body, missing fields, wrong method) was exercised with `curl`.
4. **Documentation** — `README.md` (API usage) and this file were written last, after the behavior was confirmed working.

## Conventions to follow

- **Stdlib only.** Do not add third-party dependencies (routers, frameworks, ORMs) without the user explicitly asking. Go 1.22+ `ServeMux` patterns (`"GET /books/{id}"` + `r.PathValue("id")`) cover routing needs.
- **Every response body is JSON**, written through `writeJSON`. Errors go through `writeError` and always have the shape `{"error": "<message>"}` with an appropriate 4xx status.
- **Concurrency safety.** All access to `Store.books` and `Store.nextID` must hold `Store.mu` (`RLock` for reads, `Lock` for writes). List responses copy the slice before releasing the lock.
- **Strict request parsing.** `POST` bodies are decoded with `DisallowUnknownFields`, and all four book fields are required (whitespace-only values count as missing). Keep new write endpoints equally strict.
- **IDs are server-assigned** via `Store.nextID`; never trust an `id` from a request body.
- **JSON field names are lowercase** (`id`, `title`, `author`, `genre`, `isbn`) via struct tags.

## How to verify changes

There are no automated tests yet (adding `httptest`-based handler tests would be a good first improvement). Until then, verify by hand:

```sh
go vet ./...
go build ./...
PORT=8091 go run .   # in one shell
```

Then exercise the endpoints, including error paths:

```sh
curl -s http://localhost:8091/books
curl -s http://localhost:8091/books/2
curl -s -w ' [%{http_code}]' http://localhost:8091/books/99        # expect 404
curl -s -w ' [%{http_code}]' http://localhost:8091/books/abc       # expect 400
curl -s -w ' [%{http_code}]' -X POST http://localhost:8091/books \
  -H 'Content-Type: application/json' \
  -d '{"title":"T","author":"A","genre":"G","isbn":"I"}'           # expect 201
curl -s -w ' [%{http_code}]' -X POST http://localhost:8091/books \
  -H 'Content-Type: application/json' -d '{"title":"only"}'        # expect 400
```

## When extending the API

- Register new routes in `main.go`; keep handler logic in a resource-specific file (like `books.go`).
- Update `README.md` whenever endpoints, fields, or status codes change — it is the API contract.
- If persistence is ever added, keep the `Store` interface surface (list/get/create) so handlers don't need rewriting.
