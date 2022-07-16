package infrastructure

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/users-api/src/domain"
	"io/ioutil"
	"net/http"
	"testing"
)

type HttpUserRepositoryUnitSuite struct {
	suite.Suite
	client              *MockClient
	userRepository      *HttpUserRepository
	errorClient         *MockErrorClient
	userErrorRepository *HttpUserRepository
}

func (suite *HttpUserRepositoryUnitSuite) SetupTest() {
	suite.client = new(MockClient)
	suite.errorClient = new(MockErrorClient)
	suite.userRepository = NewHttpUserRepository(suite.client)
	suite.userErrorRepository = NewHttpUserRepository(suite.errorClient)
}

func TestUnit(t *testing.T) {
	suite.Run(t, new(HttpUserRepositoryUnitSuite))
}

type MockClient struct {
	mock.Mock
}

func (mock *MockClient) Get(string) (response *http.Response, err error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*http.Response), nil
}

type MockErrorClient struct {
	mock.Mock
}

func (mock *MockErrorClient) Get(string) (response *http.Response, err error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*http.Response), errors.New("A error has ocurred. ")
}

func (suite *HttpUserRepositoryUnitSuite) TestGet() {
	suite.client.On("Get").Return(Get())

	actual, err := suite.userRepository.GetUser(1)

	suite.NotNil(actual)
	suite.NoError(err)
	suite.Equal(1, actual.Id)
	suite.Equal("John Doe", actual.Name)
	suite.Equal("john@doe.com", actual.Email)
}

func (suite *HttpUserRepositoryUnitSuite) TestGet_NotFound() {
	suite.client.On("Get").Return(GetNotFound())

	actual, err := suite.userRepository.GetUser(1)

	suite.Nil(actual)
	suite.Error(err)
	if e, ok := err.(*fiber.Error); ok {
		suite.Equal(http.StatusNotFound, e.Code)
		suite.Equal("Couldn't retreive user with id 1 not found. ", e.Message)
	}
}

func (suite *HttpUserRepositoryUnitSuite) TestGet_InternalServerError() {
	suite.errorClient.On("Get").Return(&http.Response{})

	actual, err := suite.userErrorRepository.GetUser(1)

	suite.Nil(actual)
	suite.Error(err)
	if e, ok := err.(*fiber.Error); ok {
		suite.Equal(http.StatusInternalServerError, e.Code)
		suite.Equal("A error has ocurred. ", e.Message)
	}
}

func (suite *HttpUserRepositoryUnitSuite) TestGetUsers_NotFound() {
	suite.client.On("Get").Return(GetNotFound())

	actual, err := suite.userRepository.GetUsers()

	suite.Nil(actual)
	suite.Error(err)
	if e, ok := err.(*fiber.Error); ok {
		suite.Equal(http.StatusNotFound, e.Code)
	}
}

func (suite *HttpUserRepositoryUnitSuite) TestGetUsers_InternalServerError() {
	suite.errorClient.On("Get").Return(&http.Response{})

	actual, err := suite.userErrorRepository.GetUsers()

	suite.Nil(actual)
	suite.Error(err)
	if e, ok := err.(*fiber.Error); ok {
		suite.Equal(http.StatusInternalServerError, e.Code)
		suite.Equal("A error has ocurred. ", e.Message)
	}
}

func (suite *HttpUserRepositoryUnitSuite) TestGetAll() {
	suite.client.On("Get").Return(GetAll())

	actual, err := suite.userRepository.GetUsers()

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

func Get() (response *http.Response, err error) {
	user := domain.User{
		Id:    1,
		Name:  "John Doe",
		Email: "john@doe.com",
	}
	binary, err := json.Marshal(user)
	return buildResponse(binary)
}

func GetNotFound() (response *http.Response, err error) {
	return &http.Response{
		StatusCode: http.StatusNotFound,
	}, nil
}

func GetAll() (response *http.Response, err error) {
	user1 := domain.User{
		Id:    1,
		Name:  "John Doe",
		Email: "john@doe.com",
	}
	user2 := domain.User{
		Id:    2,
		Name:  "John Foo",
		Email: "john@foo.com",
	}
	var users = make([]domain.User, 2)
	users[0] = user1
	users[1] = user2
	binary, err := json.Marshal(users)
	return buildResponse(binary)
}

func buildResponse(binary []byte) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBuffer(binary)),
	}, nil
}
