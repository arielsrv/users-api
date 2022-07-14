package common

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
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

func (server *WebServer) GetWebServer() *fiber.App {
	return server.app
}

type WebServerBuilder struct {
	app *fiber.App
}

func NewWebServerBuilder() *WebServerBuilder {
	return &WebServerBuilder{
		app: fiber.New(fiber.Config{
			AppName:           "users-api",
			Prefork:           true,
			EnablePrintRoutes: true,
		}),
	}
}

func (builder *WebServerBuilder) EnableLog() *WebServerBuilder {
	builder.app.Use(requestid.New())
	builder.app.Use(logger.New(logger.Config{
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}\n",
	}))
	return builder
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
