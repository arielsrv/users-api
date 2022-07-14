package common

import (
	"github.com/gofiber/fiber/v2"
	"github.com/users-api/src/infrastructure"
	"reflect"
)

var (
	routes = struct {
		path map[string]string
	}{path: map[string]string{
		"*infrastructure.UserController": "/users/:id",
	}}
)

type WebServer struct {
	app *fiber.App
}

func (server WebServer) GetWebServer() *fiber.App {
	return server.app
}

type WebServerBuilder struct {
	app *fiber.App
}

func NewWebServerBuilder() *WebServerBuilder {

	routes := make(map[string]string)
	routes["*infrastructure.UserController"] = "/users/:id"

	return &WebServerBuilder{
		app: fiber.New(fiber.Config{
			AppName: "users-api",
		}),
	}
}

func (builder *WebServerBuilder) AddRouteGetUserById(controller infrastructure.IUserController) *WebServerBuilder {
	builder.app.Get(routes.path[reflect.TypeOf(controller).String()], func(ctx *fiber.Ctx) error {
		userDto := controller.GetUser(ctx)
		return ctx.JSON(userDto)
	})
	return builder
}

func (builder *WebServerBuilder) Build() *WebServer {
	return &WebServer{
		app: builder.app,
	}
}
