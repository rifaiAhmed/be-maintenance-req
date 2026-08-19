package service

import (
	"backend-test/helpers"
	"backend-test/internal/interfaces"
	"backend-test/internal/models"
	"backend-test/internal/models/dto"
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct{ Users interfaces.UserRepository }

func NewAuthService(users interfaces.UserRepository) *AuthService { return &AuthService{Users: users} }

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (string, *models.User, error) {
	user, err := s.Users.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, helpers.ErrUnauthorized
		}
		return "", nil, err
	}
	if user.Status != models.UserActive || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		return "", nil, helpers.ErrUnauthorized
	}
	token, err := helpers.GenerateToken(ctx, user, time.Now())
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

func (s *AuthService) CurrentUser(ctx context.Context, id uint) (*models.User, error) {
	user, err := s.Users.FindByID(ctx, id)
	if err != nil || user.Status != models.UserActive {
		return nil, helpers.ErrUnauthorized
	}
	return user, nil
}
