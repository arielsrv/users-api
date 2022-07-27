package infrastructure

import (
	"fmt"
	"github.com/users-api/src/domain"
)

type HTTPUserRepository struct {
	client  HTTPClient
	baseURL string
}

func NewHTTPUserRepository(client HTTPClient) *HTTPUserRepository {
	return &HTTPUserRepository{
		client:  client,
		baseURL: "https://gorest.co.in/public/v2",
	}
}

func (repository HTTPUserRepository) GetById(userID int) (*domain.User, error) {
	url := fmt.Sprintf("%s/users/%d", repository.baseURL, userID)
	user, err := Client[domain.User]{repository.client}.Get(url)
	return &user, err
}

func (repository HTTPUserRepository) GetAll() ([]domain.User, error) {
	url := fmt.Sprintf("%s/users", repository.baseURL)
	users, err := Client[[]domain.User]{repository.client}.Get(url)
	return users, err
}
