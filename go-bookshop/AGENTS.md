# AGENTS.md

Guidance for LLM/agentic tools working in this repository.

## What this project is

A deliberately small, **fully standard-library Go** REST API simulating a book
webshop, built as a backend for a frontend demo (taught to students). It must stay
dependency-free — do **not** add external frameworks (no gin, echo, chi, sqlx, etc.)
or module dependencies unless explicitly asked.

## Hard constraints

- **Go 1.22+** (uses `http.Request.PathValue` and method-based `ServeMux` routing
  like `mux.HandleFunc("GET /books/{id}", ...)`). The `go.mod` declares `go 1.22`.
- **No external dependencies.** Standard library only.
- **In-memory data store** that resets on every restart. Do not introduce a
  database unless the user asks — this is intentional for a demo.
- All handlers return JSON; errors use `{"error":"..."}`.

## Architecture / where things live (MVC)

The project follows an MVC-style layering:

- **Model** (`internal/models`, `internal/store`) — domain data and business rules.
  - `internal/models/models.go` — domain structs (`Book`, `Order`,
    `OrderItem`, `Customer`) with JSON tags. Pure data, no behaviour.
  - `internal/store/store.go` — the `Store` (thread-safe via `sync.RWMutex`),
    data access, business logic, and error sentinels (`ErrBookNotFound`,
    `ErrOrderNotFound`, `ErrCustomerNotFound`, `ErrOutOfStock`,
    `ErrInvalidQuantity`).
  - `internal/store/seed.go` — seed data. `seedBooks(n)` generates `n` books
    programmatically (currently 100). Sample orders and 3 customers are seeded.
    Book ids are stable `bk-N`; orders/customers use short random ids.
  - `internal/store/id.go` — `newShortID(prefix)` generates 10-char base62 ids
    via `crypto/rand` (e.g. `ord-s8qU6b2I1v`). Used for orders and customers.
- **View** (`internal/view/view.go`) — renders HTTP responses (`JSON`, `Error`).
  Knows nothing about routing or business logic.
- **Controller** (`internal/handlers/handlers.go`) — HTTP handlers and route
  table. Binds request DTOs (`CreateOrderRequest`, `CreateCustomerRequest`),
  calls the Model (store), and renders via the View. Keeps no business logic.

`main.go` wires the three layers together (Model -> Controller -> View) and
boots the server on `:8080`.

## Conventions

- ID scheme: books `bk-N`, orders `ord-<10 base62>`, customers `cust-<10 base62>`.
- Handlers call `writeJSON(w, status, v)` / `writeError(w, status, msg)`.
- Routing uses the stdlib `GET /path/{param}` syntax; read params with
  `r.PathValue("param")`.
- Store methods take/return models and are the single source of business rules
  (stock checks, point awards). Keep logic in the store, not handlers.
- Loyalty rule: placing an order with a valid `customer_id` awards 1 point per
  full €1 of the order total (`int(total)`).

## Endpoints (current)

`GET /health`, `GET /books`, `GET /books/{id}`, `GET /orders`, `GET /orders/{id}`,
`POST /orders`, `GET /customers`, `GET /customers/{id}`, `POST /customers`.

## Build / verify before finishing

```bash
go build ./... && go vet ./...
```

Run with `go run .` (port 8080). If a previous process is still bound to `:8080`,
kill it before re-running (PowerShell: `Get-NetTCPConnection -LocalPort 8080 |
ForEach-Object { Stop-Process -Id $_.OwningProcess -Force }`).

## History of this conversation (for context)

1. Initial request: build a pure-Go book webshop demo API with books + orders.
2. Added `isbn` to Book and a loyalty programme (customers with personal details
   + points; orders award points to a linked `customer_id`).
3. Replaced sequential numeric order/customer ids with short random base62 ids
   (UUID-like but short); expanded seed to 100 generated books.
4. Added `README.md` (endpoints + run instructions) and this `AGENTS.md`.

## Open extension ideas (not yet implemented, may be requested)

- Persistence (SQLite or JSON file) instead of in-memory.
- `POST /books` to add books; `PATCH` to update stock.
- Redeeming loyalty points at checkout.
- Pagination / filtering on `GET /books`.
