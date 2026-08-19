package dto

import "backend-test/internal/models"

type CreateUserRequest struct {
	Name     string            `json:"name" binding:"required,min=2,max=120"`
	Email    string            `json:"email" binding:"required,email,max=180"`
	Password string            `json:"password" binding:"required,min=8,max=72"`
	Role     models.UserRole   `json:"role" binding:"required,oneof=Admin Supervisor Operator"`
	Status   models.UserStatus `json:"status" binding:"omitempty,oneof=Active Inactive"`
}

type UpdateUserRequest struct {
	Name   string            `json:"name" binding:"required,min=2,max=120"`
	Email  string            `json:"email" binding:"required,email,max=180"`
	Role   models.UserRole   `json:"role" binding:"required,oneof=Admin Supervisor Operator"`
	Status models.UserStatus `json:"status" binding:"required,oneof=Active Inactive"`
}
