package repository

import (
	"context"

	"github.com/philiplambok/task-api/internal/user/common/datamodel"
	"gorm.io/gorm"
)

// FindUserByEmail finds a user by email address
func (r *Repository) FindUserByEmail(ctx context.Context, email string) (*datamodel.User, error) {
	var user datamodel.User

	err := r.db.WithContext(ctx).
		Table("users").
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
