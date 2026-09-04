package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/akimsavvin/gonet/v2/di"
)

func main() {
	// Create instances manually before container
	cacheRepo := NewCacheRepositoryImpl()
	dbRepo := NewDBRepositoryImpl()

	c, err := di.NewContainer(
		// Register first repository with WithService passing the instance
		di.WithService[UserRepository](cacheRepo),
		// Register second repository with WithValue, the only difference is that it is type-safe
		di.WithValue[UserRepository](dbRepo),
		// Factory automatically receives []UserRepository with all registered implementations
		di.WithFactory(NewUserService),
	)
	if err != nil {
		log.Fatalf("could not create the container: %s", err)
	}

	service := c.MustGetService[*UserService]()

	fmt.Println("\n=== Creating user ===")
	userID := service.CreateUser("Akim")

	fmt.Println("\n=== Fetching user (cache first, then DB) ===")
	user := service.GetUser(userID)

	fmt.Println("\n=== Fetching non-existent user ===")
	user2 := service.GetUser(rand.Int())

	fmt.Println("\n=== Results ===")
	fmt.Println(user.Name == "Akim") // true
	fmt.Println(user2 == nil)        // true
}
