package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Failled to hash password: %v", err)
	}

	return string(hashedpassword), nil
}

func CheckPassword(password string, hashpassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashpassword), []byte(password))
}
