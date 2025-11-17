package profile

import (
	"context"

	"github.com/philiplambok/task-api/internal/user/shared/datamodel"
)

//go:generate go run go.uber.org/mock/mockgen@latest -destination=mock_querier_test.go -package=profile -source=service.go

// Querier defines the interface for querying user data
type Querier interface {
	FindUserByEmail(ctx context.Context, email string) (*datamodel.User, error)
}

// Service handles the get profile business logic
type Service struct {
	querier Querier
}

// NewService creates a new profile service
func NewService(querier Querier) *Service {
	return &Service{
		querier: querier,
	}
}

// GetProfile retrieves the user profile by email
func (s *Service) GetProfile(ctx context.Context, email string) (*ProfileResultDTO, error) {
	// Find user by email
	user, err := s.querier.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// Check if user exists
	if user == nil {
		return nil, ErrUserNotFound
	}

	return &ProfileResultDTO{
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}
