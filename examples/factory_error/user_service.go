package main

import "math/rand"

type User struct {
	ID   int
	Name string
}

type UserRepository interface {
	InsertUser(User)
	GetUser(id int) (User, bool)
}

type UserService struct {
	repository UserRepository
}

func NewUserService(repository UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (service *UserService) GetUser(id int) *User {
	user, ok := service.repository.GetUser(id)
	if !ok {
		return nil
	}
	return &user
}

func (service *UserService) CreateUser(name string) int {
	user := User{
		ID:   rand.Int(),
		Name: name,
	}

	service.repository.InsertUser(user)
	return user.ID
}

