package channel

import (
	"github.com/akimsavvin/gonet/generic"
	"log"
	"reflect"
)

const defaultName = "__GoNet_Default"

func Add[T any](size ...int) {
	AddNamed[T](defaultName, size...)
}

func AddNamed[T any](name string, size ...int) {
	typ := generic.GetType[T]()
	chSize := 0
	if len(size) > 0 {
		chSize = size[0]
	}

	val := reflect.ValueOf(make(chan T, chSize))

	defaultContainer.add(name, typ, &val)
}

func Get[T any]() chan T {
	return GetNamed[T](defaultName)
}

func GetNamed[T any](name string) chan T {
	typ := generic.GetType[T]()
	val := defaultContainer.get(name, typ)
	if val == nil {
		log.Panicf("Could not find channel with name \"%s\" and type %v\n", name, typ)
	}

	return val.Interface().(chan T)
}
