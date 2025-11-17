package repository

import (
	"context"

	commondatamodel "github.com/philiplambok/task-api/internal/common/datamodel"
	"github.com/philiplambok/task-api/internal/task/common/datamodel"
)

// CreateTask creates a new task in the database and returns the task ID and created_at timestamp
// The task will be assigned to the user's default list
func (r *Repository) CreateTask(ctx context.Context, params datamodel.CreateTask) (*datamodel.CreateTaskResult, error) {
	var result datamodel.CreateTaskResult

	// First, get the user's default list ID
	var listID int64
	listQuery := `
		SELECT id FROM lists
		WHERE user_id = $1 AND name = $2
		LIMIT 1
	`
	err := r.db.WithContext(ctx).Raw(listQuery, params.UserID, commondatamodel.DefaultListName).Scan(&listID).Error
	if err != nil {
		return nil, err
	}

	// Insert the task and return the id, uuid, and created_at timestamp
	taskQuery := `
		INSERT INTO tasks (user_id, list_id, title, description, deadline, is_done, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, uuid, created_at
	`

	err = r.db.WithContext(ctx).Raw(
		taskQuery,
		params.UserID,
		listID,
		params.Title,
		params.Description,
		params.Deadline,
		false, // is_done defaults to false
	).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return &result, nil
}
