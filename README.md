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

GoNet provides advanced DI tools. There are two ways of adding a new service:

- gonet.AddScoped
- gonet.AddSingleton

The difference is in the service lifecycle. **AddSingleton** registers a service which is created only once and used
every time its requested, while transient service is created newly for each injection.

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

2. Create a provider for the needed service (Multiplier)

```go
package utils

type Multiplier struct{}

func NewMultiplier() *Multiplier {
	return &Multiplier{}
}

func (m *Multiplier) Multiply(x, y int) int {
	return x * y
}
```

3. Now you can register services and the GoNet will panic if there are any wrong implementations

```go
package main

import (
	"github.com/akimsavvin/gonet"
	"services"
	"utils"
)

func main() {
	gonet.AddService[services.Multiplier](utils.NewMultiplier)
	gonet.AddService[services.UserService](services.NewUserService)
}
```

3. Now you can use you service

```go
package main

import (
	"github.com/akimsavvin/gonet"
	"services"
	"utils"
)

func main() {
	gonet.AddService[services.Multiplier](utils.NewMultiplier)
	gonet.AddService[services.UserService](services.NewUserService)

	service := gonet.GetService[services.UserService]()
	service.LogSquared(4) // 16
	service.LogSquared(8) // 64
}
```

## Routing

GoNet provides advanced routing system, based on controllers.

Difference between services and controllers is that AddController method does not have a required generic and a
controller struct must have a **Register(router gonet.Router)** method.

**gonet.Router** is an interface which has all the necessary methods to define a routing system, including Use and all
the http methods methods. You can also add a controller prefix (as in .NET and other frameworks) with *
*router.AddPrefix(prefix string)** method.

All of the controller handlers must return a **gonet.ActionResult**. To access them you must use a *
*gonet.ControllerBase** Methods, so you must inject it in your controller struct.

In gonet you can not return any errors in route handler, so you must handle all of your errors with the needed payload.
But if you still want to return an error, you can use gonet.ControllerBase **Err(statusCode int, err error)** method

The **Content-Type** header is also set automatically according to your payload argument

### Example

1. Declare a controller and define a constructor which has all required dependencies in its parameters

```go
package controllers

import "github.com/akimsavvin/gonet"

type UserController struct {
	gonet.ControllerBase
	svc *services.UserService
}

func NewUserController(svc *services.UserService) *UserController {
	return &UserController{
		svc: svc,
	}
}

func (c *UserController) Register(router gonet.Router) {
	router.SetPrefix("/users")
	router.POST("/", c.create)
}

func (c *UserController) create(ctx gonet.Context) gonet.ActionResult {
	return c.Created(gonet.H{
		"id": "some-id",
	})
}

```

2. Add your controller and start an http server

```go
package main

import (
	"controllers"
	"github.com/akimsavvin/gonet"
	"services"
	"utils"
)

func main() {
	/* ............... */
	/* adding services */

	gonet.AddController(controllers.NewUserController)
	// gonet.StartHttp(":3000")
	gonet.StartHttp()
}
```



