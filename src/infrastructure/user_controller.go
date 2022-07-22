package infrastructure

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/users-api/src/application"
	"net/http"
)

type IUserController interface {
	GetUser(ctx *fiber.Ctx) error
	GetUsers(ctx *fiber.Ctx) error
}

type UserController struct {
	userService application.IUserService
}

func NewUserController(userService application.IUserService) *UserController {
	return &UserController{userService: userService}
}

func (userController UserController) GetUser(ctx *fiber.Ctx) error {
	userID, err := ctx.ParamsInt("id")
	if err != nil {
		err := NewBadRequest(fmt.Sprintf("Invalid format for userId, %s", ctx.Params("id")))
		return err
	}

	var result, _ = userController.
		userService.
		GetUser(userID)

	return ctx.JSON(result)
}

func NewBadRequest(message string) error {
	err := fiber.NewError(http.StatusBadRequest, message)
	return err
}

func (userController UserController) GetUsers(ctx *fiber.Ctx) error {
	result, err := userController.
		userService.
		GetUsers()

	if err != nil {
		return err
	}

	return ctx.JSON(result)
}
