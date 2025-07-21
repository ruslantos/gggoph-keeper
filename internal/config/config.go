package config

import (
	"fmt"
	"os"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/jwa"
)

var TokenAuth *jwtauth.JWTAuth

// Config структура для хранения конфигурации из TOML
type Config struct {
	Database DatabaseConfig `toml:"database"`
	Server   ServerConfig   `toml:"server"`
	Other    OtherConfig    `toml:"other"`
}

type DatabaseConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Name     string `toml:"name"`
	URL      string `toml:"url"`
}

type ServerConfig struct {
	Port int    `toml:"port"`
	Env  string `toml:"env"`
}

type OtherConfig struct {
	Secret string `toml:"secret"`
}

func SetEnvFromConfig(config Config) {
	// Устанавливаем переменные окружения для базы данных
	os.Setenv("DB_HOST", config.Database.Host)
	os.Setenv("DB_PORT", fmt.Sprintf("%d", config.Database.Port))
	os.Setenv("DB_USER", config.Database.User)
	os.Setenv("DB_PASSWORD", config.Database.Password)
	os.Setenv("DB_NAME", config.Database.Name)
	os.Setenv("DATABASE_URL", config.Database.URL)

	// Устанавливаем переменные окружения для сервера
	os.Setenv("SERVER_PORT", fmt.Sprintf("%d", config.Server.Port))
	os.Setenv("SERVER_ENV", config.Server.Env)

	os.Setenv("JWT_SECRET", config.Other.Secret)
}

func InitJWT() {
	secret := os.Getenv("JWT_SECRET")
	if len(secret) < 32 {
		panic("JWT_SECRET must be at least 32 chars")
	}

	TokenAuth = jwtauth.New(
		jwa.HS256.String(),
		[]byte(secret),
		nil,
	)
}
