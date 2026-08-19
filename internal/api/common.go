package api

import (
	"backend-test/helpers"
	"backend-test/internal/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func actor(c *gin.Context) *models.User {
	value, exists := c.Get("currentUser")
	if !exists {
		return nil
	}
	user, _ := value.(*models.User)
	return user
}
func pagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "25"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 25
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit
}
func respond(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, gin.H{"message": message, "data": data})
}
func respondList(c *gin.Context, data interface{}, meta models.Pagination) {
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": data, "meta": meta})
}
func respondError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	message := "internal server error"
	switch {
	case errors.Is(err, helpers.ErrUnauthorized):
		status, message = http.StatusUnauthorized, "invalid credentials or session"
	case errors.Is(err, helpers.ErrForbidden):
		status, message = http.StatusForbidden, err.Error()
	case errors.Is(err, helpers.ErrNotFound):
		status, message = http.StatusNotFound, err.Error()
	case errors.Is(err, helpers.ErrConflict):
		status, message = http.StatusConflict, err.Error()
	case errors.Is(err, helpers.ErrInvalid):
		status, message = http.StatusBadRequest, err.Error()
	}
	c.JSON(status, gin.H{"message": message})
}
