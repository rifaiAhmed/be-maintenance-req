package api

import (
	"backend-test/internal/interfaces"
	"backend-test/internal/models/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct{ Service interfaces.UserService }

func NewUserHandler(service interfaces.UserService) *UserHandler {
	return &UserHandler{Service: service}
}
func userID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user id"})
		return 0, false
	}
	return uint(id), true
}
func (h *UserHandler) List(c *gin.Context) {
	page, limit := pagination(c)
	users, meta, err := h.Service.List(c.Request.Context(), actor(c), c.Query("search"), c.Query("role"), c.Query("status"), page, limit)
	if err != nil {
		respondError(c, err)
		return
	}
	respondList(c, users, meta)
}
func (h *UserHandler) Get(c *gin.Context) {
	id, ok := userID(c)
	if !ok {
		return
	}
	user, err := h.Service.Get(c.Request.Context(), actor(c), id)
	if err != nil {
		respondError(c, err)
		return
	}
	respond(c, http.StatusOK, "success", user)
}
func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user data", "errors": err.Error()})
		return
	}
	user, err := h.Service.Create(c.Request.Context(), actor(c), req)
	if err != nil {
		respondError(c, err)
		return
	}
	respond(c, http.StatusCreated, "user created successfully", user)
}
func (h *UserHandler) Update(c *gin.Context) {
	id, ok := userID(c)
	if !ok {
		return
	}
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user data", "errors": err.Error()})
		return
	}
	user, err := h.Service.Update(c.Request.Context(), actor(c), id, req)
	if err != nil {
		respondError(c, err)
		return
	}
	respond(c, http.StatusOK, "user updated successfully", user)
}
