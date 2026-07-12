package main

import (
	"log"
	"net/http"

	"bookshop/internal/handlers"
	"bookshop/internal/store"
)

func main() {
	s := store.New()
	h := handlers.New(s)

	addr := ":8080"
	log.Printf("bookshop API listening on %s", addr)
	if err := http.ListenAndServe(addr, h.Routes()); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
