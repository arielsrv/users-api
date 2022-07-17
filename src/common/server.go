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
)

type WebServer struct {
	app *fiber.App
}

type Controllers struct {
	userController *infrastructure.UserController
}

func NewControllers(userController *infrastructure.UserController) *Controllers {
	return &Controllers{userController: userController}
}

func (server *WebServer) GetWebServer() *fiber.App {
	return server.app
}

type WebServerBuilder struct {
	app         *fiber.App
	controllers *Controllers
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

func (builder *WebServerBuilder) AddRoutes() *WebServerBuilder {
	builder.app.Get("/users/:id", builder.GetUserById())
	builder.app.Get("/users", builder.GetUsers())
	return builder
}

func (builder *WebServerBuilder) GetUsers() func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		usersDto, err := builder.controllers.userController.GetUsers()
		if e, ok := err.(*fiber.Error); ok {
			return ctx.Status(e.Code).SendString(err.Error())
		}
		return ctx.JSON(usersDto)
	}
}

func (builder *WebServerBuilder) GetUserById() func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		userDto, err := builder.controllers.userController.GetUser(ctx)
		if e, ok := err.(*fiber.Error); ok {
			return ctx.Status(e.Code).SendString(err.Error())
		}
		return ctx.JSON(userDto)
	}
}

func (builder *WebServerBuilder) AddControllers(controllers *Controllers) *WebServerBuilder {
	builder.controllers = controllers
	return builder
}

func (builder *WebServerBuilder) Build() *WebServer {
	return &WebServer{
		app: builder.app,
	}
}
