package helper

import (
	"errors"
	"portal_api/config"

	"github.com/google/uuid"
)

func GenerateUniqueId() string {
	return uuid.NewString()
}

func GetSecretKey() (string, error) {
	secretKey := config.Cfg.SecretKey
	if secretKey == "" {
		return "", errors.New("secret key is empty")
	}

	return secretKey, nil
}
