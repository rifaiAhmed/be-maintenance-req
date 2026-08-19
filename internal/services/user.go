package service

import (
	"backend-test/helpers"
	"backend-test/internal/interfaces"
	"backend-test/internal/models"
	"backend-test/internal/models/dto"
	"context"
	"errors"
	"math"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct{ Users interfaces.UserRepository }

func NewUserService(users interfaces.UserRepository) *UserService { return &UserService{Users: users} }
func requireAdmin(actor *models.User) error {
	if actor == nil || actor.Role != models.RoleAdmin {
		return helpers.ErrForbidden
	}
	return nil
}

func (s *UserService) List(ctx context.Context, actor *models.User, search, role, status string, page, limit int) ([]models.User, models.Pagination, error) {
	if err := requireAdmin(actor); err != nil {
		return nil, models.Pagination{}, err
	}
	users, total, err := s.Users.List(ctx, search, role, status, page, limit)
	return users, models.Pagination{Page: page, Limit: limit, Total: total, TotalPages: int(math.Ceil(float64(total) / float64(limit)))}, err
}
func (s *UserService) Get(ctx context.Context, actor *models.User, id uint) (*models.User, error) {
	if err := requireAdmin(actor); err != nil {
		return nil, err
	}
	user, err := s.Users.FindByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, helpers.ErrNotFound
	}
	return user, err
}
func (s *UserService) Create(ctx context.Context, actor *models.User, req dto.CreateUserRequest) (*models.User, error) {
	if err := requireAdmin(actor); err != nil {
		return nil, err
	}
	if _, err := s.Users.FindByEmail(ctx, req.Email); err == nil {
		return nil, helpers.ErrConflict
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	status := req.Status
	if status == "" {
		status = models.UserActive
	}
	user := &models.User{Name: req.Name, Email: req.Email, PasswordHash: string(hash), Role: req.Role, Status: status}
	if err := s.Users.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
func (s *UserService) Update(ctx context.Context, actor *models.User, id uint, req dto.UpdateUserRequest) (*models.User, error) {
	if err := requireAdmin(actor); err != nil {
		return nil, err
	}
	user, err := s.Users.FindByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, helpers.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if existing, findErr := s.Users.FindByEmail(ctx, req.Email); findErr == nil && existing.ID != id {
		return nil, helpers.ErrConflict
	} else if findErr != nil && !errors.Is(findErr, gorm.ErrRecordNotFound) {
		return nil, findErr
	}
	user.Name, user.Email, user.Role, user.Status = req.Name, req.Email, req.Role, req.Status
	if err := s.Users.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
