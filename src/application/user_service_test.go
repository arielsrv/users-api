package application_test

import (
	"github.com/users-api/src/application"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/users-api/src/domain"
)

type UserServiceUnitTestSuite struct {
	suite.Suite
	userRepository *MockUserRepository
	userService    *application.UserService
}

type MockUserRepository struct {
	mock.Mock
}

func TestUnit(t *testing.T) {
	suite.Run(t, new(UserServiceUnitTestSuite))
}

func (suite *UserServiceUnitTestSuite) SetupTest() {
	suite.userRepository = new(MockUserRepository)
	suite.userService = application.NewUserService(suite.userRepository)
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

func (suite *UserServiceUnitTestSuite) TestGetUser() {
	suite.userRepository.On("GetUser").Return(GetUser())

	actual, err := suite.userService.GetUser(1)

	suite.NotNil(actual)
	suite.NoError(err)
	suite.Equal(1, actual.Id)
	suite.Equal("John Doe", actual.Name)
	suite.Equal("john@doe.com", actual.Email)
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
