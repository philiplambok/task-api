package datamodel

import "time"

// CreateUser represents the data needed to create a new user
type CreateUser struct {
	Email          string
	PasswordDigest string
}

// User represents a user retrieved from the database
type User struct {
	ID             int64
	Email          string
	PasswordDigest string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
