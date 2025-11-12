package main

import "fmt"

type DBRepositoryImpl struct {
	users map[int]User
}

var _ UserRepository = (*DBRepositoryImpl)(nil)

func NewDBRepositoryImpl() *DBRepositoryImpl {
	fmt.Println("Creating DBRepositoryImpl")
	return &DBRepositoryImpl{
		users: make(map[int]User),
	}
}

func (dr *DBRepositoryImpl) GetUser(id int) (User, bool) {
	user, ok := dr.users[id]
	if ok {
		fmt.Printf("DB: Found user ID %d\n", id)
	} else {
		fmt.Printf("DB: User ID %d not found\n", id)
	}
	return user, ok
}

func (dr *DBRepositoryImpl) InsertUser(user User) {
	fmt.Printf("DB: Saving user: %s (ID: %d)\n", user.Name, user.ID)
	dr.users[user.ID] = user
}
