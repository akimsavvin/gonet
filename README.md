# GoNet framework

___

- 🔥 GoNet is the first full-fledged framework made for Golang!
- ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
- 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

___

## Getting started

### Prerequisites

- **[Go](https://go.dev/)**: version 1.18.0 or higher.

### Getting GoNet

With [Go module](https://github.com/golang/go/wiki/Modules) support, simply add the following import

```go
import gonet "github.com/akimsavvin/gonet/v2"
```

to your code, and then `go [build|run|test]` will automatically fetch the necessary dependencies.

Otherwise, run the following Go command to install the `gonet` package:

```sh
$ go get -u github.com/akimsavvin/gonet/v2
```

## Dependency injection

GoNet provides advanced tools to deal with dependency injection.

### Example

1. Declare a service and define a constructor which has all required dependencies in its parameters

```go
package greeter

import "fmt"

type UserRepo interface {
	GetNameByID(id int) string
}

type Greeter struct {
	repo UserRepo
}

func NewGreeter(repo UserRepo) *Greeter {
	return &Greeter{
		repo: repo,
	}
}

func (g *Greeter) Greet(id int) {
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
	"myproject/greeter"
	"myproject/storage"
)

func main() {
	di.AddService[greeter.UserRepo](storage.NewUserRepo)
	di.AddService[*greeter.Greeter](greeter.NewGreeter)

	g := di.GetRequiredService[*greeter.Greeter]()
	g.Greet(17) // stdout: Hello, Akim!
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

## Environment

GoNet provides the tools to interact with the ENVIRONMENT variable with the **env** package.

### Examples

```go
package main

import (
	"fmt"
	"github.com/akimsavvin/gonet/env"
	"os"
)

func main() {
	os.Setenv("ENVIRONMENT", "Staging")

	curEnv, ok := env.Current()
	fmt.Println(ok) // true
	fmt.Println(curEnv == env.Staging) // true

	os.Clearenv()

	curEnv, ok = env.Current()
	fmt.Println(ok) // false
	fmt.Println(curEnv == "") // true

	curEnv = env.CurrentOrDefault()
	fmt.Println(curEnv == env.Development) // true
	
	os.Setenv("ENVIRONMENT", "Production")

	curEnv = env.CurrentOrDefault()
	fmt.Println(curEnv == env.Production) // true
}
```

## Graceful shutdown

GoNet provides the tools to deal with graceful shutdown with **graceful** package.

### Example

```go
package main

import (
	"app"
	"context"
	"fmt"
	"github.com/akimsavvin/gonet/graceful"
)

func main() {
	// ctx will be cancelled on os.Interrupt or os.Kill
	ctx, cancel := graceful.Context()
	defer cancel()

	go app.Start(ctx, ...)

	// current goroutine will be blocked and wait until the application is stopped
	graceful.WaitShutdown()
}
```