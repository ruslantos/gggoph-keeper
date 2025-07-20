package get_credentials

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"

	middleware "goph-keeper/internal/middleware/auth"
	"goph-keeper/internal/models"
)

type CredsService interface {
	GetPasswords(userID uuid.UUID) ([]models.Password, error)
	GetCards(userID uuid.UUID) ([]models.Card, error)
	GetPlainTexts(userID uuid.UUID) ([]models.PlainText, error)
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

	field := r.URL.Query().Get("field")

	var response GetCredsResponse
	var err error

	switch field {
	case "passwords":
		response.Passwords, err = h.service.GetPasswords(userID)
	case "cards":
		response.Cards, err = h.service.GetCards(userID)
	case "plainTexts":
		response.PlainTexts, err = h.service.GetPlainTexts(userID)
	default:
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid field parameter"})
		return
	}

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to get data"})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}
