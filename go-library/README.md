# go-booksapi

A small JSON REST API for books, written in Go using only the standard library (`net/http`). Books are stored in memory and the server starts with three seeded books, so it works out of the box — note that any books you add are lost when the server restarts.

## Requirements

- Go 1.22 or newer (uses `http.ServeMux` method/wildcard route patterns; built and tested with Go 1.26)

## Running

```sh
go run .
```

The server listens on port `8080` by default. Override it with the `PORT` environment variable:

```sh
PORT=3000 go run .
```

## The Book resource

```json
{
  "id": 1,
  "title": "The Hobbit",
  "author": "J.R.R. Tolkien",
  "genre": "Fantasy",
  "isbn": "978-0-618-96863-3"
}
```

| Field    | Type   | Notes                                   |
| -------- | ------ | --------------------------------------- |
| `id`     | int    | Assigned by the server, auto-increments |
| `title`  | string | Required on create                      |
| `author` | string | Required on create                      |
| `genre`  | string | Required on create                      |
| `isbn`   | string | Required on create                      |

## Endpoints

### `GET /books` — list all books

```sh
curl http://localhost:8080/books
```

Returns `200 OK` with a JSON array of all books.

### `GET /books/{id}` — get one book

```sh
curl http://localhost:8080/books/2
```

- `200 OK` with the book as a JSON object.
- `400 Bad Request` if `{id}` is not an integer.
- `404 Not Found` if no book has that id.

### `POST /books` — create a book

```sh
curl -X POST http://localhost:8080/books \
  -H 'Content-Type: application/json' \
  -d '{"title":"1984","author":"George Orwell","genre":"Dystopian","isbn":"978-0-452-28423-4"}'
```

- `201 Created` with the stored book, including its server-assigned `id`.
- `400 Bad Request` if the body is not valid JSON, contains unknown fields, or is missing any of `title`, `author`, `genre`, `isbn`.

### `GET /health` — health check

```sh
curl http://localhost:8080/health
```

Returns `200 OK` with `{"status":"ok"}`.

## Error responses

All errors are JSON objects with a single `error` field, for example:

```json
{"error": "book not found"}
```

Requests using an unsupported method on a known path (e.g. `DELETE /books`) get a `405 Method Not Allowed` from the router.

## Deployment

The repo includes a multi-stage `Dockerfile` that produces a small Alpine-based image running as a non-root user on port 8080. It is set up for [Coolify](https://coolify.io) (build pack: Dockerfile, exposed port 8080, health check `GET /health`), but works with any container platform:

```sh
docker build -t booksapi .
docker run -p 8080:8080 booksapi
```

## Project layout

| File       | Purpose                                                        |
| ---------- | -------------------------------------------------------------- |
| `main.go`  | Server setup: routes, port handling, seed data, health handler |
| `books.go` | `Book` type, in-memory `Store`, HTTP handlers, JSON helpers    |
