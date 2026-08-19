package api

import (
	"backend-test/internal/interfaces"
	"backend-test/internal/models/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct{ Service interfaces.AuthService }

func NewAuthHandler(service interfaces.AuthService) *AuthHandler {
	return &AuthHandler{Service: service}
}
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "valid email and password are required", "errors": err.Error()})
		return
	}
	token, user, err := h.Service.Login(c.Request.Context(), req)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "signed in successfully", "data": gin.H{"token": token, "user": user}})
}
func (h *AuthHandler) Me(c *gin.Context) { respond(c, http.StatusOK, "success", actor(c)) }
func (h *AuthHandler) Logout(c *gin.Context) {
	respond(c, http.StatusOK, "signed out successfully", nil)
}
