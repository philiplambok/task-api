package datamodel

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

// User represents a user retrieved from the database
type User struct {
	ID             int64
	Email          string
	PasswordDigest string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
