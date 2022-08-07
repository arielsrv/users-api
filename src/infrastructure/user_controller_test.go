package infrastructure_test

import (
	"github.com/users-api/src/infrastructure"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/users-api/src/application"
	"github.com/users-api/src/common"
)

type UserControllerIntegrationSuite struct {
	suite.Suite
	app         *fiber.App
	userService *MockUserService
}

type MockUserService struct {
	mock.Mock
}

func (suite *UserControllerIntegrationSuite) SetupTest() {
	suite.userService = new(MockUserService)
	userController := infrastructure.NewUserController(suite.userService)
	builder := common.NewWebServerBuilder()
	suite.app = builder.
		AddRoute(http.MethodGet, "/users/:id", userController.GetUser).
		AddRoute(http.MethodGet, "/users", userController.GetAll).
		Build().
		App()
}

func (mock *MockUserService) GetByID(int) (*application.UserDto, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*application.UserDto), nil
}

func (mock *MockUserService) GetAll() ([]application.UserDto, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]application.UserDto), nil
}

func (mock *MockUserService) MultiGetByID([]int) (*[]application.MultiGetDto[application.UserDto], error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*[]application.MultiGetDto[application.UserDto]), nil
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(UserControllerIntegrationSuite))
}

func (suite *UserControllerIntegrationSuite) Test_Get_User_By_Id() {
	suite.userService.On("GetByID").Return(GetUser())

	request := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	response, err := suite.app.Test(request)
	body, _ := io.ReadAll(response.Body)

	suite.NotNil(response)
	suite.NoError(err)
	suite.Equal(http.StatusOK, response.StatusCode)
	suite.Equal(`{"id":1,"name":"John Doe","email":"john@doe.com"}`, string(body))
}

func (suite *UserControllerIntegrationSuite) Test_Get_User_By_Id_Bad_Request() {
	suite.userService.On("GetByID").Return(GetUser())

	request := httptest.NewRequest(http.MethodGet, "/users/a", nil)
	response, err := suite.app.Test(request)

	suite.NotNil(response)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, response.StatusCode)
}

func (suite *UserControllerIntegrationSuite) Test_Get_Users() {
	suite.userService.On("GetAll").Return(GetUsers())

	request := httptest.NewRequest(http.MethodGet, "/users", nil)
	response, err := suite.app.Test(request)
	body, _ := io.ReadAll(response.Body)

	suite.NotNil(response)
	suite.NoError(err)
	suite.Equal(http.StatusOK, response.StatusCode)
	suite.Equal(`[{"id":1,"name":"John Doe","email":"john@doe.com"},{"id":2,"name":"John Foo","email":"john@foo.com"}]`, string(body))
}

func GetUsers() []application.UserDto {
	usersDto := make([]application.UserDto, 2)
	usersDto[0] = application.UserDto{
		ID:    1,
		Name:  "John Doe",
		Email: "john@doe.com",
	}
	usersDto[1] = application.UserDto{
		ID:    2,
		Name:  "John Foo",
		Email: "john@foo.com",
	}
	return usersDto
}

func GetUser() *application.UserDto {
	return &application.UserDto{
		ID:    1,
		Name:  "John Doe",
		Email: "john@doe.com",
	}
}
