package repository

import (
	"context"

	"github.com/philiplambok/task-api/internal/task/common/datamodel"
)

// ListTasks retrieves all tasks for a user with optional filtering by is_done status
func (r *Repository) ListTasks(ctx context.Context, params datamodel.ListTasksParams) ([]datamodel.Task, error) {
	var tasks []datamodel.Task

	query := `
		SELECT id, uuid, user_id, title, description, deadline, is_done, created_at, updated_at
		FROM tasks
		WHERE user_id = $1
	`

	args := []interface{}{params.UserID}

	// Add is_done filter if provided
	if params.IsDone != nil {
		query += " AND is_done = $2"
		args = append(args, *params.IsDone)
	}

	// Order by created_at descending (newest first)
	query += " ORDER BY created_at DESC"

	err := r.db.WithContext(ctx).Raw(query, args...).Scan(&tasks).Error
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
