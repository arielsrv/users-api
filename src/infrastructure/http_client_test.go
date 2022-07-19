package infrastructure

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type HttpClientUnitSuite struct {
	suite.Suite
	client      *MockHttpClient
	proxy       *HttpClientProxy
	errorClient *MockHttpErrorClient
	errorProxy  *HttpClientProxy
}

func (suite *HttpClientUnitSuite) SetupTest() {
	suite.client = new(MockHttpClient)
	suite.errorClient = new(MockHttpErrorClient)
	suite.proxy = NewHttpClientProxy(suite.client)
	suite.errorProxy = NewHttpClientProxy(suite.errorClient)
}

func TestHttpClientUnit(t *testing.T) {
	suite.Run(t, new(HttpClientUnitSuite))
}

type MockHttpClient struct {
	mock.Mock
}

type MockHttpErrorClient struct {
	mock.Mock
}

func (mock *MockHttpClient) Get(string) (response *http.Response, err error) {
	args := mock.Called()
	return args.Get(0).(*http.Response), err
}

func (mock *MockHttpErrorClient) Get(string) (response *http.Response, err error) {
	args := mock.Called()
	return args.Get(0).(*http.Response), args.Get(1).(error)
}

func (suite *HttpClientUnitSuite) TestGet() {
	suite.client.On("Get").Return(Get())
	actual, err := suite.proxy.Get("foo.com/users/1")
	suite.NotNil(actual)
	suite.NoError(err)
}

func (suite *HttpClientUnitSuite) TestGetError() {
	suite.errorClient.On("Get").Return(GetError())
	actual, err := suite.errorProxy.Get("foo.com/users/1")
	suite.Nil(actual)
	suite.Error(err)
}

func GetError() (*http.Response, error) {
	return nil, fiber.NewError(http.StatusInternalServerError, "error has ocurred. ")
}
