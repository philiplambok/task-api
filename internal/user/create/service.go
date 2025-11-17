package create

import (
	"context"
	"time"

	"github.com/philiplambok/task-api/internal/user/common/datamodel"
	"github.com/philiplambok/task-api/internal/user/common/domain"
)

// Querier defines the interface for querying user data
type Querier interface {
	CreateUser(ctx context.Context, params datamodel.CreateUser) (time.Time, error)
}

// Service handles the create user business logic
type Service struct {
	querier Querier
}

// NewService creates a new create user service
func NewService(querier Querier) *Service {
	return &Service{
		querier: querier,
	}
}

// CreateUser performs the user creation operation
func (s *Service) CreateUser(ctx context.Context, dto *CreateUserDTO) (*CreateUserResultDTO, error) {
	// Hash the password using domain service
	hashedPassword, err := domain.HashPassword(dto.Password)
	if err != nil {
		return nil, err
	}

	// Convert DTO to datamodel
	params := datamodel.CreateUser{
		Email:          dto.Email,
		PasswordDigest: hashedPassword,
	}

	createdAt, err := s.querier.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return &CreateUserResultDTO{
		Email:     dto.Email,
		CreatedAt: createdAt,
	}, nil
}
