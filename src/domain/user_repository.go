package domain

type UserRepository interface {
	GetById(id int) (*User, error)
	GetAll() ([]User, error)
}
