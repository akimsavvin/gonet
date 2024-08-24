# GoNet framework
___
GoNet is the first Golang framework with generic-based dependency injection, advanced routing, gRPC, kafka and
websockets support and an extensive collection of other tools.

*GoNet is inspired by the .NET and NestJS frameworks, so it may look a bit similar to them*

## Getting started

### Prerequisites

- **[Go](https://go.dev/)**: version 1.22.0 or higher (we use it innovations).

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

GoNet provides advanced DI tools.

### Example

1. Declare a service and define a constructor which has all required dependencies in its parameters

```go
package services

import (
	"fmt"
)

type Multiplier interface {
	Multiply(x, y int) int
}

type UserService struct {
	m Multiplier
}

func NewUserService(m Multiplier) *UserService {
	return &UserService{
		m: m,
	}
}

func (s *UserService) LogSquared(x int) {
	fmt.Println(s.m.Multiply(x, x))
}
```

2. Create a provider for the needed service (_Multiplier in this case_)

```go
package util

type Multiplier struct{}

func NewMultiplier() *Multiplier {
	return &Multiplier{}
}

func (m *Multiplier) Multiply(x, y int) int {
	return x * y
}
```

3. Now you can register services and the GoNet will panic if any service does not implement provided type

```go
package main

import (
	"github.com/akimsavvin/gonet/di"
	"services"
	"util"
)

func main() {
	di.AddService[services.Multiplier](util.NewMultiplier)
	di.AddService[services.UserService](services.NewUserService)
}
```

3. Now you can get your service as follows

```go
package main

import (
	"github.com/akimsavvin/gonet/di"
	"services"
	"util"
)

func main() {
	di.AddService[*services.Multiplier](util.NewMultiplier)
	di.AddService[*services.UserService](services.NewUserService)

	service := di.GetService[*services.UserService]()
	service.LogSquared(4) // 16
	service.LogSquared(8) // 64
}
```

```go
package main

import (
	"fmt"
	"github.com/akimsavvin/gonet/di"
)

type MyService struct{}

func (s *MyService) Do() {...}

func NewMyService() *MyService {
	fmt.Println("Created new MyService instance")
	return new(MyService)
}

func main() {
	di.AddTransient[*MyService](NewMyService)

	// Creates a new MyService instance
	serv := di.GetService[*MyService]() // Created new MyService instance
	serv.Do()
	serv2 := di.GetService[*MyService]() // Returns an existing instance
	serv2.Do()

	fmt.Println(serv == serv2) // false
}
```
