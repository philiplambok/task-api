package repository

import (
	"context"

	"github.com/philiplambok/task-api/internal/auth/common/datamodel"
	commondatamodel "github.com/philiplambok/task-api/internal/common/datamodel"
	"gorm.io/gorm"
)

// FindUserByEmail retrieves a user by email from the database
func (r *Repository) FindUserByEmail(ctx context.Context, email string) (*datamodel.User, error) {
	var user datamodel.User

	query := `
		SELECT id, email, password_digest, created_at, updated_at
		FROM users
		WHERE email = $1 AND status = $2
	`

	err := r.db.WithContext(ctx).Raw(query, email, commondatamodel.UserActive).Scan(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	// Check if user was found
	if user.ID == 0 {
		return nil, nil
	}

	return &user, nil
}
