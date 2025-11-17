package datamodel

import "time"

// CreateTask represents the data needed to create a new task
type CreateTask struct {
	UserID      int64
	Title       string
	Description *string
	Deadline    *time.Time
}

// CreateTaskResult represents the result of creating a new task
type CreateTaskResult struct {
	ID        int64
	UUID      string
	CreatedAt time.Time
}

// Task represents a task retrieved from the database
type Task struct {
	ID          int64
	UUID        string
	UserID      int64
	Title       string
	Description *string
	Deadline    *time.Time
	IsDone      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ListTasksParams represents the parameters for listing tasks
type ListTasksParams struct {
	UserID int64
	IsDone *bool // nil means no filter, otherwise filter by the value
}
