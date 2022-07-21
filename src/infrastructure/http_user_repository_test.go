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

type HTTPUserRepositoryUnitSuite struct {
	suite.Suite
	client              *MockClient
	errorClient         *MockErrorClient
	userRepository      *HTTPUserRepository
	userErrorRepository *HTTPUserRepository
}

func (suite *HTTPUserRepositoryUnitSuite) SetupTest() {
	suite.client = new(MockClient)
	suite.errorClient = new(MockErrorClient)
	suite.userRepository = NewHTTPUserRepository(suite.client)
	suite.userErrorRepository = NewHTTPUserRepository(suite.errorClient)
}

func TestUnit(t *testing.T) {
	suite.Run(t, new(HTTPUserRepositoryUnitSuite))
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

func (suite *HTTPUserRepositoryUnitSuite) TestGet() {
	suite.client.On("Get").Return(Get())

	actual, err := suite.userRepository.GetUser(1)

	suite.NotNil(actual)
	suite.NoError(err)
	suite.Equal(1, actual.ID)
	suite.Equal("John Doe", actual.Name)
	suite.Equal("john@doe.com", actual.Email)
}

func (suite *HTTPUserRepositoryUnitSuite) TestGetUsers() {
	suite.client.On("Get").Return(GetUsers())

	actual, err := suite.userRepository.GetUsers()

	suite.NotNil(actual)
	suite.NoError(err)
	suite.Len(actual, 2)
	suite.NotNil(actual[0])
	suite.Equal(1, actual[0].ID)
	suite.Equal("John Doe", actual[0].Name)
	suite.Equal("john@doe.com", actual[0].Email)
	suite.NotNil(actual[1])
	suite.Equal(2, actual[1].ID)
	suite.Equal("John Foo", actual[1].Name)
	suite.Equal("john@foo.com", actual[1].Email)
}

func GetUsers() (*http.Response, error) {
	users := make([]domain.User, 2)
	users[0] = domain.User{
		ID:    1,
		Name:  "John Doe",
		Email: "john@doe.com",
	}
	users[1] = domain.User{
		ID:    2,
		Name:  "John Foo",
		Email: "john@foo.com",
	}
	binary, _ := json.Marshal(users)
	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBuffer(binary)),
	}, nil
}

func (suite *HTTPUserRepositoryUnitSuite) TestError() {
	suite.errorClient.On("Get").Return(GetError())

	_, err := suite.userErrorRepository.GetUser(1)

	suite.Error(err)
}

func (suite *HTTPUserRepositoryUnitSuite) TestNotOk() {
	suite.client.On("Get").Return(GetNotFound())

	_, err := suite.userRepository.GetUser(1)

	suite.Error(err)
}

func GetNotFound() (*http.Response, error) {
	return &http.Response{StatusCode: http.StatusNotFound}, nil
}

func Get() (*http.Response, error) {
	user := domain.User{
		ID:    1,
		Name:  "John Doe",
		Email: "john@doe.com",
	}
	binary, _ := json.Marshal(user)
	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBuffer(binary)),
	}, nil
}
