package create

import "time"

// CreateUserDTO represents the data transfer object for the create user use case.
// It contains the validated data needed to create a new user.
type CreateUserDTO struct {
	Email                string `validate:"required,email"`
	Password             string `validate:"required,min=8"`
	PasswordConfirmation string `validate:"required,eqfield=Password"`
}

// CreateUserResultDTO represents the result of creating a user
type CreateUserResultDTO struct {
	Email     string
	CreatedAt time.Time
}
