package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/render"

	"goph-keeper/internal/config"
	"goph-keeper/internal/models"
)

type service interface {
	Authenticate(login, password string) (*models.User, error)
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

	// Аутентификация через сервис
	user, err := h.service.Authenticate(req.Login, req.Password)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "Invalid get_credentials"})
		return
	}

	// Генерация JWT
	_, tokenString, _ := config.TokenAuth.Encode(map[string]interface{}{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	render.JSON(w, r, map[string]string{
		"token":         tokenString,
		"refresh_token": "generated-refresh-token", // Реализуйте отдельно
	})
}
