package handlers

import (
	"net/http"

	"github.com/go-chi/render"

	"goph-keeper/internal/models"
)

type service interface {
	Register(login, password string) (*models.User, error)
}

type Handler struct {
	service service
}

func New(service service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request"})
		return
	}

	user, err := h.service.Register(req.Login, req.Password)
	if err != nil {
		render.Status(r, http.StatusConflict)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	render.JSON(w, r, user)
}
