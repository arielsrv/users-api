package application

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/users-api/src/domain"
	"testing"
)

type UserServiceUnitTestSuite struct {
	suite.Suite
	userRepository *MockUserRepository
	userService    *UserService
}

type MockUserRepository struct {
	mock.Mock
}

func TestUnit(t *testing.T) {
	suite.Run(t, new(UserServiceUnitTestSuite))
}

func (suite *UserServiceUnitTestSuite) SetupTest() {
	suite.userRepository = new(MockUserRepository)
	suite.userService = NewUserService(suite.userRepository)
}

func (mock *MockUserRepository) GetUser(id int) *domain.User {
	args := mock.Called()
	result := args.Get(0)
	return result.(*domain.User)
}

func (suite *UserServiceUnitTestSuite) TestGetUser() {
	suite.userRepository.On("GetUser").Return(GetUser())

	actual := suite.userService.GetUser(1)

	suite.NotNil(actual)
	suite.Equal(1, actual.Id)
	suite.Equal("John Doe", actual.Name)
	suite.Equal("john@doe.com", actual.Email)
}

func GetUser() *domain.User {
	return &domain.User{
		Id:    1,
		Name:  "John Doe",
		Email: "john@doe.com",
	}
}
