package repository

import (
	"backend-test/internal/models"
	"context"
	"strings"

	"gorm.io/gorm"
)

type UserRepository struct{ Database *gorm.DB }

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.Database.WithContext(ctx).Where("LOWER(email) = ?", strings.ToLower(email)).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	if err := r.Database.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) List(ctx context.Context, search, role, status string, page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64
	query := r.Database.WithContext(ctx).Model(&models.User{})
	if search != "" {
		term := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ?", term, term)
	}
	if role != "" {
		query = query.Where("role = ?", role)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("created_at ASC").Offset((page - 1) * limit).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return r.Database.WithContext(ctx).Create(user).Error
}
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	return r.Database.WithContext(ctx).Save(user).Error
}
