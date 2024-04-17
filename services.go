// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package gonet

import (
	"fmt"
)

func AddSingleton[TService any](constructor any) {
	addService[TService](singleton, constructor)
}

func AddScoped_DO_NOT_USE[TService any](constructor any) {
	addService[TService](scoped, constructor)
}

func AddTransient[TService any](constructor any) {
	addService[TService](transient, constructor)
}

func addService[TService any](lifetime Lifetime, constructor any) {
	validateConstructor(constructor)
	addProvider[TService](servicePr, lifetime, constructor)
}

func GetService[TService any]() TService {
	p := getProvider[TService]()
	instance, ok := p.getInstance().Interface().(TService)

	if p.typ != servicePr || !ok {
		panic(fmt.Sprintf("no service found for type: %v", p.valTyp))
	}

	return instance
}
