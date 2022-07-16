package infrastructure

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/users-api/src/application"
	"net/http"
)

type IUserController interface {
	GetUser(ctx *fiber.Ctx) *application.UserDto
	GetUsers(ctx *fiber.Ctx) []application.UserDto
}

type UserController struct {
	userService application.IUserService
}

func NewUserController(userService application.IUserService) *UserController {
	return &UserController{userService: userService}
}

func (userController UserController) GetUser(ctx *fiber.Ctx) (*application.UserDto, error) {
	userId, err := ctx.ParamsInt("id")
	if err != nil {
		err = NewBadRequest(fmt.Sprintf("Invalid format for userId, %s", ctx.Params("id")))
		return nil, err
	}
	return userController.
		userService.
		GetUser(userId)
}

func NewBadRequest(message string) error {
	err := fiber.NewError(http.StatusBadRequest, message)
	return err
}

func (userController UserController) GetUsers() ([]application.UserDto, error) {
	return userController.
		userService.
		GetUsers()
}
