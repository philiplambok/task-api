package registration

import (
	"context"
	"time"
)

// Querier defines the interface for querying user data
type Querier interface {
	CreateUser(ctx context.Context, params CreateUser) (time.Time, error)
}

// Service handles the registration business logic
type Service struct {
	querier Querier
}

// NewService creates a new registration service
func NewService(querier Querier) *Service {
	return &Service{
		querier: querier,
	}
}

// Register performs the user registration operation
func (s *Service) Register(ctx context.Context, dto *RegisterDTO) (*RegisterResultDTO, error) {
	// Hash the password using domain service
	hashedPassword, err := HashPassword(dto.Password)
	if err != nil {
		return nil, err
	}

	// Convert DTO to datamodel
	params := CreateUser{
		Email:          dto.Email,
		PasswordDigest: hashedPassword,
	}

	createdAt, err := s.querier.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return &RegisterResultDTO{
		Email:     dto.Email,
		CreatedAt: createdAt,
	}, nil
}
