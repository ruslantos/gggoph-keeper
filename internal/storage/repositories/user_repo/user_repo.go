package user_repo

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"goph-keeper/internal/models"
)

type UserRepo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Save(ctx context.Context, user *models.User) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (id, login, password_hash) 
		VALUES ($1, $2, $3) 
		ON CONFLICT (id) DO UPDATE 
		SET login = $2, password_hash = $3`,
		user.ID, user.Login, user.Password)

	return err
}

func (r *UserRepo) FindByLogin(ctx context.Context, login string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRowxContext(ctx,
		`SELECT id, login, password_hash FROM users WHERE login = $1`,
		login).Scan(&user.ID, &user.Login, &user.Password)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &user, err
}

func (r *UserRepo) FindByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRowxContext(ctx,
		`SELECT id, login, password_hash FROM users WHERE id = $1`,
		id).Scan(&user.ID, &user.Login, &user.Password)

	return &user, err
}
