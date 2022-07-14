package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/users-api/src/application"
	"github.com/users-api/src/infrastructure"
)

type WebServer struct {
	app *fiber.App
}

func NewWebServer() *WebServer {
	return &WebServer{
		app: fiber.New(),
	}
}

func (server WebServer) GetWebServer() *fiber.App {
	app := fiber.New()

	userRepository := infrastructure.NewUserRepository()
	userService := application.NewUserService(userRepository)
	userController := infrastructure.NewUserController(userService)

	app.Get("/users/:id", func(ctx *fiber.Ctx) error {
		var userDto = userController.GetUser(ctx)
		return ctx.JSON(userDto)
	})

	return app
}
