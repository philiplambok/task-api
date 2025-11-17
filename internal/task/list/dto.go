package list

import "time"

// ListTasksDTO represents the query parameters for listing tasks
type ListTasksDTO struct {
	IsDone *bool // nil means no filter
}

// TaskDTO represents a task in the list response
type TaskDTO struct {
	UUID        string
	Title       string
	Description *string
	Deadline    *time.Time
	IsDone      bool
	CreatedAt   time.Time
}

// ListTasksResultDTO represents the result of listing tasks
type ListTasksResultDTO struct {
	Tasks []TaskDTO
}
