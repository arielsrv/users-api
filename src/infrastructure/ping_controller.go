package infrastructure

import "github.com/gofiber/fiber/v2"

type IPingController interface {
	Ping(ctx *fiber.Ctx) error
}

type PingController struct {
}

func NewPingController() *PingController {
	return &PingController{}
}

func (pingController PingController) Ping(ctx *fiber.Ctx) error {
	return ctx.SendString("pong")
}
