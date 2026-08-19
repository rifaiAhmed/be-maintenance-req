package interfaces

import (
	"backend-test/internal/models"
	"backend-test/internal/models/dto"
	"context"
)

type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (string, *models.User, error)
	CurrentUser(ctx context.Context, id uint) (*models.User, error)
}

type UserService interface {
	List(ctx context.Context, actor *models.User, search, role, status string, page, limit int) ([]models.User, models.Pagination, error)
	Get(ctx context.Context, actor *models.User, id uint) (*models.User, error)
	Create(ctx context.Context, actor *models.User, req dto.CreateUserRequest) (*models.User, error)
	Update(ctx context.Context, actor *models.User, id uint, req dto.UpdateUserRequest) (*models.User, error)
}

type RequestService interface {
	List(ctx context.Context, actor *models.User, filter dto.RequestFilter) ([]models.MaintenanceRequest, models.Pagination, error)
	Get(ctx context.Context, actor *models.User, id string) (*models.MaintenanceRequest, error)
	Create(ctx context.Context, actor *models.User, req dto.CreateMaintenanceRequest) (*models.MaintenanceRequest, error)
	Update(ctx context.Context, actor *models.User, id string, req dto.UpdateMaintenanceRequest) (*models.MaintenanceRequest, error)
	Review(ctx context.Context, actor *models.User, id string, req dto.ReviewMaintenanceRequest) (*models.MaintenanceRequest, error)
	Delete(ctx context.Context, actor *models.User, id string) error
	Dashboard(ctx context.Context, actor *models.User) (map[string]int64, []models.MaintenanceRequest, error)
}
