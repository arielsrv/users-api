package infrastructure

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http"
)

type HTTPClient interface {
	Get(url string) (response *http.Response, err error)
}

type HTTPClientProxy struct {
	client HTTPClient
}

func NewHTTPClientProxy(client HTTPClient) *HTTPClientProxy {
	return &HTTPClientProxy{client: client}
}

func (customHttpClient HTTPClientProxy) Get(url string) (response *http.Response, err error) {
	return customHttpClient.client.Get(url)
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
