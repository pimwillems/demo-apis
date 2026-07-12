package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	store := NewStore(seedBooks())

	mux := http.NewServeMux()
	mux.HandleFunc("GET /books", store.handleListBooks)
	mux.HandleFunc("GET /books/{id}", store.handleGetBook)
	mux.HandleFunc("POST /books", store.handleCreateBook)
	mux.HandleFunc("GET /health", handleHealth)

	addr := ":" + port()
	log.Printf("books API listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func port() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}
	return "8080"
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func seedBooks() []Book {
	return []Book{
		{ID: 1, Title: "The Hobbit", Author: "J.R.R. Tolkien", Genre: "Fantasy", ISBN: "978-0-618-96863-3"},
		{ID: 2, Title: "Dune", Author: "Frank Herbert", Genre: "Science Fiction", ISBN: "978-0-441-17271-9"},
		{ID: 3, Title: "The Name of the Rose", Author: "Umberto Eco", Genre: "Historical Mystery", ISBN: "978-0-15-144647-6"},
	}
}
