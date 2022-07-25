package main

import (
	"fmt"
	"github.com/users-api/src/application"
	"github.com/users-api/src/common"
	"github.com/users-api/src/infrastructure"
	http "net/http"
	"os"
)

func main() {
	httpClientProxy := infrastructure.NewHTTPClientProxy(&http.Client{})
	userRepository := infrastructure.NewHTTPUserRepository(httpClientProxy)
	userService := application.NewUserService(userRepository)
	userController := infrastructure.NewUserController(userService)

	pingController := infrastructure.NewPingController()

	builder := common.NewWebServerBuilder()

	builder.AddRoute(http.MethodGet, "/ping", pingController.Ping)
	builder.AddRoute(http.MethodGet, "/users", userController.GetUsers)
	builder.AddRoute(http.MethodGet, "/users/:id", userController.GetUser)

	builder.UseLog()
	builder.UseRecover()
	builder.UseNewRelic()

	address, port := os.LookupEnv("PORT")
	if port {
		builder.Listen(fmt.Sprintf(":%s", address))
	} else {
		builder.UseDefaultAddress()
	}

	_, prefork := os.LookupEnv("PREFORK")
	if prefork {
		builder.Prefork()
	}

	builder.
		Build().
		Start()
}
