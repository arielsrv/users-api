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
		AddRoute(http.MethodGet, "/users", userController.GetAll).
		AddRoute(http.MethodGet, "/users/multi-get", userController.MultiGet).
		AddRoute(http.MethodGet, "/users/:id", userController.GetUser).
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

func (mock *MockUserService) MultiGetByID([]int) ([]application.MultiGetDto[application.UserDto], error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]application.MultiGetDto[application.UserDto]), nil
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

func (suite *UserControllerIntegrationSuite) Test_MultiGet_Users() {
	suite.userService.On("MultiGetByID").Return(GetMultiGetUsers())

	request := httptest.NewRequest(http.MethodGet, "/users/multi-get?ids=1,2,,2", nil)
	response, err := suite.app.Test(request)
	body, _ := io.ReadAll(response.Body)

	suite.NotNil(response)
	suite.NoError(err)
	suite.Equal(http.StatusOK, response.StatusCode)
	suite.Equal(`[{"code":200,"body":{"id":1,"name":"John Doe","email":"john@doe.com"}},{"code":200,"body":{"id":2,"name":"John Foo","email":"john@foo.com"}}]`, string(body))
}

func (suite *UserControllerIntegrationSuite) Test_MultiGet_Users_BadRequest_1() {
	suite.userService.On("MultiGetByID").Return(GetMultiGetUsers())

	request := httptest.NewRequest(http.MethodGet, "/users/multi-get?ids=", nil)
	response, err := suite.app.Test(request)

	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, response.StatusCode)
}

func (suite *UserControllerIntegrationSuite) Test_MultiGet_Users_BadRequest_2() {
	suite.userService.On("MultiGetByID").Return(GetMultiGetUsers())

	request := httptest.NewRequest(http.MethodGet, "/users/multi-get?ids=a,1", nil)
	response, err := suite.app.Test(request)

	suite.NoError(err)
	suite.Equal(http.StatusBadRequest, response.StatusCode)
}

func GetMultiGetUsers() []application.MultiGetDto[application.UserDto] {
	usersDto := GetUsers()
	result := make([]application.MultiGetDto[application.UserDto], 2)
	result[0].Code = 200
	result[0].Body = usersDto[0]
	result[1].Code = 200
	result[1].Body = usersDto[1]
	return result
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
