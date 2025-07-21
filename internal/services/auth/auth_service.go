package auth

import (
	"context"
	"errors"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/jwa"
	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"

	"goph-keeper/internal/config"
	"goph-keeper/internal/models"
)

type UserRepository interface {
	Save(ctx context.Context, user *models.User) error
	FindByLogin(ctx context.Context, login string) (*models.User, error)
	FindByID(ctx context.Context, id string) (*models.User, error)
}

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New(
		jwa.HS256.String(),                     // Алгоритм подписи
		[]byte("your-32-byte-secret-key-here"), // Секретный ключ
		nil,                                    // Дополнительные опции
	)
}

type Service struct {
	userRepo UserRepository
}

func New(repo UserRepository) *Service {
	return &Service{userRepo: repo}
}

func (s *Service) Register(login, password string) (*models.User, error) {
	// Проверка существования пользователя
	existing, err := s.userRepo.FindByLogin(context.Background(), login)
	if existing != nil {
		return nil, errors.New("user already exists")
	}

	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	id, _ := uuid.NewUUID()
	user := &models.User{
		ID:       id.String(),
		Login:    login,
		Password: string(hashedPassword),
	}

	// Сохранение в репозитории
	if err := s.userRepo.Save(context.Background(), user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Authenticate(login, password string) (*models.User, error) {
	user, err := s.userRepo.FindByLogin(context.Background(), login)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	// Сравнение паролей
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

func (s *Service) GenerateToken(userID string) (string, error) {
	_, token, err := config.TokenAuth.Encode(map[string]interface{}{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	return token, err
}
