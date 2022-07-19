package infrastructure

import (
	"github.com/essentialkaos/librato/v10"
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http"
	"os"
	"time"
)

type WrappedClient interface {
	Get(url string) (response *http.Response, err error)
}

type HttpClient interface {
	Get(url string) (response *Response, err error)
}

type CustomClient struct {
	client WrappedClient
	name   string
}

func NewCustomClient(client WrappedClient, name string) *CustomClient {
	return &CustomClient{
		client: client,
		name:   name,
	}
}

type Response struct {
	Raw  http.Response
	Data []byte
}

func (c CustomClient) Get(url string) (r *Response, err error) {
	start := time.Now()
	response, err := c.client.Get(url)
	elapsed := time.Since(start)

	c.recordElapsedTime(elapsed)

	if err != nil {
		err = fiber.NewError(http.StatusInternalServerError, err.Error())
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		err = fiber.NewError(response.StatusCode, "Couldn't retreive. ")
		return nil, err
	}

	body, _ := io.ReadAll(response.Body)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	return &Response{
		Raw:  *response,
		Data: body,
	}, err
}

func (c CustomClient) recordElapsedTime(elapsed time.Duration) {
	libratoToken, env1 := os.LookupEnv("LIBRATO_TOKEN")
	libratoUser, env2 := os.LookupEnv("LIBRATO_USER")
	if env1 && env2 {
		librato.Mail = libratoUser
		librato.Token = libratoToken
		_ = librato.AddMetric(
			librato.Gauge{
				Name:  "httpClient:" + c.name,
				Value: elapsed,
			},
		)
	}
}
