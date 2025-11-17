package create

import (
	"context"
	"time"

	"github.com/philiplambok/task-api/internal/task/common/datamodel"
)

// Querier defines the interface for querying task data
type Querier interface {
	CreateTask(ctx context.Context, params datamodel.CreateTask) (*datamodel.CreateTaskResult, error)
}

// Service handles the create task business logic
type Service struct {
	querier Querier
}

// NewService creates a new create task service
func NewService(querier Querier) *Service {
	return &Service{
		querier: querier,
	}
}

// CreateTask performs the task creation operation
func (s *Service) CreateTask(ctx context.Context, userID int64, dto *CreateTaskDTO) (*CreateTaskResultDTO, error) {
	// Parse deadline if provided
	var deadline *time.Time
	if dto.Deadline != nil {
		parsedTime, err := time.Parse("2006-01-02", *dto.Deadline)
		if err != nil {
			return nil, err
		}
		deadline = &parsedTime
	}

	// Convert DTO to datamodel
	params := datamodel.CreateTask{
		UserID:      userID,
		Title:       dto.Title,
		Description: dto.Description,
		Deadline:    deadline,
	}

	result, err := s.querier.CreateTask(ctx, params)
	if err != nil {
		return nil, err
	}

	return &CreateTaskResultDTO{
		ID:          result.ID,
		Title:       dto.Title,
		Description: dto.Description,
		Deadline:    deadline,
		IsDone:      false,
		CreatedAt:   result.CreatedAt,
	}, nil
}
