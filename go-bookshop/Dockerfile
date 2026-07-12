# ---- build stage ----
FROM golang:1.22-alpine AS build

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/bookshop .

# ---- runtime stage ----
FROM alpine:3.20

WORKDIR /app

COPY --from=build /out/bookshop /app/bookshop

EXPOSE 8080

ENTRYPOINT ["/app/bookshop"]
