package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// Book is the resource served by the API.
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Genre  string `json:"genre"`
	ISBN   string `json:"isbn"`
}

// Store holds books in memory, safe for concurrent handlers.
type Store struct {
	mu     sync.RWMutex
	books  []Book
	nextID int
}

func NewStore(seed []Book) *Store {
	nextID := 1
	for _, b := range seed {
		if b.ID >= nextID {
			nextID = b.ID + 1
		}
	}
	return &Store{books: seed, nextID: nextID}
}

func (s *Store) handleListBooks(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	books := make([]Book, len(s.books))
	copy(books, s.books)
	s.mu.RUnlock()

	writeJSON(w, http.StatusOK, books)
}

func (s *Store) handleGetBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "id must be an integer")
		return
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, b := range s.books {
		if b.ID == id {
			writeJSON(w, http.StatusOK, b)
			return
		}
	}
	writeError(w, http.StatusNotFound, "book not found")
}

func (s *Store) handleCreateBook(w http.ResponseWriter, r *http.Request) {
	var in Book
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body: "+err.Error())
		return
	}

	var missing []string
	for _, f := range []struct{ name, value string }{
		{"title", in.Title},
		{"author", in.Author},
		{"genre", in.Genre},
		{"isbn", in.ISBN},
	} {
		if strings.TrimSpace(f.value) == "" {
			missing = append(missing, f.name)
		}
	}
	if len(missing) > 0 {
		writeError(w, http.StatusBadRequest, "missing required fields: "+strings.Join(missing, ", "))
		return
	}

	s.mu.Lock()
	in.ID = s.nextID
	s.nextID++
	s.books = append(s.books, in)
	s.mu.Unlock()

	writeJSON(w, http.StatusCreated, in)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
