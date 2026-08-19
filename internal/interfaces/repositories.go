package interfaces

import (
	"backend-test/internal/models"
	"backend-test/internal/models/dto"
	"context"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id uint) (*models.User, error)
	List(ctx context.Context, search, role, status string, page, limit int) ([]models.User, int64, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
}

type RequestRepository interface {
	List(ctx context.Context, filter dto.RequestFilter, viewer *models.User) ([]models.MaintenanceRequest, int64, error)
	FindByID(ctx context.Context, id string) (*models.MaintenanceRequest, error)
	Create(ctx context.Context, request *models.MaintenanceRequest, activity *models.RequestActivity) error
	Update(ctx context.Context, request *models.MaintenanceRequest, activity *models.RequestActivity) error
	Review(ctx context.Context, request *models.MaintenanceRequest, activity *models.RequestActivity) error
	Delete(ctx context.Context, id string) error
	DashboardCounts(ctx context.Context, viewer *models.User) (map[string]int64, error)
	DB() *gorm.DB
}
