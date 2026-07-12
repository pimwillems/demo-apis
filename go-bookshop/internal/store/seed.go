package store

import (
	"fmt"
	"math/rand"
	"time"

	"bookshop/internal/models"
)

func (s *Store) seed() {
	rand.Seed(time.Now().UnixNano())
	s.seedBooks(100)
	s.seedOrders()
	s.seedCustomers()
}

func (s *Store) seedBooks(n int) {
	adjectives := []string{"The", "A", "Silent", "Hidden", "Last", "Broken", "Eternal", "Lost", "Secret", "Burning", "Crimson", "Forgotten", "Wild", "Quiet", "Distant", "Golden", "Frozen", "Ancient", "Restless", "Hidden"}
	nouns := []string{"Mountain", "River", "Shadow", "Empire", "Garden", "Storm", "Library", "Horizon", "Machine", "Forest", "Ocean", "Castle", "Comet", "Labyrinth", "Republic", "Wanderer", "Algorithm", "Kingdom", "Prophecy", "Circuit"}
	authors := []string{"Jane Doe", "John Smith", "A. Rivera", "M. Chen", "L. Patel", "R. Okafor", "S. Müller", "T. Nakamura", "E. Larsson", "C. Rossi", "B. Haddad", "N. Petrov", "K. Andersson", "F. Silva", "D. Kim"}
	genres := []string{"Fiction", "Sci-Fi", "Fantasy", "History", "Programming", "Mystery", "Biography", "Philosophy", "Romance", "Thriller"}

	for i := 1; i <= n; i++ {
		adj := adjectives[rand.Intn(len(adjectives))]
		noun := nouns[rand.Intn(len(nouns))]
		author := authors[rand.Intn(len(authors))]
		genre := genres[rand.Intn(len(genres))]
		price := 9.99 + rand.Float64()*40
		price = float64(int(price*100)) / 100
		stock := rand.Intn(120) + 5
		isbn := fmt.Sprintf("978%010d", (i*137+41)%10000000000)

		s.books[fmt.Sprintf("bk-%d", i)] = models.Book{
			ID:          fmt.Sprintf("bk-%d", i),
			ISBN:        isbn,
			Title:       fmt.Sprintf("%s %s", adj, noun),
			Author:      author,
			Description: fmt.Sprintf("A %s tale of %s by %s.", genre, noun, author),
			Price:       price,
			Stock:       stock,
			Genre:       genre,
		}
	}
}

func (s *Store) seedOrders() {
	s.orders["ord-seed1"] = models.Order{
		ID:        "ord-seed1",
		Customer:  "Alice",
		Items:     []models.OrderItem{{BookID: "bk-1", Title: s.books["bk-1"].Title, Quantity: 1, Price: s.books["bk-1"].Price}},
		Total:     s.books["bk-1"].Price,
		CreatedAt: time.Now().Add(-48 * time.Hour),
		Status:    "shipped",
	}
	s.orders["ord-seed2"] = models.Order{
		ID:       "ord-seed2",
		Customer: "Bob",
		Items: []models.OrderItem{
			{BookID: "bk-4", Title: s.books["bk-4"].Title, Quantity: 2, Price: s.books["bk-4"].Price},
			{BookID: "bk-5", Title: s.books["bk-5"].Title, Quantity: 1, Price: s.books["bk-5"].Price},
		},
		Total:     s.books["bk-4"].Price*2 + s.books["bk-5"].Price,
		CreatedAt: time.Now().Add(-24 * time.Hour),
		Status:    "pending",
	}
}

func (s *Store) seedCustomers() {
	s.customers["cust-seed1"] = models.Customer{ID: "cust-seed1", Name: "Alice Anderson", Email: "alice@example.com", Phone: "+1-555-0101", Points: 120, CreatedAt: time.Now().Add(-720 * time.Hour)}
	s.customers["cust-seed2"] = models.Customer{ID: "cust-seed2", Name: "Bob Brown", Email: "bob@example.com", Phone: "+1-555-0102", Points: 53, CreatedAt: time.Now().Add(-480 * time.Hour)}
	s.customers["cust-seed3"] = models.Customer{ID: "cust-seed3", Name: "Carol Clark", Email: "carol@example.com", Phone: "+1-555-0103", Points: 0, CreatedAt: time.Now().Add(-24 * time.Hour)}
}
