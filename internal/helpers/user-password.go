package helpers

import (
	"errors"

	"github.com/Shodocan/UserService/internal/configs/engine"
	"golang.org/x/crypto/bcrypt"
)

func ComparePasswordToHash(hash, password string) *engine.Error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return engine.ErrInvalidPassword()
		}
		return engine.ErrInternalFailure()
	}

	return nil
}

func HashPassword(password string) (string, *engine.Error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", engine.ErrInternalFailure()
	}

	return string(hash), nil
}
