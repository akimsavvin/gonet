# GoNet framework

___

- 🔥 GoNet is the first full-fledged framework made for Golang!
- ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
- 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

___

## Getting started

### Prerequisites

- **[Go](https://go.dev/)**: version 1.20.0 or higher.

### Getting GoNet

With [Go module](https://github.com/golang/go/wiki/Modules) support, simply add the following import

```go
import "github.com/akimsavvin/gonet"
```

to your code, and then `go [build|run|test]` will automatically fetch the necessary dependencies.

Otherwise, run the following Go command to install the `gonet` package:

```sh
$ go get -u github.com/akimsavvin/gonet
```

## Dependency injection

GoNet provides advanced tools to deal with dependency injection.

### Example

1. Declare a service and define a constructor which has all required dependencies in its parameters

```go
package greeter

import (
	"fmt"
)

type UserRepo interface {
	GetNameByID(id int) string
}

type Greeter struct {
	repo UserRepo
}

func New(repo UserRepo) *Greeter {
	return &Greeter{
		repo: repo,
	}
}

func (g *Greeter) GreetByID(id int) {
	fmt.Printf("Hello, %s!\n", g.repo.GetNameByID(id))
}
```

2. Create a service for the **UserRepo**

```go
package storage

type UserRepo struct {
	data map[int]string
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
       data: map[int]string{
		   17: "Akim",
       },
    }
}

func (repo *UserRepo) GetNameByID(id int) string {
	return repo.data[id]
}
```

3. Now you can add your services to the default collection

```go
package main

import (
	"github.com/akimsavvin/gonet/di"
	"greeter"
	"storage"
)

func main() {
	di.AddService[greeter.UserRepo](storage.NewUserRepo)
	di.AddService[*greeter.Greeter](greeter.New)
}
```

3. Now you must build the service provider and you can get your service as follows

```go
package main

import (
	"github.com/akimsavvin/gonet/di"
	"greeter"
	"storage"
)

func main() {
	di.AddService[greeter.UserRepo](storage.NewUserRepo)
	di.AddService[*greeter.Greeter](greeter.New)

	// Build function build the service provider instance (and check the services dependencies if the future),  
	// which is then used to get services
	di.Build()

	service := di.GetRequiredService[*greeter.Greeter]()
	service.Greet(17) // Hello, Akim!
}
```

