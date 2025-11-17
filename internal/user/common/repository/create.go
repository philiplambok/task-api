package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	commondatamodel "github.com/philiplambok/task-api/internal/common/datamodel"
	"github.com/philiplambok/task-api/internal/user/common/datamodel"
	"github.com/philiplambok/task-api/internal/user/common/domain"
)

// CreateUser creates a new user in the database and returns the created_at timestamp
func (r *Repository) CreateUser(ctx context.Context, params datamodel.CreateUser) (time.Time, error) {
	var result datamodel.CreateUserResult

	// Use raw SQL to insert and return the user id and created_at timestamp
	userQuery := `
		INSERT INTO users (email, password_digest, status, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, created_at
	`

	err := r.db.WithContext(ctx).Raw(userQuery, params.Email, params.PasswordDigest, commondatamodel.UserActive).Scan(&result).Error
	if err != nil {
		// Check if it's a unique constraint violation on email
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && isEmailDuplicateError(pgErr) {
			return time.Time{}, domain.ErrDuplicateEmail
		}
		return time.Time{}, err
	}

	// Create default list for the user
	listQuery := `
		INSERT INTO lists (user_id, name, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
	`

	err = r.db.WithContext(ctx).Exec(listQuery, result.ID, commondatamodel.DefaultListName).Error
	if err != nil {
		return time.Time{}, err
	}

	return result.CreatedAt, nil
}

// isEmailDuplicateError checks if the PostgreSQL error is a unique constraint violation on the email column
func isEmailDuplicateError(pgErr *pgconn.PgError) bool {
	// PostgreSQL unique violation error code is 23505
	// Check both the error code and the constraint name to identify which column failed
	return pgErr.Code == "23505" && pgErr.ConstraintName == "users_email_key"
}
