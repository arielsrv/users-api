package infrastructure_test

import (
	"github.com/users-api/src/infrastructure"
	"io/ioutil"
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
		AddRoute("GET", "/users/:id", userController.GetUser).
		AddRoute("GET", "/users", userController.GetUsers).
		Build().
		App()
}

func (mock *MockUserService) GetUser(int) (*application.UserDto, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*application.UserDto), nil
}

func (mock *MockUserService) GetUsers() ([]application.UserDto, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]application.UserDto), nil
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(UserControllerIntegrationSuite))
}

func (suite *UserControllerIntegrationSuite) Test_Get_User_By_Id() {
	suite.userService.On("GetUser").Return(GetUser())

	request := httptest.NewRequest("GET", "/users/1", nil)
	response, err := suite.app.Test(request)
	body, _ := ioutil.ReadAll(response.Body)

	suite.NotNil(response)
	suite.NoError(err)
	suite.Equal(http.StatusOK, response.StatusCode)
	suite.Equal(`{"id":1,"name":"John Doe","email":"john@doe.com"}`, string(body))
}

func (suite *UserControllerIntegrationSuite) Test_Get_User_By_Id_Bad_Request() {
	suite.userService.On("GetUser").Return(GetUser())

	request := httptest.NewRequest("GET", "/users/a", nil)
	response, err := suite.app.Test(request)

	suite.NotNil(response)
	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, response.StatusCode)
}

func (suite *UserControllerIntegrationSuite) Test_Get_Users() {
	suite.userService.On("GetUsers").Return(GetUsers())

	request := httptest.NewRequest("GET", "/users", nil)
	response, err := suite.app.Test(request)
	body, _ := ioutil.ReadAll(response.Body)

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
