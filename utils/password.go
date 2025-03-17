package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HasghPassword return the becrypt hash of the password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash the password")
	}
	return string(hashedPassword), nil
}

// CheckPassword checks if provied passwrod maches with hash
func CheckPassword(password string, passwordHash string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
}
