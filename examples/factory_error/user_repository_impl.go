package main

import "errors"

type UserRepositoryImpl struct {
	users map[int]User
}

var _ UserRepository = (*UserRepositoryImpl)(nil)

func NewUserRepositoryImpl() (*UserRepositoryImpl, error) {
	return nil, errors.ErrUnsupported
}

func (ur *UserRepositoryImpl) GetUser(id int) (User, bool) {
	user, ok := ur.users[id]
	return user, ok
}

func (ur *UserRepositoryImpl) InsertUser(user User) {
	ur.users[user.ID] = user
}
