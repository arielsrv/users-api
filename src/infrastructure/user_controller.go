package infrastructure

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/users-api/src/application"
	"github.com/users-api/src/common"
	"net/http"
)

type IUserController interface {
	GetUser(ctx *fiber.Ctx) error
	GetAll(ctx *fiber.Ctx) error
	MultiGet(ctx *fiber.Ctx) error
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
		return NewBadRequest(fmt.Sprintf("Invalid format for userId, %s", ctx.Params("id")))
	}

	result, _ := userController.
		userService.
		GetByID(userID)

	return ctx.JSON(result)
}

func (userController UserController) MultiGet(ctx *fiber.Ctx) error {
	param := ctx.Query("ids")
	if common.IsEmptyString(param) {
		return NewBadRequest("Invalid format for ids.")
	}

	ids, err := common.ToEnumerableIds(param)
	if err != nil {
		return NewBadRequest("Invalid format for ids.")
	}

	result, _ := userController.
		userService.
		MultiGetByID(ids)

	return ctx.JSON(result)
}

func NewBadRequest(message string) error {
	err := fiber.NewError(http.StatusBadRequest, message)
	return err
}

func (userController UserController) GetAll(ctx *fiber.Ctx) error {
	result, err := userController.
		userService.
		GetAll()

	if err != nil {
		return err
	}

	return ctx.JSON(result)
}
