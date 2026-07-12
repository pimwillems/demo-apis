package store

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"bookshop/internal/models"
)

var (
	ErrBookNotFound     = errors.New("book not found")
	ErrOrderNotFound    = errors.New("order not found")
	ErrCustomerNotFound = errors.New("customer not found")
	ErrOutOfStock       = errors.New("not enough stock for book")
	ErrInvalidQuantity  = errors.New("quantity must be greater than zero")
)

type Store struct {
	mu       sync.RWMutex
	books    map[string]models.Book
	orders   map[string]models.Order
	customers map[string]models.Customer
}

func New() *Store {
	s := &Store{
		books:     make(map[string]models.Book),
		orders:    make(map[string]models.Order),
		customers: make(map[string]models.Customer),
	}
	s.seed()
	return s
}

func (s *Store) ListBooks() []models.Book {
	s.mu.RLock()
	defer s.mu.RUnlock()
	books := make([]models.Book, 0, len(s.books))
	for _, b := range s.books {
		books = append(books, b)
	}
	return books
}

func (s *Store) GetBook(id string) (models.Book, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.books[id]
	return b, ok
}

func (s *Store) ListOrders() []models.Order {
	s.mu.RLock()
	defer s.mu.RUnlock()
	orders := make([]models.Order, 0, len(s.orders))
	for _, o := range s.orders {
		orders = append(orders, o)
	}
	return orders
}

func (s *Store) GetOrder(id string) (models.Order, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	o, ok := s.orders[id]
	return o, ok
}

func (s *Store) CreateOrder(customer, customerID string, items []struct {
	BookID   string `json:"book_id"`
	Quantity int    `json:"quantity"`
}) (models.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(items) == 0 {
		return models.Order{}, errors.New("order must contain at least one item")
	}

	orderItems := make([]models.OrderItem, 0, len(items))
	total := 0.0

	for _, it := range items {
		if it.Quantity <= 0 {
			return models.Order{}, ErrInvalidQuantity
		}
		book, ok := s.books[it.BookID]
		if !ok {
			return models.Order{}, fmt.Errorf("%w: %s", ErrBookNotFound, it.BookID)
		}
		if book.Stock < it.Quantity {
			return models.Order{}, fmt.Errorf("%w: %s (requested %d, available %d)", ErrOutOfStock, book.Title, it.Quantity, book.Stock)
		}
		book.Stock -= it.Quantity
		s.books[it.BookID] = book

		orderItems = append(orderItems, models.OrderItem{
			BookID:   book.ID,
			Title:    book.Title,
			Quantity: it.Quantity,
			Price:    book.Price,
		})
		total += book.Price * float64(it.Quantity)
	}

	id := newOrderID()
	order := models.Order{
		ID:         id,
		Customer:   customer,
		CustomerID: customerID,
		Items:      orderItems,
		Total:      total,
		CreatedAt:  time.Now(),
		Status:     "pending",
	}
	s.orders[id] = order

	if customerID != "" {
		if c, ok := s.customers[customerID]; ok {
			c.Points += int(total)
			s.customers[customerID] = c
		}
	}
	return order, nil
}

func (s *Store) ListCustomers() []models.Customer {
	s.mu.RLock()
	defer s.mu.RUnlock()
	customers := make([]models.Customer, 0, len(s.customers))
	for _, c := range s.customers {
		customers = append(customers, c)
	}
	return customers
}

func (s *Store) GetCustomer(id string) (models.Customer, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.customers[id]
	return c, ok
}

func (s *Store) CreateCustomer(name, email, phone string) (models.Customer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if name == "" || email == "" {
		return models.Customer{}, errors.New("name and email are required")
	}

	id := newCustomerID()
	c := models.Customer{
		ID:        id,
		Name:      name,
		Email:     email,
		Phone:     phone,
		Points:    0,
		CreatedAt: time.Now(),
	}
	s.customers[id] = c
	return c, nil
}
