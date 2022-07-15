package common

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/nobuyo/nrfiber"
	"github.com/users-api/src/infrastructure"
	"os"
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
			Prefork:           false,
			EnablePrintRoutes: true,
		}),
	}
}

func (builder *WebServerBuilder) EnableRecover() *WebServerBuilder {
	var config = recover.Config{
		EnableStackTrace: true,
	}
	builder.app.Use(recover.New(config))
	return builder
}

func (builder *WebServerBuilder) EnableLog() *WebServerBuilder {
	builder.app.Use(requestid.New())
	builder.app.Use(logger.New(logger.Config{
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}\n",
	}))
	return builder
}

func (builder *WebServerBuilder) EnableNewRelic() *WebServerBuilder {
	nrapp, _ := newrelic.NewApplication(
		newrelic.ConfigAppName("golang-users-api"),
		newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
		newrelic.ConfigDebugLogger(os.Stdout),
	)

	builder.app.Use(nrfiber.New(nrfiber.Config{
		NewRelicApp: nrapp,
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
