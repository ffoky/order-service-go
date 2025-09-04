package http

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type StaticHandler struct{}

func NewStaticHandler() *StaticHandler {
	return &StaticHandler{}
}

func (h *StaticHandler) WithStaticHandlers(r chi.Router) {
	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Get("/", h.serveIndexHandler)
}

func (h *StaticHandler) serveIndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}
