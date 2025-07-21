package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID       string
	Login    string
	Password string
}

type Password struct {
	ID       int       `json:"id"`
	UserID   uuid.UUID `json:"-"`
	Resource string    `json:"resource"`
	Login    string    `json:"login"`
	Password string    `json:"password"`
}

type Card struct {
	ID         int       `json:"id"`
	UserID     uuid.UUID `json:"-"`
	Bank       string    `json:"bank"`
	PAN        string    `json:"pan"`
	ValidThru  string    `json:"validThru"`
	Cardholder string    `json:"cardholder"`
}

type PlainText struct {
	ID      int       `json:"id"`
	UserID  uuid.UUID `json:"-"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
}
