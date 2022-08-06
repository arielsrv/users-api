package infrastructure

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/newrelic/go-agent/v3/newrelic"
	"io"
	"net/http"
	"os"
	"time"
)

type HTTPClient interface {
	Get(url string) (response *http.Response, err error)
}

type HTTPClientProxy struct {
	client HTTPClient
	name   string
}

func NewHTTPClientProxy(client HTTPClient, name string) *HTTPClientProxy {
	return &HTTPClientProxy{
		client: client,
		name:   name,
	}
}

func (customHttpClient HTTPClientProxy) Get(url string) (response *http.Response, err error) {
	var start = time.Now()
	response, err = customHttpClient.client.Get(url)
	var end = time.Since(start)
	nrapp, _ := newrelic.NewApplication(
		newrelic.ConfigAppName("golang-users-api"),
		newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
		newrelic.ConfigDebugLogger(os.Stdout),
	)
	metric := fmt.Sprintf("%s-client", customHttpClient.name)
	nrapp.RecordCustomMetric(metric, float64(end))
	return response, err
}

type Client[T any] struct {
	client HTTPClient
}

func (httpClient Client[T]) Get(url string) (T, error) {
	var reference T
	response, err := httpClient.client.Get(url)

	if err != nil {
		err = fiber.NewError(http.StatusInternalServerError, err.Error())
		return reference, err
	}

	if response.StatusCode != http.StatusOK {
		err = fiber.NewError(response.StatusCode, "Couldn't retreive. ")
		return reference, err
	}

	body, err := io.ReadAll(response.Body)
	defer func(Body io.ReadCloser) {
		err = Body.Close()
	}(response.Body)
	err = json.Unmarshal(body, &reference)

	return reference, err
}
