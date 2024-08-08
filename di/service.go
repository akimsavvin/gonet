// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

// Singleton, Scoped, Transient

// Collection: Singleton, Service: Singleton => Get or create service instance
// Collection: Singleton, Service: Scoped => Panic
// Collection: Singleton, Service: Transient => Create service instance

// Collection: Scoped, Service: Singleton => Panic
// Collection: Scoped, Service: Scoped =>  Get or create service instance
// Collection: Scoped, Service: Transient => Panic

package di

import (
	"github.com/akimsavvin/gonet/generic"
	"log"
	"reflect"
)

// servDescriptor service description
type servDescriptor struct {
	typ    reflect.Type
	constr *constructor
	value  *reflect.Value
}

// newVal crete a new service instance or panics if err
func (sd servDescriptor) newVal(deps []reflect.Value) reflect.Value {
	vals := sd.constr.val.Call(deps)

	if sd.constr.hasErr && !vals[1].IsNil() {
		err := vals[1].Interface()
		log.Panicf("Could not create instance of %v due to error: %v\n", sd.typ, err)
	}

	return vals[0]
}

type servColl struct {
	sds map[reflect.Type]*servDescriptor
}

// newServColl creates a new service collection
func newServColl() *servColl {
	return &servColl{
		sds: make(map[reflect.Type]*servDescriptor),
	}
}

// addSD adds a new service descriptor to the collection
func (sc *servColl) addSD(sd *servDescriptor) {
	sc.sds[sd.typ] = sd
}

// getTypSD returns a service descriptor for the provided generic
func (sc *servColl) getTypSD(typ reflect.Type) *servDescriptor {
	return sc.sds[typ]
}

func (sc *servColl) resolveSDDeps(sd *servDescriptor) []reflect.Value {
	deps := make([]reflect.Value, sd.constr.depsCount)

	for i := range sd.constr.depsCount {
		typ := sd.constr.typ.In(i)
		val := sc.getTypVal(typ)
		deps[i] = val
	}

	return deps
}

// getTypSD returns a service descriptor for the provided generic
func (sc *servColl) getTypVal(typ reflect.Type) reflect.Value {
	sd := sc.getTypSD(typ)

	if sd == nil {
		log.Panicf("Could not find service descriptor for generic %v\n", typ)
	}

	if sd.value != nil {
		return *sd.value
	}

	deps := sc.resolveSDDeps(sd)
	val := sd.newVal(deps)

	return val
}

// AddService adds a constructor for the provided type
func AddService[T any](constr any) {
	typ := generic.GetType[T]()
	AddServiceType(typ, reflect.ValueOf(constr))
}

// AddServiceType adds a constructor for the provided type
func AddServiceType(typ reflect.Type, constrVal reflect.Value) {
	mustValidateConstrVal(typ, constrVal)
	sd := &servDescriptor{
		typ:    typ,
		constr: newConstructor(constrVal),
	}

	defaultContainer.addSD(sd)
}

// GetService gets service for T
func GetService[T any]() T {
	typ := generic.GetType[T]()
	return GetServiceByType(typ).Interface().(T)
}

// GetServiceByType gets service by the provided type
func GetServiceByType(typ reflect.Type) reflect.Value {
	return defaultContainer.getTypVal(typ)
}
