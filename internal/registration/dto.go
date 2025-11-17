package registration

import "time"

// RegisterDTO represents the data transfer object for the registration use case.
// It contains the validated data needed to create a new user.
type RegisterDTO struct {
	Email                string `validate:"required,email"`
	Password             string `validate:"required,min=8"`
	PasswordConfirmation string `validate:"required,eqfield=Password"`
}

// RegisterResultDTO represents the result of registering a user
type RegisterResultDTO struct {
	Email     string
	CreatedAt time.Time
}
