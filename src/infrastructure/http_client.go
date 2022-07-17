package infrastructure

import (
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http"
)

type HttpClient interface {
	Get(url string) (response *Response, err error)
}

type CustomClient struct {
	client http.Client
}

func NewCustomClient(client http.Client) *CustomClient {
	return &CustomClient{client: client}
}

type Response struct {
	Raw  http.Response
	Data []byte
}

func (c CustomClient) Get(url string) (r *Response, err error) {
	response, err := c.client.Get(url)

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
