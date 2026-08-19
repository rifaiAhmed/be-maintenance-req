package cmd

import (
	"backend-test/helpers"
	"backend-test/internal/api"
	"backend-test/internal/interfaces"
	repository "backend-test/internal/repositories"
	service "backend-test/internal/services"
)

type Dependency struct {
	AuthAPI     *api.AuthHandler
	UserAPI     *api.UserHandler
	RequestAPI  *api.RequestHandler
	AuthService interfaces.AuthService
}

func dependencyInject() Dependency {
	userRepository := &repository.UserRepository{Database: helpers.DB}
	requestRepository := &repository.RequestRepository{Database: helpers.DB}
	authService := service.NewAuthService(userRepository)
	userService := service.NewUserService(userRepository)
	requestService := service.NewRequestService(requestRepository)
	return Dependency{AuthAPI: api.NewAuthHandler(authService), UserAPI: api.NewUserHandler(userService), RequestAPI: api.NewRequestHandler(requestService), AuthService: authService}
}
