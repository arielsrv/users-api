package infrastructure

import (
	"fmt"
	"github.com/users-api/src/domain"
)

type HttpUserRepository struct {
	client  HttpClient
	baseUrl string
}

func NewHttpUserRepository(client HttpClient) *HttpUserRepository {
	return &HttpUserRepository{
		client:  client,
		baseUrl: "https://gorest.co.in/public/v2",
	}
}

func (repository HttpUserRepository) GetUser(userId int) (*domain.User, error) {
	url := fmt.Sprintf("%s/users/%d", repository.baseUrl, userId)
	user, err := Client[domain.User]{client: repository.client}.Get(url)
	return &user, err
}

func (repository HttpUserRepository) GetUsers() ([]domain.User, error) {
	url := fmt.Sprintf("%s/users", repository.baseUrl)
	users, err := Client[[]domain.User]{client: repository.client}.Get(url)
	return users, err
}
