package infrastructure_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/users-api/src/application"
	"github.com/users-api/src/common"
	"github.com/users-api/src/infrastructure"
)

type UserControllerIntegrationSuite struct {
	suite.Suite
	app         *fiber.App
	userService *MockUserService
	port        string
}

type MockUserService struct {
	mock.Mock
}

func (suite *UserControllerIntegrationSuite) SetupTest() {
	suite.userService = new(MockUserService)
	builder := common.NewWebServerBuilder()
	suite.app = builder.
		AddControllers(common.NewControllers(
			infrastructure.NewUserController(suite.userService),
		)).
		AddRoutes().
		Build().
		GetWebServer()
}

func (mock *MockUserService) GetUser(int) *application.UserDto {
	args := mock.Called()
	result := args.Get(0)
	return result.(*application.UserDto)
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

func GetUser() *application.UserDto {
	return &application.UserDto{
		Id:    1,
		Name:  "John Doe",
		Email: "john@doe.com",
	}
}
