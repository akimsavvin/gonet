// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

import (
	"log"
	"reflect"
)

type constructor struct {
	val       reflect.Value
	depsCount int
	typ       reflect.Type
	retTyp    reflect.Type
	hasErr    bool
}

// newConstructor creates a new constructor struct
// must be called after mustValidateConstrVal
func newConstructor(val reflect.Value) *constructor {
	typ := val.Type()
	return &constructor{
		val:       val,
		typ:       typ,
		depsCount: typ.NumIn(),
		retTyp:    typ.Out(0),
		hasErr:    typ.NumOut() == 2,
	}
}

func mustValidateConstrVal(typ reflect.Type, constrVal reflect.Value) {
	if constrVal.Kind() != reflect.Func {
		log.Panicf("%v constructor is not a function\n", typ)
	}

	numOut := constrVal.Type().NumOut()

	if numOut > 2 {
		log.Panicf("%v constructor has to many return values\n", typ)
	} else if numOut < 1 {
		log.Panicf("%v constructor has to few return values\n", typ)
	}

	errTyp := getGenericType[error]()
	if numOut == 2 && !constrVal.Type().Out(1).Implements(errTyp) {
		log.Panicf("%v constructor second return value is of type error\n", typ)
	}

	if constrVal.Type().Out(0).AssignableTo(typ) {
		log.Panicf("%v constructor return value is not assignable to this type\n", typ)
	}
}
