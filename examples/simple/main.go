package main

import (
	"fmt"
	"math/rand"

	"github.com/akimsavvin/gonet/v2/di"
)

func main() {
	c := di.NewContainer(
		di.WithService[UserRepository](NewUserRepositoryImpl),
		di.WithFactory(NewUserService),
	)

	service := di.MustGetService[*UserService](c)

	userID := service.CreateUser("Akim")
	user := service.GetUser(userID)
	user2 := service.GetUser(rand.Int())

	fmt.Println(user.Name == "Akim") // true
	fmt.Println(user2 == nil)        // true
}
