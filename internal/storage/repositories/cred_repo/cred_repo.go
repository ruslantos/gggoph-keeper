package cred_repo

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"goph-keeper/internal/models"
)

type PostgresCredsRepository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *PostgresCredsRepository {
	return &PostgresCredsRepository{db: db}
}

func (r *PostgresCredsRepository) SavePassword(password models.Password) error {
	_, err := r.db.Exec(
		"INSERT INTO passwords (user_id, resource, login, password) VALUES ($1, $2, $3, $4)",
		password.UserID, password.Resource, password.Login, password.Password,
	)
	return err
}

func (r *PostgresCredsRepository) SaveCard(card models.Card) error {
	_, err := r.db.Exec(
		"INSERT INTO cards (user_id, bank, pan, valid_thru, cardholder) VALUES ($1, $2, $3, $4, $5)",
		card.UserID, card.Bank, card.PAN, card.ValidThru, card.Cardholder,
	)
	return err
}

func (r *PostgresCredsRepository) SavePlainText(text models.PlainText) error {
	_, err := r.db.Exec(
		"INSERT INTO plain_texts (user_id, title, content) VALUES ($1, $2, $3)",
		text.UserID, text.Title, text.Content,
	)
	return err
}

func (r *PostgresCredsRepository) GetPasswords(userID uuid.UUID) ([]models.Password, error) {
	rows, err := r.db.Query("SELECT id, resource, login, password FROM passwords WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var passwords []models.Password
	for rows.Next() {
		var p models.Password
		if err := rows.Scan(&p.ID, &p.Resource, &p.Login, &p.Password); err != nil {
			return nil, err
		}
		passwords = append(passwords, p)
	}
	return passwords, nil
}

func (r *PostgresCredsRepository) GetCards(userID uuid.UUID) ([]models.Card, error) {
	rows, err := r.db.Query("SELECT id, bank, pan, valid_thru, cardholder FROM cards WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []models.Card
	for rows.Next() {
		var c models.Card
		if err := rows.Scan(&c.ID, &c.Bank, &c.PAN, &c.ValidThru, &c.Cardholder); err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}
	return cards, nil
}

func (r *PostgresCredsRepository) GetPlainTexts(userID uuid.UUID) ([]models.PlainText, error) {
	rows, err := r.db.Query("SELECT id, title, content FROM plain_texts WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var texts []models.PlainText
	for rows.Next() {
		var t models.PlainText
		if err := rows.Scan(&t.ID, &t.Title, &t.Content); err != nil {
			return nil, err
		}
		texts = append(texts, t)
	}
	return texts, nil
}
