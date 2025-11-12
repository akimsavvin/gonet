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
	repositories []UserRepository
}

func NewUserService(repositories []UserRepository) *UserService {
	return &UserService{
		repositories: repositories,
	}
}

func (service *UserService) GetUser(id int) *User {
	// Try all repositories in order
	for _, repo := range service.repositories {
		user, ok := repo.GetUser(id)
		if ok {
			return &user
		}
	}
	return nil
}

func (service *UserService) CreateUser(name string) int {
	user := User{
		ID:   rand.Int(),
		Name: name,
	}

	// Store in all repositories
	for _, repo := range service.repositories {
		repo.InsertUser(user)
	}

	return user.ID
}
