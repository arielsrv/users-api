package application

import (
	"github.com/users-api/src/domain"
)

type IUserService interface {
	GetUser(id int) (*UserDto, error)
	GetAll() ([]UserDto, error)
}
type UserService struct {
	userRepository domain.UserRepository
}

func NewUserService(userRepository domain.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (service UserService) GetUser(id int) (*UserDto, error) {
	user, err := service.userRepository.GetUser(id)
	if err != nil {
		return nil, err
	}
	userDto := UserDto{ID: user.ID, Name: user.Name, Email: user.Email}
	return &userDto, err
}

func (service UserService) GetAll() ([]UserDto, error) {
	users, err := service.userRepository.GetAll()
	if err != nil {
		return nil, err
	}

	var usersDto = make([]UserDto, len(users))
	for i, user := range users {
		usersDto[i] = UserDto{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		}
	}
	return usersDto, nil
}
