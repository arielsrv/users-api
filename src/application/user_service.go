package application

import (
	"github.com/users-api/src/domain"
)

type IUserService interface {
	GetUser(id int) (*UserDto, error)
	GetUsers() ([]UserDto, error)
}
type UserService struct {
	userRepository domain.UserRepository
}

func NewUserService(userRepository domain.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (service UserService) GetUser(id int) (*UserDto, error) {
	user, err := service.userRepository.GetUser(id)
	userDto := UserDto{Id: user.Id, Name: user.Name, Email: user.Email}
	return &userDto, err
}

func (service UserService) GetUsers() ([]UserDto, error) {
	users, _ := service.userRepository.GetUsers()
	var usersDto = make([]UserDto, len(users))

	for i, user := range users {
		var userDto UserDto
		userDto.Id = user.Id
		userDto.Name = user.Name
		userDto.Email = user.Email
		usersDto[i] = userDto
	}

	return usersDto, nil
}
