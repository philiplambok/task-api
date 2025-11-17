package registration

import "time"

// CreateUser represents the data needed to create a new user
type CreateUser struct {
	Email          string
	PasswordDigest string
}

// CreateUserResult represents the result of creating a new user
type CreateUserResult struct {
	ID        int64
	CreatedAt time.Time
}
