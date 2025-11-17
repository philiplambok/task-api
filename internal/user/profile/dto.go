package profile

import (
	"errors"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

// ProfileResultDTO represents the result of getting a user profile
type ProfileResultDTO struct {
	Email     string
	CreatedAt time.Time
}
