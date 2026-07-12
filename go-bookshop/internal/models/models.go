package models

import "time"

type Book struct {
	ID          string  `json:"id"`
	ISBN        string  `json:"isbn"`
	Title       string  `json:"title"`
	Author      string  `json:"author"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Genre       string  `json:"genre"`
}

type Customer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Points    int       `json:"points"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderItem struct {
	BookID   string  `json:"book_id"`
	Title    string  `json:"title"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type Order struct {
	ID         string      `json:"id"`
	Customer   string      `json:"customer"`
	CustomerID string      `json:"customer_id,omitempty"`
	Items      []OrderItem `json:"items"`
	Total      float64     `json:"total"`
	CreatedAt  time.Time   `json:"created_at"`
	Status     string      `json:"status"`
}

