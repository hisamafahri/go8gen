package http

import (
	"github.com/go-chi/chi"

	"abc/internal/domain/health"
)

func RegisterHTTPEndPoints(router *chi.Mux, uc health.UseCase) {
	h := NewHandler(uc)

	router.Route("/health", func(router chi.Router) {
		router.Get("/liveness", h.Liveness)
		router.Get("/readiness", h.Readiness)
	})
}