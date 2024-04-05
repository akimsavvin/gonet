package gonet

import (
	"fmt"
	"reflect"
)

func AddService[TService any](constructor any, lifetime Lifetime) {
	validateConstructor[TService](constructor)

	returnType := reflect.TypeOf(constructor).Out(0)
	typ := reflect.TypeOf((*TService)(nil)).Elem()

	if !typ.Out(0).Implements(typ) {
		panic(fmt.Sprintf("constructor return type %v does not implement %v", returnType, typ))
	}

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
