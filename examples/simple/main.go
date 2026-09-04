package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/akimsavvin/gonet/v2/di"
)

func main() {
	c, err := di.NewContainer(
		di.WithService[UserRepository](NewUserRepositoryImpl),
		di.WithFactory(NewUserService),
	)
	if err != nil {
		log.Fatalf("could not create the container: %s", err)
	}

	service := c.MustGetService[*UserService]()

	userID := service.CreateUser("Akim")
	user := service.GetUser(userID)
	user2 := service.GetUser(rand.Int())

	fmt.Println(user.Name == "Akim") // true
	fmt.Println(user2 == nil)        // true
}
