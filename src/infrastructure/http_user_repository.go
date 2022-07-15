package infrastructure

import (
	"encoding/json"
	"github.com/users-api/src/domain"
	"io"
	"log"
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

func (repository HttpUserRepository) GetUser(userId int) *domain.User {
	url := "https://gorest.co.in/public/v2/users/" + strconv.Itoa(userId)
	response, err := repository.client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode == http.StatusOK {
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(response.Body)
		body, _ := io.ReadAll(response.Body)
		user := domain.User{}
		_ = json.Unmarshal(body, &user)
		return &user
	}
	return nil
}

func (repository HttpUserRepository) GetUsers() []domain.User {
	url := "https://gorest.co.in/public/v2/users/"
	response, err := repository.client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode == http.StatusOK {
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(response.Body)
		body, _ := io.ReadAll(response.Body)
		var users []domain.User
		_ = json.Unmarshal(body, &users)
		return users
	}
	return nil
}
