package main

import (
	"github.com/users-api/src/application"
	"github.com/users-api/src/common"
	"github.com/users-api/src/infrastructure"
	"net/http"
	"os"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	controllers := common.NewControllers(
		GetUserController(),
	)

	prefork := os.Getenv("PREFORK")
	builder := common.NewWebServerBuilder(prefork)
	_ = builder.
		EnableRecover().
		EnableNewRelic().
		EnableLog().
		AddControllers(controllers).
		AddRoutes().
		Build().
		GetWebServer().
		Listen(":" + port)
}

func GetUserController() *infrastructure.UserController {
	httpClientProxy := infrastructure.NewHttpClientProxy(http.Client{})
	userRepository := infrastructure.NewHttpUserRepository(httpClientProxy)
	userService := application.NewUserService(userRepository)
	userController := infrastructure.NewUserController(userService)
	return userController
}
