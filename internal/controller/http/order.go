package http

import (
	"WBTECH_L0/internal/controller/http/types"
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

// getOrderHandler получает заказ по ID
// @Summary      Получить заказ по ID
// @Description  Возвращает информацию о заказе по его уникальному идентификатору
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID заказа"  example("b563feb7b2b84b6test")
// @Success      200  {object}  types.GetOrderResponse
// @Failure      400  {object}  types.ErrorResponse
// @Failure      404  {object}  types.ErrorResponse
// @Failure      500  {object}  types.ErrorResponse
// @Router       /order/{id} [get]
func (h *OrderHandler) getOrderHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	request, err := types.CreateGetOrderRequest(r)
	if err != nil {
		h.sendErrorResponse(w, "Invalid request: "+err.Error(), http.StatusBadRequest, "INVALID_REQUEST")
		return
	}

	fromCache := h.service.IsInCache(r.Context(), request.OrderID)

	orderObject, err := h.service.Get(r.Context(), request.OrderID)
	if err != nil {
		logrus.Errorf("Failed to get order %s: %v", request.OrderID, err)
		h.sendErrorResponse(w, "Order not found", http.StatusNotFound, "ORDER_NOT_FOUND")
		return
	}

	source := "database"
	if fromCache {
		source = "cache"
	}

	response := types.GetOrderResponse{
		Order: orderObject,
	}

	h.sendSuccessResponse(w, response)

	logrus.Infof("Order %s received from %s in %s", request.OrderID, source, time.Since(startTime))
}

func (h *OrderHandler) sendSuccessResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logrus.Errorf("Failed to encode success response: %v", err)
		h.sendErrorResponse(w, "Internal server error", http.StatusInternalServerError, "ENCODING_ERROR")
	}
}

func (h *OrderHandler) sendErrorResponse(w http.ResponseWriter, message string, statusCode int, code string) {
	response := types.ErrorResponse{
		Error: message,
		Code:  code,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logrus.WithError(err).Error("failed to encode error response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *OrderHandler) WithOrderHandlers(r chi.Router) {
	r.Route("/order", func(r chi.Router) {
		r.Get("/{id}", h.getOrderHandler)
	})
}
