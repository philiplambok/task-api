package create

import "time"

// CreateTaskDTO represents the data transfer object for the create task use case.
// It contains the validated data needed to create a new task.
type CreateTaskDTO struct {
	Title       string  `validate:"required,max=255"`
	Description *string `validate:"omitempty,max=1000"`
	Deadline    *string `validate:"omitempty,datetime=2006-01-02"`
}

// CreateTaskResultDTO represents the result of creating a task
type CreateTaskResultDTO struct {
	UUID        string
	Title       string
	Description *string
	Deadline    *time.Time
	IsDone      bool
	CreatedAt   time.Time
}
