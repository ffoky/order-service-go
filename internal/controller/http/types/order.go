package types

import (
	"WBTECH_L0/internal/domain"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type GetOrderRequest struct {
	OrderID string `json:"order_id" validate:"required"`
}

func CreateGetOrderRequest(r *http.Request) (GetOrderRequest, error) {
	orderID := chi.URLParam(r, "id")
	if orderID == "" {
		return GetOrderRequest{}, errors.New("order_id is required")
	}
	return GetOrderRequest{OrderID: orderID}, nil
}

type GetOrderResponse struct {
	*domain.Order
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}

type ValidationErrorResponse struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields"`
}
