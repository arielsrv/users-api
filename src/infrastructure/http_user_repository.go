package infrastructure

import (
	"github.com/users-api/src/domain"
	"strconv"
)

type HttpUserRepository struct {
	client HttpClient
}

func NewHttpUserRepository(client HttpClient) *HttpUserRepository {
	return &HttpUserRepository{
		client: client,
	}
}

func (repository HttpUserRepository) GetUser(userId int) (*domain.User, error) {
	url := "https://gorest.co.in/public/v2/users/" + strconv.Itoa(userId)
	user, err := Client[domain.User]{client: repository.client}.Get(url)
	if err != nil {
		return &user, err
	}
	return &user, nil
}

func (repository HttpUserRepository) GetUsers() ([]domain.User, error) {
	url := "https://gorest.co.in/public/v2/users/"
	users, err := Client[[]domain.User]{client: repository.client}.Get(url)
	if err != nil {
		return users, err
	}
	return users, nil
}
