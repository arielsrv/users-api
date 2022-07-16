package infrastructure

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/users-api/src/domain"
	"io"
	"net/http"
	"strconv"
)

type HttpUserRepository struct {
	client HttpClient
}

func NewHttpUserRepository(httpClient HttpClient) *HttpUserRepository {
	return &HttpUserRepository{
		client: httpClient,
	}
}

func (repository HttpUserRepository) GetUser(userId int) (*domain.User, error) {
	url := "https://gorest.co.in/public/v2/users/" + strconv.Itoa(userId)
	response, err := repository.client.Get(url)

	if err != nil {
		err = fiber.NewError(http.StatusInternalServerError, err.Error())
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		err = fiber.NewError(response.StatusCode, fmt.Sprintf("Couldn't retreive user with id %d not found. ", userId))
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	body, _ := io.ReadAll(response.Body)
	user := domain.User{}
	_ = json.Unmarshal(body, &user)
	return &user, nil
}

func (repository HttpUserRepository) GetUsers() ([]domain.User, error) {
	url := "https://gorest.co.in/public/v2/users/"
	response, err := repository.client.Get(url)
	if err != nil {
		err = fiber.NewError(http.StatusInternalServerError, err.Error())
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		err = fiber.NewError(response.StatusCode, "Couldn't retreive. ")
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	body, _ := io.ReadAll(response.Body)
	var users []domain.User
	_ = json.Unmarshal(body, &users)
	return users, nil
}
