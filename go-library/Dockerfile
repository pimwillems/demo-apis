FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /booksapi .

FROM alpine:3
RUN adduser -D -u 10001 app
USER app
COPY --from=build /booksapi /usr/local/bin/booksapi
EXPOSE 8080
ENTRYPOINT ["booksapi"]
