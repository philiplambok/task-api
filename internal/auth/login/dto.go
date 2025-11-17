package login

import (
	"errors"
	"time"

	"github.com/philiplambok/task-api/internal/auth/common/datamodel"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

// LoginDTO represents the data transfer object for the login use case.
type LoginDTO struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

// LoginResultDTO represents the data sent after a successful login.
type LoginResultDTO struct {
	User      *datamodel.User
	Token     string
	ExpiresAt time.Time
}
