package common

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/nobuyo/nrfiber"
	"os"
)

type WebServer struct {
	app  *fiber.App
	addr string
}

func (server *WebServer) App() *fiber.App {
	return server.app
}

func (server *WebServer) Start() {
	_ = server.app.Listen(server.addr)
}

type WebServerBuilder struct {
	routes         []Route
	enableLog      bool
	enableNewRelic bool
	enableRecover  bool
	addr           string
	prefork        bool
}

func NewWebServerBuilder() *WebServerBuilder {
	return &WebServerBuilder{
		routes: make([]Route, 0),
	}
}

func (builder *WebServerBuilder) EnableRecover() *WebServerBuilder {
	builder.enableRecover = true
	return builder
}

func (builder *WebServerBuilder) EnableLog() *WebServerBuilder {
	builder.enableLog = true
	return builder
}

func (builder *WebServerBuilder) EnableNewRelic() *WebServerBuilder {
	builder.enableNewRelic = true
	return builder
}

type Route struct {
	Method string
	Path   string
	Action func(ctx *fiber.Ctx) error
}

func (builder *WebServerBuilder) AddRoute(method string, path string, action func(ctx *fiber.Ctx) error) *WebServerBuilder {
	builder.routes = append(builder.routes, Route{
		Method: method,
		Path:   path,
		Action: action,
	})
	return builder
}

func (builder *WebServerBuilder) Listen(address string) *WebServerBuilder {
	builder.addr = address
	return builder
}

func (builder *WebServerBuilder) Prefork() *WebServerBuilder {
	builder.prefork = true
	return builder
}

func (builder *WebServerBuilder) Build() *WebServer {
	app := fiber.New(fiber.Config{
		AppName:           "users-api",
		Prefork:           builder.prefork,
		EnablePrintRoutes: false,
	})

	if builder.enableLog {
		app.Use(requestid.New())
		app.Use(logger.New(logger.Config{
			Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}\n",
		}))
	}

	if builder.enableNewRelic {
		nrapp, _ := newrelic.NewApplication(
			newrelic.ConfigAppName("golang-users-api"),
			newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
			newrelic.ConfigDebugLogger(os.Stdout),
		)

		app.Use(nrfiber.New(nrfiber.Config{
			NewRelicApp: nrapp,
		}))
	}

	if builder.enableRecover {
		app.Use(recover.New(recover.Config{
			EnableStackTrace: true,
		}))
	}

	for _, route := range builder.routes {
		app.Add(route.Method, route.Path, route.Action)
	}

	return &WebServer{
		app:  app,
		addr: builder.addr,
	}
}
