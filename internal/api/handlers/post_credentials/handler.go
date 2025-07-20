package post_credentials

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"

	middleware "goph-keeper/internal/middleware/auth"
)

type CredsService interface {
	SavePassword(userID uuid.UUID, resource, login, password string) error
	SaveCard(userID uuid.UUID, bank, pan, validThru, cardholder string) error
	SavePlainText(userID uuid.UUID, title, content string) error
}

type CredsHandler struct {
	service CredsService
}

func New(service CredsService) *CredsHandler {
	return &CredsHandler{service: service}
}

func (h *CredsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	var req CredsRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request"})
		return
	}

	if req.Passwords != nil {
		err := h.service.SavePassword(
			userID,
			req.Passwords.Resource,
			req.Passwords.Login,
			req.Passwords.Password,
		)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to save password"})
			return
		}
	}

	if req.Cards != nil {
		err := h.service.SaveCard(
			userID,
			req.Cards.Bank,
			req.Cards.PAN,
			req.Cards.ValidThru,
			req.Cards.Cardholder,
		)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to save card"})
			return
		}
	}

	if req.PlainText != nil {
		err := h.service.SavePlainText(
			userID,
			req.PlainText.Title,
			req.PlainText.Content,
		)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to save plain text"})
			return
		}
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"status": "success"})
}
