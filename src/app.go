package main

import (
	"github.com/users-api/src/application"
	"github.com/users-api/src/common"
	"github.com/users-api/src/infrastructure"
	"net/http"
	"os"
)

func main() {

	userController := GetUserController()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	builder := common.NewWebServerBuilder()
	_ = builder.
		EnableRecover().
		EnableNewRelic().
		EnableLog().
		AddRouteGetUserById(userController).
		Build().
		GetWebServer().
		Listen(":" + port)
}

func GetUserController() *infrastructure.UserController {
	userHttpClient := &http.Client{}
	userRepository := infrastructure.NewHttpUserRepository(userHttpClient)
	userService := application.NewUserService(userRepository)
	userController := infrastructure.NewUserController(userService)
	return userController
}
