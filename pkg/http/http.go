package http

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
)

func CreateAndRunServer(r chi.Router, addr string) error {
	httpServer := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	return httpServer.ListenAndServe()
}

func CreateServerWithShutdown(r chi.Router, addr string) (*http.Server, error) {
	httpServer := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("HTTP server failed: %v", err)
		}
	}()

	return httpServer, nil
}

func ShutdownServer(ctx context.Context, server *http.Server) error {
	return server.Shutdown(ctx)
}
