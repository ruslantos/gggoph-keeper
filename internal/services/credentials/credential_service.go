package credentials

import (
	"github.com/google/uuid"

	"goph-keeper/internal/models"
)

type CredsRepository interface {
	SavePassword(password models.Password) error
	SaveCard(card models.Card) error
	SavePlainText(text models.PlainText) error
	GetPasswords(userID uuid.UUID) ([]models.Password, error)
	GetCards(userID uuid.UUID) ([]models.Card, error)
	GetPlainTexts(userID uuid.UUID) ([]models.PlainText, error)
}

type CredsService struct {
	repo CredsRepository
}

func New(repo CredsRepository) *CredsService {
	return &CredsService{repo: repo}
}

func (s *CredsService) SavePassword(userID uuid.UUID, resource, login, password string) error {
	return s.repo.SavePassword(models.Password{
		UserID:   userID,
		Resource: resource,
		Login:    login,
		Password: password,
	})
}

func (s *CredsService) SaveCard(userID uuid.UUID, bank, pan, validThru, cardholder string) error {
	return s.repo.SaveCard(models.Card{
		UserID:     userID,
		Bank:       bank,
		PAN:        pan,
		ValidThru:  validThru,
		Cardholder: cardholder,
	})
}

func (s *CredsService) SavePlainText(userID uuid.UUID, title, content string) error {
	return s.repo.SavePlainText(models.PlainText{
		UserID:  userID,
		Title:   title,
		Content: content,
	})
}

func (s *CredsService) GetPasswords(userID uuid.UUID) ([]models.Password, error) {
	return s.repo.GetPasswords(userID)
}

func (s *CredsService) GetCards(userID uuid.UUID) ([]models.Card, error) {
	return s.repo.GetCards(userID)
}

func (s *CredsService) GetPlainTexts(userID uuid.UUID) ([]models.PlainText, error) {
	return s.repo.GetPlainTexts(userID)
}
