package main

import (
	"errors"
	"fmt"

	"github.com/akimsavvin/gonet/v2/di"
)

func main() {
	c := di.NewContainer(
		di.WithService[UserRepository](NewUserRepositoryImpl),
		di.WithFactory(NewUserService),
	)

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
