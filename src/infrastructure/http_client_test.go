package infrastructure

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type HttpClientUnitSuite struct {
	suite.Suite
	client      *MockHttpClient
	errorClient *MockHttpErrorClient
}

func (suite *HttpClientUnitSuite) SetupTest() {
	suite.client = new(MockHttpClient)
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
