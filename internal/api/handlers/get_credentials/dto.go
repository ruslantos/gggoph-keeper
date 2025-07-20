package get_credentials

import (
	"goph-keeper/internal/models"
)

type GetCredsResponse struct {
	Passwords  []models.Password  `json:"passwords,omitempty"`
	Cards      []models.Card      `json:"cards,omitempty"`
	PlainTexts []models.PlainText `json:"plainTexts,omitempty"`
}
