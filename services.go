package gonet

import (
	"fmt"
)

func AddSingleton[TService any](constructor any) {
	addService[TService](singleton, constructor)
}

func AddTransient[TService any](constructor any) {
	addService[TService](transient, constructor)
}

func addService[TService any](lifetime Lifetime, constructor any) {
	validateConstructor(constructor)
	addProvider[TService](lifetime, constructor, servicePr)
}

func GetService[TService any]() TService {
	p := getProvider[TService]()
	instance, ok := p.getInstance().Interface().(TService)

	if p.typ != servicePr || !ok {
		panic(fmt.Sprintf("no service found for type: %v", p.valTyp))
	}

	return instance
}
