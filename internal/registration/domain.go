package registration

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// Domain errors
var ErrDuplicateEmail = errors.New("email already exists")

// HashPassword hashes a plain text password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
