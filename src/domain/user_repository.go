package domain

type UserRepository interface {
	GetByID(id int) (*User, error)
	GetAll() ([]User, error)
}
