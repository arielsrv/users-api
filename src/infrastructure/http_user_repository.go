package infrastructure

import (
	"encoding/json"
	"github.com/users-api/src/domain"
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
		return nil, err
	}
	user := domain.User{}
	_ = json.Unmarshal(response.Data, &user)
	return &user, nil
}

func (repository HttpUserRepository) GetUsers() ([]domain.User, error) {
	url := "https://gorest.co.in/public/v2/users/"
	response, err := repository.client.Get(url)
	if err != nil {
		return nil, err
	}
	var users []domain.User
	_ = json.Unmarshal(response.Data, &users)
	return users, nil
}
