package gonet

import (
	"fmt"
	"reflect"
)

func validateConstructor(constructor any) {
	typ := reflect.TypeOf(constructor)

	if typ.Kind() != reflect.Func {
		panic(fmt.Sprintf("constructor %v is not a function", constructor))
	}

	if typ.NumOut() != 1 {
		panic(fmt.Sprintf("%v constructor must have only one return argument", constructor))
	}
}
