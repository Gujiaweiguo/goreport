package repository

import (
	"context"

	"github.com/gujiaweiguo/goreport/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
