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
