<p>
  <em><b>Gonet</b> is a <a href="https://dotnet.microsoft.com">.NET</a> inspired <b>DI framework</b>. Designed to <b>ease</b> things up for <b>fast</b> development with <a href="https://docs.gofiber.io/#zero-allocation"><b>zero memory allocation</b></a> and <b>performance</b> in mind.</em>
</p>

## ⚙️ Installation

Gonet requires **Go version `1.23` or higher** to run. If you need to install or upgrade Go, visit the [official Go download page](https://go.dev/dl/). To start setting up your project, create a new directory for your project and navigate into it. Then, initialize your project with Go modules by executing the following command in your terminal:

```bash
go mod init github.com/you/repo
```

To learn more about Go modules and how they work, you can check out the [Using Go Modules](https://go.dev/blog/using-go-modules) blog post.

After setting up your project, you can install Gonet with the `go get` command:

```bash
go get -u github.com/akimsavvin/gonet/v2
```

This command fetches the Gonet package and adds it to your project's dependencies, allowing you to start building your applications with Gonet.

## ⚡️ Quickstart

Getting started with Gonet is easy. Here's a basic example to create a simple di container that contains a usual dependencies of modern application. This example demonstrates initializing a new Gonet container, setting up a dependencies, and getting the services.

```go title="Example"
package main

import (
	"github.com/akimsavvin/gonet/v2/di"
	"github.com/app/config"
	"github.com/app/rest"
	"github.com/app/storage"
	"github.com/app/usecase"
)

func main() {
	// Initialize a new Gonet container
	c := di.NewContainer(
		// Define a typed service
		// usecase.UserRepo - interface used in the usecase.NewUserService function
		// storage.NewUserRepo - function creating a service which implements the usecase.UserRepo interface
		di.WithService[usecase.UserRepo](storage.NewUserRepo),

		// Define a value-typed service
		// di.WithFactory just adds the service for the value type,
		// so it's the same as di.WithService[*config.Config](config.New())
		di.WithValue(config.New()),

		// Define a factory-typed service
		// usecase.UserRepo - interface defining the repository
		// di.WithFactory just adds the service for the function return type
		// so it's the same as di.WithService[*usecase.UserService](usecase.NewUserService)
		di.WithFactory(usecase.NewUserService),

		// rest.NewUserController requires *usecase.UserService dependency
		di.WithService[rest.Controller](rest.NewUserController),
	)

	// Get all services of type rest.Controller
	controllers := di.MustGetService[[]rest.Controller](c)
	for _, contr := range controllers {
		// Do something with the received services
		contr.Init(...)
	}
}

```

This simple container is easy to set up and run. It introduces the core concepts of Gonet: container initialization, dependencies definition, and receiving the services.
