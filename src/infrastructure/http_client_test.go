package infrastructure

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/users-api/src/domain"
	"io/ioutil"
	"net/http"
	"testing"
)

type HttpClientUnitSuite struct {
	suite.Suite
	client            *MockHttpClient
	customClient      *CustomClient
	errorClient       *MockHttpErrorClient
	customErrorClient *CustomClient
}

func (suite *HttpClientUnitSuite) SetupTest() {
	suite.client = new(MockHttpClient)
	suite.customClient = NewCustomClient(suite.client)
	suite.errorClient = new(MockHttpErrorClient)
	suite.customErrorClient = NewCustomClient(suite.errorClient)
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
	suite.client.On("Get").Return(GetHttpResponse())
	response, err := suite.customClient.Get("foo.com/users1")

	suite.NotNil(response)
	suite.NoError(err)
}

func (suite *HttpClientUnitSuite) TestGetGenericError() {
	suite.errorClient.On("Get").Return(GenericErrorResponse())
	_, err := suite.customErrorClient.Get("foo.com/users1")

	suite.Error(err)
}

func (suite *HttpClientUnitSuite) TestGetError() {
	suite.client.On("Get").Return(ErrorResponse())
	_, err := suite.customClient.Get("foo.com/users1")

	suite.Error(err)
}

func GenericErrorResponse() (*http.Response, error) {
	return nil, fiber.NewError(http.StatusInternalServerError, "")
}

func ErrorResponse() (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       nil,
	}, nil
}

func GetHttpResponse() (*http.Response, error) {
	user := domain.User{
		Id:    1,
		Name:  "John Doe",
		Email: "john@doe.com",
	}
	binary, err := json.Marshal(user)
	var httpResponse = http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBuffer(binary)),
	}
	return &httpResponse, err
}
