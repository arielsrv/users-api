package infrastructure

import "github.com/users-api/src/domain"

type HttpUserRepository struct {
}

func NewUserRepository() *HttpUserRepository {
	return &HttpUserRepository{}
}

func (repository HttpUserRepository) GetUser(userId int) domain.User {
	user := domain.User{Id: userId, Name: "Steve Jobs", Email: "stevejobs@apple.com"}
	return user
}
