package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email,max=180"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}
