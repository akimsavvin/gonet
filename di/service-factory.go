// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

import (
	"github.com/akimsavvin/gonet/generic"
	"log"
	"reflect"
)

// serviceFactory is a service factory function description
type serviceFactory struct {
	// Type is factory type
	Type reflect.Type

	// Value is factory value
	Value reflect.Value

	// DepsCount is number of the factory dependencies
	DepsCount int

	// ReturnType is the return type of the factory function
	ReturnType reflect.Type

	// HasErr is true if the factory returns an error as the second return argument
	HasErr bool
}

// newServiceFactory creates a new serviceFactory for the provided factory function
func newServiceFactory(factory any) *serviceFactory {
	val := reflect.ValueOf(factory)
	typ := val.Type()

	if typ.Kind() != reflect.Func {
		log.Panicf("[%t]: service factory must be a function\n", factory)
	}

	numOut := typ.NumOut()
	switch numOut {
	case 0:
		log.Panicf("[%t]: service factory must return at least one value\n", factory)
	case 1:
	case 2:
		if typ.Out(1) != generic.TypeOf[error]() {
			log.Panicf("[%t]: second service factory return value must be an error\n", factory)
		}
	default:
		log.Panicf("[%t]: service factory returns too many values\n", factory)
	}

	return &serviceFactory{
		Type:       typ,
		Value:      val,
		DepsCount:  typ.NumIn(),
		ReturnType: typ.Out(0),
		HasErr:     numOut == 2,
	}
}
