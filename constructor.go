package gonet

import (
	"fmt"
	"reflect"
)

func validateConstructor(constructor any) {
	constructorType := reflect.TypeOf(constructor)

	if constructorType.Kind() != reflect.Func {
		panic(fmt.Sprintf("constructor %v is not a function", constructor))
	}

	if constructorType.NumOut() != 1 {
		panic(fmt.Sprintf("%v constructor must have only one return argument", constructor))
	}
}
