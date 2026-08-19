package models

import "time"

type UserRole string

type UserStatus string

const (
	RoleOperator   UserRole = "Operator"
	RoleSupervisor UserRole = "Supervisor"
	RoleAdmin      UserRole = "Admin"

	UserActive   UserStatus = "Active"
	UserInactive UserStatus = "Inactive"
)

type User struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Name         string     `gorm:"size:120;not null" json:"name"`
	Email        string     `gorm:"size:180;uniqueIndex;not null" json:"email"`
	PasswordHash string     `gorm:"size:255;not null" json:"-"`
	Role         UserRole   `gorm:"type:varchar(20);not null" json:"role"`
	Status       UserStatus `gorm:"type:varchar(20);not null;default:Active" json:"status"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}
