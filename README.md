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

3. Now you can register services with required lifetime (_see lifetimes later_) and the GoNet will panic if any service does not implement provided type

```go
package main

import (
	"github.com/akimsavvin/gonet/di"
	"services"
	"util"
)

func main() {
	di.AddTransient[services.Multiplier](util.NewMultiplier)
	di.AddTransient[services.UserService](services.NewUserService)
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
	di.AddTransient[services.Multiplier](util.NewMultiplier)
	di.AddTransient[services.UserService](services.NewUserService)

	service := di.GetService[services.UserService]()
	service.LogSquared(4) // 16
	service.LogSquared(8) // 64
}
```

### Lifetimes

You can add your services with different lifetimes\
Service lifetime determines when the service will be created and how long will it live

#### Singleton
Services with this lifetime are only created once and every time it's requested it takes an existing instance. You can singletons as follows:

```go
package main

import (
	"abstract"
	"fmt"
	"github.com/akimsavvin/gonet/di"
	"lib"
)

type MyService struct {
	// Some external interface
	otherServ abstract.MyInterface
}

func (s *MyService) Do() {
	s.otherServ.Call()
}

func NewMyService(otherServ abstract.MyInterface) *MyService {
	fmt.Println("Created new MyService instance")
	return &MyService{
		otherServ: otherServ,
	}
}

func main() {
	// Some kind of interface implementation
	di.AddSingleton[abstract.MyInterface](lib.NewIntegration)
	di.AddSingleton[*MyService](NewMyService)

	// Creates a new MyService instance
	serv := di.GetService[*MyService]() // Created new MyService instance
	serv.Do()

	// Takes the existing instance
	serv2 := di.GetService[*MyService]() // ...
	serv2.Do()

	fmt.Println(serv == serv2) // true
}
```

#### Scoped
Services with this lifetime are similiar with singleton, but only inside of a scope. On the other hand, you can not get scoped services outside the scope. A scope is basically created for each request, but it can be created manually with di.NewScope()

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
	di.AddScoped[*MyService](NewMyService)
	
	// Create new scope
	scope := di.NewScope()

	// Panics
	di.GetService[*MyService]()

	// Creates a new MyService instance
	serv := di.GetScopedService[*MyService](scope) // Created new MyService instance
	serv.Do()

	// Takes the existing instance
	serv2 := di.GetScopedService[*MyService](scope)
	serv2.Do()

	// Create a second scope
	scope2 := di.NewScope()

	// The new service is created for this scope
	serv3 := di.GetScopedService[*MyService](scope2)
	serv3.Do()

	fmt.Println(serv == serv2) // true
	fmt.Println(serv == serv3) // false
}
```

#### Transient
Services with this lifetime are created on each request

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

	// Creates a new MyService instance
	serv2 := di.GetService[*MyService]() // Created new MyService instance
	serv2.Do()

	fmt.Println(serv == serv2) // false
}
```

## Routing

GoNet provides advanced routing system, based on controllers.

Difference between services and controllers is that AddController method does not have a required generic and a
controller struct must have a **Builder(cb routig.ControllerBuilder)** method.

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

func (c *UserController) Register(cb gonet.ConrtollerBuilder) {
	cb.SetPrefix("/users")
	cb.POST("/", c.create)
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



