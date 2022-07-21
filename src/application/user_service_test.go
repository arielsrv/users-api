package application_test

import (
	"github.com/gofiber/fiber/v2"
	"github.com/users-api/src/application"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/users-api/src/domain"
)

type UserServiceUnitTestSuite struct {
	suite.Suite
	userRepository      *MockUserRepository
	userService         *application.UserService
	userErrorRepository *MockUserErrorRepository
	userErrorService    *application.UserService
}

type MockUserRepository struct {
	mock.Mock
}

type MockUserErrorRepository struct {
	mock.Mock
}

func TestUnit(t *testing.T) {
	suite.Run(t, new(UserServiceUnitTestSuite))
}

func (suite *UserServiceUnitTestSuite) SetupTest() {
	suite.userRepository = new(MockUserRepository)
	suite.userErrorRepository = new(MockUserErrorRepository)
	suite.userService = application.NewUserService(suite.userRepository)
	suite.userErrorService = application.NewUserService(suite.userErrorRepository)
}

func (mock *MockUserRepository) GetUser(int) (*domain.User, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*domain.User), nil
}

func (mock *MockUserRepository) GetUsers() ([]domain.User, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]domain.User), nil
}

func (mock *MockUserErrorRepository) GetUser(int) (*domain.User, error) {
	args := mock.Called()
	return args.Get(0).(*domain.User), args.Get(1).(error)
}

func (mock *MockUserErrorRepository) GetUsers() ([]domain.User, error) {
	args := mock.Called()
	return args.Get(0).([]domain.User), args.Get(1).(error)
}

func (suite *UserServiceUnitTestSuite) TestGetUser() {
	suite.userRepository.On("GetUser").Return(GetUser())

	actual, err := suite.userService.GetUser(1)

	suite.NotNil(actual)
	suite.NoError(err)
	suite.Equal(1, actual.Id)
	suite.Equal("John Doe", actual.Name)
	suite.Equal("john@doe.com", actual.Email)
}

func (suite *UserServiceUnitTestSuite) TestGetError() {
	suite.userErrorRepository.On("GetUser").Return(GetUserError())

	actual, err := suite.userErrorService.GetUser(1)

	suite.Nil(actual)
	suite.Error(err)
}

func (suite *UserServiceUnitTestSuite) TestGetUsersError() {
	suite.userErrorRepository.On("GetUsers").Return(GetUsersError())

	actual, err := suite.userErrorService.GetUsers()

	suite.Nil(actual)
	suite.Error(err)
}

func GetUserError() (*domain.User, error) {
	return nil, fiber.NewError(http.StatusInternalServerError, "error has ocurred")
}

func GetUsersError() ([]domain.User, error) {
	return nil, fiber.NewError(http.StatusInternalServerError, "error has ocurred")
}

func (suite *UserServiceUnitTestSuite) TestGetUsers() {
	suite.userRepository.On("GetUsers").Return(GetUsers())

	actual, err := suite.userService.GetUsers()

	suite.NotNil(actual)
	suite.NoError(err)
	suite.Len(actual, 2)
	suite.NotNil(actual[0])
	suite.Equal(1, actual[0].Id)
	suite.Equal("John Doe", actual[0].Name)
	suite.Equal("john@doe.com", actual[0].Email)
	suite.NotNil(actual[1])
	suite.Equal(2, actual[1].Id)
	suite.Equal("John Foo", actual[1].Name)
	suite.Equal("john@foo.com", actual[1].Email)
}

func GetUsers() []domain.User {
	users := make([]domain.User, 2)
	users[0] = domain.User{
		Id:    1,
		Name:  "John Doe",
		Email: "john@doe.com",
	}
	users[1] = domain.User{
		Id:    2,
		Name:  "John Foo",
		Email: "john@foo.com",
	}
	return users
}

func GetUser() *domain.User {
	return &domain.User{
		Id:    1,
		Name:  "John Doe",
		Email: "john@doe.com",
	}
}
