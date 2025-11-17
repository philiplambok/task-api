package list

import (
	"context"

	"github.com/philiplambok/task-api/internal/task/common/datamodel"
)

// Querier defines the interface for querying task data
type Querier interface {
	ListTasks(ctx context.Context, params datamodel.ListTasksParams) ([]datamodel.Task, error)
}

// Service handles the list tasks business logic
type Service struct {
	querier Querier
}

// NewService creates a new list tasks service
func NewService(querier Querier) *Service {
	return &Service{
		querier: querier,
	}
}

// ListTasks retrieves all tasks for a user with optional filtering
func (s *Service) ListTasks(ctx context.Context, userID int64, dto *ListTasksDTO) (*ListTasksResultDTO, error) {
	// Build query parameters
	params := datamodel.ListTasksParams{
		UserID: userID,
		IsDone: dto.IsDone,
	}

	tasks, err := s.querier.ListTasks(ctx, params)
	if err != nil {
		return nil, err
	}

	// Convert datamodel tasks to DTOs
	taskDTOs := make([]TaskDTO, len(tasks))
	for i, task := range tasks {
		taskDTOs[i] = TaskDTO{
			UUID:        task.UUID,
			Title:       task.Title,
			Description: task.Description,
			Deadline:    task.Deadline,
			IsDone:      task.IsDone,
			CreatedAt:   task.CreatedAt,
		}
	}

	return &ListTasksResultDTO{
		Tasks: taskDTOs,
	}, nil
}
