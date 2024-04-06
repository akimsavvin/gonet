package gonet

import (
	"fmt"
)

func AddService[TService any](constructor any, lifetime Lifetime) {
	validateConstructor(constructor)
	addProvider[TService](constructor, servicePr, lifetime)
}

func GetService[TService any]() TService {
	p := getProvider[TService]()
	instance, ok := p.getInstance().Interface().(TService)

	if p.kind != servicePr || !ok {
		panic(fmt.Sprintf("no service found for type: %v", p.typ))
	}

	return instance
}
