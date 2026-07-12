package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"bookshop/internal/store"
	"bookshop/internal/view"
)

// CreateOrderRequest is the controller's input binding for POST /orders.
type CreateOrderRequest struct {
	Customer   string `json:"customer"`
	CustomerID string `json:"customer_id,omitempty"`
	Items      []struct {
		BookID   string `json:"book_id"`
		Quantity int    `json:"quantity"`
	} `json:"items"`
}

// CreateCustomerRequest is the controller's input binding for POST /customers.
type CreateCustomerRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type Handler struct {
	store *store.Store
}

func New(store *store.Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("GET /books", h.listBooks)
	mux.HandleFunc("GET /books/{id}", h.getBook)
	mux.HandleFunc("GET /orders", h.listOrders)
	mux.HandleFunc("GET /orders/{id}", h.getOrder)
	mux.HandleFunc("POST /orders", h.createOrder)
	mux.HandleFunc("GET /customers", h.listCustomers)
	mux.HandleFunc("GET /customers/{id}", h.getCustomer)
	mux.HandleFunc("POST /customers", h.createCustomer)
	return mux
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	view.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) listBooks(w http.ResponseWriter, r *http.Request) {
	view.JSON(w, http.StatusOK, h.store.ListBooks())
}

func (h *Handler) getBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	book, ok := h.store.GetBook(id)
	if !ok {
		view.Error(w, http.StatusNotFound, "book not found")
		return
	}
	view.JSON(w, http.StatusOK, book)
}

func (h *Handler) listOrders(w http.ResponseWriter, r *http.Request) {
	view.JSON(w, http.StatusOK, h.store.ListOrders())
}

func (h *Handler) getOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	order, ok := h.store.GetOrder(id)
	if !ok {
		view.Error(w, http.StatusNotFound, "order not found")
		return
	}
	view.JSON(w, http.StatusOK, order)
}

func (h *Handler) createOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		view.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	order, err := h.store.CreateOrder(req.Customer, req.CustomerID, req.Items)
	if err != nil {
		if errors.Is(err, store.ErrBookNotFound) || errors.Is(err, store.ErrOutOfStock) || errors.Is(err, store.ErrInvalidQuantity) {
			view.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		log.Printf("create order error: %v", err)
		view.Error(w, http.StatusInternalServerError, "internal error")
		return
	}
	view.JSON(w, http.StatusCreated, order)
}

func (h *Handler) listCustomers(w http.ResponseWriter, r *http.Request) {
	view.JSON(w, http.StatusOK, h.store.ListCustomers())
}

func (h *Handler) getCustomer(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c, ok := h.store.GetCustomer(id)
	if !ok {
		view.Error(w, http.StatusNotFound, "customer not found")
		return
	}
	view.JSON(w, http.StatusOK, c)
}

func (h *Handler) createCustomer(w http.ResponseWriter, r *http.Request) {
	var req CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		view.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	c, err := h.store.CreateCustomer(req.Name, req.Email, req.Phone)
	if err != nil {
		view.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	view.JSON(w, http.StatusCreated, c)
}
