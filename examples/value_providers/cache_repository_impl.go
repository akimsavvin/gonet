package main

import "fmt"

type CacheRepositoryImpl struct {
	users map[int]User
}

var _ UserRepository = (*CacheRepositoryImpl)(nil)

func NewCacheRepositoryImpl() *CacheRepositoryImpl {
	fmt.Println("Creating CacheRepositoryImpl")
	return &CacheRepositoryImpl{
		users: make(map[int]User),
	}
}

func (cr *CacheRepositoryImpl) GetUser(id int) (User, bool) {
	user, ok := cr.users[id]
	if ok {
		fmt.Printf("Cache HIT for user ID %d\n", id)
	} else {
		fmt.Printf("Cache MISS for user ID %d\n", id)
	}
	return user, ok
}

func (cr *CacheRepositoryImpl) InsertUser(user User) {
	fmt.Printf("Caching user: %s (ID: %d)\n", user.Name, user.ID)
	cr.users[user.ID] = user
}
