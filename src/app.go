package main

import (
	"fmt"
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

	pingController := infrastructure.PingController{}

	builder := common.NewWebServerBuilder()

	builder.
		AddRoute("GET", "/users/:id", userController.GetUser).
		AddRoute("GET", "/users", userController.GetUsers).
		AddRoute("GET", "/ping", pingController.Ping)

	builder.
		EnableLog().
		EnableRecover().
		EnableNewRelic()

	address, port := os.LookupEnv("PORT")
	if port {
		builder.Listen(fmt.Sprintf(":%s", address))
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
