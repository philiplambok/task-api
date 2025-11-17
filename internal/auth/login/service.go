package login

import (
	"context"

	"github.com/philiplambok/task-api/internal/auth/common/datamodel"
	"github.com/philiplambok/task-api/internal/auth/common/domain"
)

// Querier defines the interface for querying user data
type Querier interface {
	FindUserByEmail(ctx context.Context, email string) (*datamodel.User, error)
}

// Service handles the login business logic
type Service struct {
	querier       Querier
	jwtSecret     string
	tokenExpHours int
}

// NewService creates a new login service
func NewService(querier Querier, jwtSecret string, tokenExpHours int) *Service {
	return &Service{
		querier:       querier,
		jwtSecret:     jwtSecret,
		tokenExpHours: tokenExpHours,
	}
}

// CreateUserAuthToken performs the login operation and creates an authentication token
func (s *Service) CreateUserAuthToken(ctx context.Context, dto *LoginDTO) (*LoginResultDTO, error) {
	// Find user by email
	user, err := s.querier.FindUserByEmail(ctx, dto.Email)
	if err != nil {
		return nil, err
	}

	// Check if user exists
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Compare password
	err = domain.ComparePassword(user.PasswordDigest, dto.Password)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate token
	token, expiresAt, err := domain.GenerateToken(user.ID, user.Email, s.jwtSecret, s.tokenExpHours)
	if err != nil {
		return nil, err
	}

	return &LoginResultDTO{
		User:      user,
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}
