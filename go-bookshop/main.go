package main

import (
	"log"
	"net/http"
	"os"

	"bookshop/internal/handlers"
	"bookshop/internal/store"
)

func main() {
	s := store.New()
	h := handlers.New(s)

	addr := ":" + port()
	log.Printf("bookshop API listening on %s", addr)
	if err := http.ListenAndServe(addr, h.Routes()); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func port() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}
	return "8080"
}
