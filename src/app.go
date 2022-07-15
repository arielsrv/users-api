package main

import (
	"github.com/users-api/src/application"
	"github.com/users-api/src/common"
	"github.com/users-api/src/infrastructure"
	"net/http"
)

func main() {

	userController := GetUserController()

	builder := common.NewWebServerBuilder()
	_ = builder.
		EnableRecover().
		EnableLog().
		AddRouteGetUserById(userController).
		Build().
		GetWebServer().
		Listen(":3000")
}

func GetUserController() *infrastructure.UserController {
	userHttpClient := &http.Client{}
	userRepository := infrastructure.NewHttpUserRepository(userHttpClient)
	userService := application.NewUserService(userRepository)
	userController := infrastructure.NewUserController(userService)
	return userController
}
