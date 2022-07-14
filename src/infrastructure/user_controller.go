package infrastructure

import (
	"github.com/gofiber/fiber/v2"
	"github.com/users-api/src/application"
	"log"
)

type IUserController interface {
	GetUser(ctx *fiber.Ctx) *application.UserDto
}

type UserController struct {
	userService application.IUserService
}

func NewUserController(userService application.IUserService) *UserController {
	return &UserController{userService: userService}
}

func (userController UserController) GetUser(ctx *fiber.Ctx) *application.UserDto {
	userId, err := ctx.ParamsInt("id")
	if err != nil {
		log.Printf("bad request")
	}
	return userController.
		userService.
		GetUser(userId)
}
