package application

import (
	"github.com/users-api/src/domain"
	"net/http"
)

type IUserService interface {
	GetByID(id int) (*UserDto, error)
	MultiGetByID(ids []int) ([]MultiGetDto[UserDto], error)
	GetAll() ([]UserDto, error)
}
type UserService struct {
	userRepository domain.UserRepository
}

func NewUserService(userRepository domain.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (service UserService) GetByID(id int) (*UserDto, error) {
	user, err := service.userRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	userDto := UserDto{ID: user.ID, Name: user.Name, Email: user.Email}
	return &userDto, err
}

func (service UserService) MultiGetByID(ids []int) ([]MultiGetDto[UserDto], error) {
	channel := make(chan *UserDto, len(ids))
	var result []MultiGetDto[UserDto]

	for _, id := range ids {
		user, _ := service.GetByID(id)
		channel <- user
	}
	close(channel)

	for userDto := range channel {
		result = append(result, MultiGetDto[UserDto]{
			Code: http.StatusOK,
			Body: *userDto,
		})
	}

	return result, nil
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
