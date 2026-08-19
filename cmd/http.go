package cmd

import (
	"backend-test/helpers"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ServeHTTP() {
	dependency := dependencyInject()
	router := gin.New()
	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatal(err)
	}
	router.Use(gin.Logger(), gin.Recovery(), MiddlewareCORS())
	router.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	apiGroup := router.Group("/api")
	apiGroup.POST("/auth/login", dependency.AuthAPI.Login)
	secured := apiGroup.Group("")
	secured.Use(AuthMiddleware(dependency.AuthService))
	secured.GET("/auth/me", dependency.AuthAPI.Me)
	secured.POST("/auth/logout", dependency.AuthAPI.Logout)
	secured.GET("/dashboard", dependency.RequestAPI.Dashboard)
	secured.GET("/requests", dependency.RequestAPI.List)
	secured.POST("/requests", dependency.RequestAPI.Create)
	secured.GET("/requests/:id", dependency.RequestAPI.Get)
	secured.PUT("/requests/:id", dependency.RequestAPI.Update)
	secured.PATCH("/requests/:id/review", dependency.RequestAPI.Review)
	secured.DELETE("/requests/:id", dependency.RequestAPI.Delete)
	secured.GET("/users", dependency.UserAPI.List)
	secured.POST("/users", dependency.UserAPI.Create)
	secured.GET("/users/:id", dependency.UserAPI.Get)
	secured.PUT("/users/:id", dependency.UserAPI.Update)
	port := helpers.GetEnv("PORT", "8080")
	log.Printf("server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
