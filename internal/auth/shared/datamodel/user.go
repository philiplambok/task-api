package datamodel

import "time"

// User represents a user retrieved from the database for authentication
type User struct {
	ID             int64
	Email          string
	PasswordDigest string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
