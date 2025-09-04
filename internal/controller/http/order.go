package http

import (
	"WBTECH_L0/internal/usecases"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type OrderHandler struct {
	service usecases.Order
}

func NewOrderHandler(service usecases.Order) *OrderHandler {
	return &OrderHandler{service: service}
}

type OrderResponse struct {
	Order interface{} `json:"order,omitempty"`
	Error string      `json:"error,omitempty"`
}

func (h *OrderHandler) getOrderHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	orderID := chi.URLParam(r, "id")
	if orderID == "" {
		h.sendErrorResponse(w, "order_id is required", http.StatusBadRequest, false, startTime)
		return
	}

	fromCache := h.service.IsInCache(r.Context(), orderID)

	orderObject, err := h.service.Get(r.Context(), orderID)
	if err != nil {
		logrus.Errorf("Failed to get order %s: %v", orderID, err)
		h.sendErrorResponse(w, "order not found", http.StatusNotFound, fromCache, startTime)
		return
	}

	response := OrderResponse{
		Order: orderObject,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logrus.Errorf("Failed to encode response: %v", err)
		h.sendErrorResponse(w, "internal server error", http.StatusInternalServerError, fromCache, startTime)
		return
	}

	source := "database"
	if fromCache {
		source = "cache"
	}
	logrus.Infof("Order %s recieved from %s in %s", orderID, source, time.Since(startTime))
}

func (h *OrderHandler) getCacheStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats := h.service.GetCacheStats()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(stats); err != nil {
		logrus.Errorf("Failed to encode cache stats: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *OrderHandler) sendErrorResponse(w http.ResponseWriter, message string, statusCode int, fromCache bool, startTime time.Time) {
	response := OrderResponse{
		Error: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func (h *OrderHandler) WithOrderHandlers(r chi.Router) {
	r.Route("/order", func(r chi.Router) {
		r.Get("/{id}", h.getOrderHandler)
	})
}
