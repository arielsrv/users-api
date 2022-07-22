package main

import (
	"github.com/users-api/src/application"
	"github.com/users-api/src/common"
	"github.com/users-api/src/infrastructure"
	"net/http"
	"os"
)

func main() {
	httpClientProxy := infrastructure.NewHTTPClientProxy(&http.Client{})
	userRepository := infrastructure.NewHTTPUserRepository(httpClientProxy)
	userService := application.NewUserService(userRepository)
	userController := infrastructure.NewUserController(userService)

	builder := common.NewWebServerBuilder()

	builder.
		AddRoute("GET", "/users/:id", userController.GetUser).
		AddRoute("GET", "/users", userController.GetUsers)

	builder.
		EnableLog().
		EnableRecover().
		EnableNewRelic()

	address, port := os.LookupEnv("PORT")
	if port {
		builder.Listen(address)
	} else {
		builder.Listen(":3000")
	}

	_, prefork := os.LookupEnv("PREFORK")
	if prefork {
		builder.Prefork()
	}

	builder.
		Build().
		Start()
}
