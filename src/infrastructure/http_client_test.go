package infrastructure

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type HTTPClientUnitSuite struct {
	suite.Suite
	client      *MockHTTPClient
	proxy       *HTTPClientProxy
	errorClient *MockHTTPErrorClient
	errorProxy  *HTTPClientProxy
}

func (suite *HTTPClientUnitSuite) SetupTest() {
	suite.client = new(MockHTTPClient)
	suite.errorClient = new(MockHTTPErrorClient)
	suite.proxy = NewHTTPClientProxy(suite.client, "users")
	suite.errorProxy = NewHTTPClientProxy(suite.errorClient, "users")
}

func TestHttpClientUnit(t *testing.T) {
	suite.Run(t, new(HTTPClientUnitSuite))
}

type MockHTTPClient struct {
	mock.Mock
}

type MockHTTPErrorClient struct {
	mock.Mock
}

func (mock *MockHTTPClient) Get(string) (response *http.Response, err error) {
	args := mock.Called()
	return args.Get(0).(*http.Response), err
}

func (mock *MockHTTPErrorClient) Get(string) (response *http.Response, err error) {
	args := mock.Called()
	return args.Get(0).(*http.Response), args.Get(1).(error)
}

func (suite *HTTPClientUnitSuite) TestGet() {
	suite.client.On("Get").Return(Get())
	actual, err := suite.proxy.Get("foo.com/users/1")
	suite.NotNil(actual)
	suite.NoError(err)
}

func (suite *HTTPClientUnitSuite) TestGetError() {
	suite.errorClient.On("Get").Return(GetError())
	actual, err := suite.errorProxy.Get("foo.com/users/1")
	suite.Nil(actual)
	suite.Error(err)
}

func GetError() (*http.Response, error) {
	return nil, fiber.NewError(http.StatusInternalServerError, "error has ocurred. ")
}
