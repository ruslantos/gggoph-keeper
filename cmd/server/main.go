package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"goph-keeper/internal/api/handlers/get_credentials"
	loginHandler "goph-keeper/internal/api/handlers/login"
	"goph-keeper/internal/api/handlers/post_credentials"
	registerHandler "goph-keeper/internal/api/handlers/register"
	config2 "goph-keeper/internal/config"
	"goph-keeper/internal/initializers/logger"
	auth_middleware "goph-keeper/internal/middleware/auth"
	"goph-keeper/internal/services/auth"
	"goph-keeper/internal/services/credentials"
	"goph-keeper/internal/storage/repositories/cred_repo"
	"goph-keeper/internal/storage/repositories/user_repo"

	//cred "goph-keeper/internal/services/get_credentials"
	"goph-keeper/internal/storage"
)

func main() {
	var config config2.Config
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}
	// Установка переменных окружения
	config2.SetEnvFromConfig(config)

	config2.InitJWT()

	// Инициализация БД
	ctx := context.Background()
	db, err := sqlx.Open("pgx", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatalf("Failed to init storage: %v", err)
	}
	defer db.Close()

	psqlStorage := storage.New(db)
	if err := psqlStorage.InitSchema(ctx); err != nil {
		logger.GetLogger().Fatal("Failed to init schema", zap.Error(err))
	}

	// репозитории
	userRepo := user_repo.New(db)
	credRepo := cred_repo.New(db)

	// сервисы
	authService := auth.New(userRepo)
	credentialService := credentials.New(credRepo)

	// handlers
	RegisterHandler := registerHandler.New(authService)
	LoginHandler := loginHandler.New(authService)
	GetCredentialsHandler := get_credentials.New(credentialService)
	PostCredentialsHandler := post_credentials.New(credentialService)

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Public
	r.Group(func(r chi.Router) {
		r.Post("/register", RegisterHandler.Handle)
		r.Post("/login", LoginHandler.Handle)
	})

	// Защищенные маршруты (с JWT)
	r.Group(func(r chi.Router) {
		r.Use(auth_middleware.JWTVerifier(config2.TokenAuth))
		r.Use(auth_middleware.JWTAuthenticator)

		r.Route("/creds", func(r chi.Router) {
			r.Post("/", PostCredentialsHandler.Handle)
		})
		r.Route("/creds/{field}", func(r chi.Router) {
			r.Get("/", GetCredentialsHandler.Handle)
		})
	})

	logger.GetLogger().Info("Server started on :8080")
	http.ListenAndServe(":8080", r)
}
