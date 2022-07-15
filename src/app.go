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

	builder := common.NewWebServerBuilder()
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
	userHttpClient := &http.Client{}
	userRepository := infrastructure.NewHttpUserRepository(userHttpClient)
	userService := application.NewUserService(userRepository)
	userController := infrastructure.NewUserController(userService)
	return userController
}
