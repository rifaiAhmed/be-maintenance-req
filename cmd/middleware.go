package cmd

import (
	"backend-test/helpers"
	"backend-test/internal/interfaces"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func MiddlewareCORS() gin.HandlerFunc {
	origins := strings.Split(helpers.GetEnv("CORS_ORIGINS", "http://localhost:5173,http://localhost:5176"), ",")
	return cors.New(cors.Config{AllowOrigins: origins, AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}, AllowHeaders: []string{"Authorization", "Content-Type"}, AllowCredentials: true, MaxAge: 12 * time.Hour})
}

func AuthMiddleware(auth interfaces.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "authentication required"})
			return
		}
		claims, err := helpers.ValidateToken(strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid or expired session"})
			return
		}
		user, err := auth.CurrentUser(c.Request.Context(), claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid or inactive user"})
			return
		}
		c.Set("currentUser", user)
		c.Next()
	}
}
