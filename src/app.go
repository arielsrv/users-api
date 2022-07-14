package main

import (
	"github.com/users-api/src/application"
	"github.com/users-api/src/common"
	"github.com/users-api/src/infrastructure"
)

func main() {

	userController := GetUserController()

	builder := common.NewWebServerBuilder()
	_ = builder.
		EnableLog().
		AddRouteGetUserById(userController).
		Build().
		GetWebServer().
		Listen(":3000")
}

func GetUserController() *infrastructure.UserController {
	userRepository := infrastructure.NewHttpUserRepository()
	userService := application.NewUserService(userRepository)
	userController := infrastructure.NewUserController(userService)
	return userController
}
