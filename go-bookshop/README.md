# Bookshop API

A small, fully self-contained **Go** REST API for a book webshop. It is meant as a
backend demo to build a frontend against, with data about **books**, **orders**,
and a **customer loyalty programme**.

No external dependencies — it uses only the Go standard library (`net/http`,
`encoding/json`). All data lives in memory and is reset on every restart, which is
ideal for a frontend demo.

## Requirements

- Go 1.22 or newer

## Run

```bash
go run .
```

The server listens on `http://localhost:8080`.

Build a binary instead:

```bash
go build -o bookshop.exe .   # or just `go build .` on Linux/macOS
./bookshop
```

## Data model

- **Book** — `id` (`bk-N`), `isbn`, `title`, `author`, `description`, `price`,
  `stock`, `genre`. Seeded with 100 generated books.
- **Order** — `id` (`ord-` + 10 random base62 chars), `customer`, optional
  `customer_id`, `items`, `total`, `created_at`, `status`.
- **Customer** — `id` (`cust-` + 10 random base62 chars), `name`, `email`,
  `phone`, `points`, `created_at`.

When an order is placed with a valid `customer_id`, the customer earns **1 loyalty
point per €1 spent** (rounded down).

## Endpoints

All responses are JSON.

### Health

| Method | Path       | Description        |
| ------ | ---------- | ------------------ |
| `GET`  | `/health`  | Liveness check. Returns `{"status":"ok"}`. |

### Books

| Method | Path           | Description                          |
| ------ | -------------- | ------------------------------------ |
| `GET`  | `/books`       | List all books.                      |
| `GET`  | `/books/{id}`  | Get a single book by id (e.g. `bk-1`). Returns `404` if not found. |

### Orders

| Method | Path           | Description                          |
| ------ | -------------- | ------------------------------------ |
| `GET`  | `/orders`      | List all orders.                     |
| `GET`  | `/orders/{id}` | Get a single order by id. Returns `404` if not found. |
| `POST` | `/orders`      | Create an order.                     |

`POST /orders` request body:

```json
{
  "customer": "Zoe",
  "customer_id": "cust-abc123XYZ",
  "items": [
    { "book_id": "bk-7", "quantity": 2 }
  ]
}
```

- `customer` is a free-text name (required-ish; at least one item is required).
- `customer_id` is optional; if it matches an existing customer, loyalty points
  are awarded.
- Validates that books exist, are in stock, and quantities are positive. Returns
  `400` with a descriptive error otherwise. On success returns `201 Created`
  with the created order (stock is decremented).

### Customers (loyalty programme)

| Method | Path              | Description                                |
| ------ | ----------------- | ------------------------------------------ |
| `GET`  | `/customers`      | List all customers with their points.      |
| `GET`  | `/customers/{id}` | Get a single customer by id. `404` if missing. |
| `POST` | `/customers`      | Register a new customer.                   |

`POST /customers` request body:

```json
{
  "name": "Carol Clark",
  "email": "carol@example.com",
  "phone": "+1-555-0103"
}
```

`name` and `email` are required; `phone` is optional. Returns `201 Created` with
the created customer (starting at 0 points).

## Example requests

```bash
# Health
curl http://localhost:8080/health

# List books
curl http://localhost:8080/books

# Create an order for an existing loyalty customer
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{"customer_id":"cust-seed3","items":[{"book_id":"bk-7","quantity":1}]}'

# Register a new customer
curl -X POST http://localhost:8080/customers \
  -H "Content-Type: application/json" \
  -d '{"name":"Dave","email":"dave@example.com","phone":"+1-555-0104"}'
```

## Project layout

This project follows an MVC-style separation:

```
main.go                       server bootstrap (wires Model -> Controller -> View)
internal/models/models.go     Model: domain structs (Book, Order, OrderItem, Customer)
internal/store/store.go       Model: in-memory store + business logic + error sentinels
internal/store/seed.go        Model: seed data (100 books, sample orders, customers)
internal/store/id.go          Model: short unique id generator (base62)
internal/view/view.go         View: JSON response rendering (JSON / Error)
internal/handlers/handlers.go Controller: HTTP handlers, routing, request binding
```

## Deployment

This app is deployed together with the other demo APIs in this repo as a single
container — see the root [`README.md`](../README.md#deployment). When deployed that way,
these endpoints are reachable under the `/go-bookshop` prefix (e.g. `<url>/go-bookshop/books`).
Standalone (`go run .`) is unaffected — it still serves from `/books` on `:8080`.
