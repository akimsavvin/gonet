// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

import (
	"github.com/akimsavvin/gonet/generic"
	"log"
	"reflect"
)

// Singleton, Scoped, Transient

// Collection: Singleton, Service: Singleton => Get or create service instance
// Collection: Singleton, Service: Scoped => Panic
// Collection: Singleton, Service: Transient => Create service instance

// Collection: Scoped, Service: Singleton => Panic
// Collection: Scoped, Service: Scoped =>  Get or create service instance
// Collection: Scoped, Service: Transient => Panic

type servDescriptor struct {
	typ      reflect.Type
	constr   *constructor
	lifetime Lifetime
	// value for singleton and scoped services
	value *reflect.Value
}

func (sd servDescriptor) newVal(deps []reflect.Value) reflect.Value {
	vals := sd.constr.val.Call(deps)

	if sd.constr.hasErr && vals[1].IsNil() {
		err := vals[1].Interface()
		log.Panicf("Could not create instance of %v due to error: %v\n", sd.typ, err)
	}

	return vals[0]
}

type servColl struct {
	sds      map[reflect.Type]*servDescriptor
	lifetime Lifetime
}

func newServColl(lifetime Lifetime) *servColl {
	if lifetime == LifetimeTransient {
		log.Panicf("Can not create transient service collection\n")
	}

	return &servColl{
		sds:      make(map[reflect.Type]*servDescriptor),
		lifetime: lifetime,
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

	if sd.lifetime == LifetimeScoped && sc.lifetime != LifetimeScoped {
		log.Panicf("Can not create a scoped service instance for generic %v outside of a scope\n", typ)
	}

	if sc.lifetime == LifetimeScoped && sd.lifetime != LifetimeScoped {
		log.Panicf("Can not create a non-scoped service instance for generic %v in a scoped service collection\n", typ)
	}

	if sc.lifetime == LifetimeSingleton &&
		sd.lifetime == LifetimeSingleton &&
		sd.value != nil {
		return *sd.value
	}

	deps := sc.resolveSDDeps(sd)
	val := sd.newVal(deps)

	if sd.lifetime == LifetimeSingleton || sd.lifetime == LifetimeScoped {
		sd.value = &val
	}

	return val
}

// getScopedColl returns new service collection with only scoped service descriptors
func (sc *servColl) getScopedColl() *servColl {
	var scopedColl = newServColl(LifetimeScoped)

	for _, sd := range sc.sds {
		if sd.lifetime == LifetimeScoped {
			scopedColl.addSD(&*sd)
		}
	}

	return scopedColl
}

func AddTransient[T any](constr any) {
	AddService[T](constr, LifetimeTransient)
}

func AddScoped[T any](constr any) {
	AddService[T](constr, LifetimeScoped)
}

func AddSingleton[T any](constr any) {
	AddService[T](constr, LifetimeSingleton)
}

func AddService[T any](constr any, lifetime Lifetime) {
	typ := generic.GetType[T]()
	AddServiceType(typ, reflect.ValueOf(constr), lifetime)
}

func AddServiceType(typ reflect.Type, constrVal reflect.Value, lifetime Lifetime) {
	mustValidateConstrVal(typ, constrVal)
	sd := &servDescriptor{
		typ:      typ,
		constr:   newConstructor(constrVal),
		lifetime: lifetime,
	}

	defaultContainer.addSD(sd)
}

func GetService[T any]() T {
	typ := generic.GetType[T]()
	return GetServiceByType(typ).Interface().(T)
}

func GetServiceByType(typ reflect.Type) reflect.Value {
	return defaultContainer.getTypVal(typ)
}

func GetScopedService[T any](scope *Scope) T {
	typ := generic.GetType[T]()
	return GetScopedServiceByType(scope, typ).Interface().(T)
}

func GetScopedServiceByType(scope *Scope, typ reflect.Type) reflect.Value {
	return scope.getTypVal(typ)
}
