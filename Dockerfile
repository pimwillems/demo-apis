# Builds every demo API in this repo into a single image, fronted by Caddy,
# which path-routes /go-<name>/* to each app's internal port. See AGENTS.md
# for the recipe to add a new demo API here.

# ---- build: go-library ----
FROM golang:1.26-alpine AS build-go-library
WORKDIR /src
COPY go-library/go.mod ./
RUN go mod download
COPY go-library/ ./
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/go-library .

# ---- build: go-bookshop ----
FROM golang:1.26-alpine AS build-go-bookshop
WORKDIR /src
COPY go-bookshop/go.mod ./
RUN go mod download
COPY go-bookshop/ ./
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/go-bookshop .

# ---- runtime ----
FROM caddy:2-alpine

WORKDIR /app
COPY --from=build-go-library   /out/go-library   /app/bin/go-library
COPY --from=build-go-bookshop  /out/go-bookshop  /app/bin/go-bookshop
COPY Caddyfile /app/Caddyfile
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8080
ENTRYPOINT ["/entrypoint.sh"]
