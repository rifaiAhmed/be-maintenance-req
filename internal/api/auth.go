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

// Login godoc
// @Summary Login user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} map[string]interface{} "success"
// @Failure 400 {object} map[string]interface{} "invalid request"
// @Router /auth/login [post]
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

// Me godoc
// @Summary Get current user
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]interface{} "success"
// @Security BearerAuth
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) { respond(c, http.StatusOK, "success", actor(c)) }

// Logout godoc
// @Summary Logout current user
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]interface{} "success"
// @Security BearerAuth
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	respond(c, http.StatusOK, "signed out successfully", nil)
}
