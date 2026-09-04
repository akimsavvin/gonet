package main

import (
	"errors"
	"fmt"
	"log"

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

	service, err := c.GetService[*UserService]()
	var dependencyError *di.DependencyError
	if err != nil && errors.As(err, &dependencyError) {
		fmt.Println(
			errors.Is(
				errors.Unwrap(dependencyError),
				errors.ErrUnsupported,
			),
		) // True
	}

	fmt.Println(service == nil) // True
}
