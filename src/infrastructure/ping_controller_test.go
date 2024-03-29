package infrastructure

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"github.com/users-api/src/common"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type PingControllerIntegrationSuite struct {
	suite.Suite
	app *fiber.App
}

func (suite *PingControllerIntegrationSuite) SetupTest() {
	pingController := NewPingController()
	builder := common.NewWebServerBuilder()
	suite.app = builder.
		AddRoute(http.MethodGet, "/ping", pingController.Ping).
		Build().
		App()
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(PingControllerIntegrationSuite))
}

func (suite *PingControllerIntegrationSuite) TestPing() {
	request := httptest.NewRequest(http.MethodGet, "/ping", nil)
	response, err := suite.app.Test(request)
	body, _ := io.ReadAll(response.Body)
	suite.NotNil(response)
	suite.NoError(err)
	suite.Equal(http.StatusOK, response.StatusCode)
	suite.Equal("pong", string(body))
}
