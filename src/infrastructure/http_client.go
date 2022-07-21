package infrastructure

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http"
)

type HttpClient interface {
	Get(url string) (response *http.Response, err error)
}

type HttpClientProxy struct {
	client HttpClient
}

func NewHttpClientProxy(client HttpClient) *HttpClientProxy {
	return &HttpClientProxy{client: client}
}

func (customHttpClient HttpClientProxy) Get(url string) (response *http.Response, err error) {
	return customHttpClient.client.Get(url)
}

type Client[T any] struct {
	client HttpClient
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
