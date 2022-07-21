package infrastructure

import (
	"bytes"
	"encoding/json"
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
	errorClient         *MockErrorClient
	userRepository      *HttpUserRepository
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

type MockErrorClient struct {
	mock.Mock
}

func (mock *MockClient) Get(string) (response *http.Response, err error) {
	args := mock.Called()
	return args.Get(0).(*http.Response), err
}

func (mock *MockErrorClient) Get(string) (response *http.Response, err error) {
	args := mock.Called()
	return args.Get(0).(*http.Response), args.Get(1).(error)
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

func (suite *HttpUserRepositoryUnitSuite) TestGetUsers() {
	suite.client.On("Get").Return(GetUsers())

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

func GetUsers() (*http.Response, error) {
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
	binary, _ := json.Marshal(users)
	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBuffer(binary)),
	}, nil
}

func (suite *HttpUserRepositoryUnitSuite) TestError() {
	suite.errorClient.On("Get").Return(GetError())

	_, err := suite.userErrorRepository.GetUser(1)

	suite.Error(err)
}

func (suite *HttpUserRepositoryUnitSuite) TestNotOk() {
	suite.client.On("Get").Return(GetNotFound())

	_, err := suite.userRepository.GetUser(1)

	suite.Error(err)
}

func GetNotFound() (*http.Response, error) {
	return &http.Response{StatusCode: http.StatusNotFound}, nil
}

func Get() (*http.Response, error) {
	user := domain.User{
		Id:    1,
		Name:  "John Doe",
		Email: "john@doe.com",
	}
	binary, _ := json.Marshal(user)
	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBuffer(binary)),
	}, nil
}
