package types

import (
	"net/http"
)

type GetOrderHandlerRequest struct {
	UID string `json:"order_uid"`
}

func CreateGetOrderByUIDHanlderRequest(r *http.Request) (GetOrderHandlerRequest, error) {
	uid := r.URL.Query().Get("uid")
	if uid == "" {
		return GetOrderHandlerRequest{}, http.ErrNoLocation
	}
	return GetOrderHandlerRequest{UID: uid}, nil
}
