package http

import (
	"WBTECH_L0/internal/usecases/service"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type OrderHandler struct {
	service *service.Order
}

func NewOrderHandler(service *service.Order) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) getOrderHandler(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "id")

	orderObject, err := h.service.Get(r.Context(), orderID)
	if err != nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(orderObject); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *OrderHandler) WithOrderHandlers(r chi.Router) {
	r.Route("/order", func(r chi.Router) {
		r.Get("/{id}", h.getOrderHandler)
	})
}
