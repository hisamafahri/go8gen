package http

import (
	"net/http"

	"github.com/go-chi/render"

	"abc/internal/domain/health"
)

type Handler struct {
	useCase health.UseCase
}

func NewHandler(useCase health.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) Liveness(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
}

func (h *Handler) Readiness(w http.ResponseWriter, r *http.Request) {
	err := h.useCase.Readiness()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, nil)
	}
	render.Status(r, http.StatusOK)
}