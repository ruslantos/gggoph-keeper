package storage

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"path/filepath"
)

const keyFileName = "encryption.key"

// KeyStorage управляет ключом шифрования
type KeyStorage struct {
	keyPath string
}

func NewKeyStorage() *KeyStorage {
	configDir, _ := os.UserConfigDir()
	return &KeyStorage{
		keyPath: filepath.Join(configDir, "gophkeeper", keyFileName),
	}
}

// GenerateAndSave создает новый ключ и сохраняет его
func (s *KeyStorage) GenerateAndSave() ([]byte, error) {
	key := make([]byte, 32) // AES-256
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(filepath.Dir(s.keyPath), 0700); err != nil {
		return nil, err
	}

	encoded := base64.StdEncoding.EncodeToString(key)
	return key, os.WriteFile(s.keyPath, []byte(encoded), 0600)
}

// Load загружает ключ из хранилища
func (s *KeyStorage) Load() ([]byte, error) {
	data, err := os.ReadFile(s.keyPath)
	if err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(string(data))
}
