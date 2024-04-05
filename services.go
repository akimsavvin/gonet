package gonet

import (
	"fmt"
	"reflect"
)

type Scope int

const (
	ScopeSingleton          Scope = iota + 1
	ScopeRequest_DO_NOT_USE       = iota + 1
	ScopeTransient                = iota + 1
)

type service struct {
	instance    *reflect.Value
	constructor reflect.Value
	typ         reflect.Type
	lifetime    Scope
}

func (s *service) getInstance() reflect.Value {
	if s.lifetime == ScopeSingleton && s.instance != nil {
		return *s.instance
	}

	depsCount := s.constructor.Type().NumIn()
	deps := make([]reflect.Value, depsCount)

	for i := 0; i < depsCount; i++ {
		deps[i] = getServiceUnsafe(s.constructor.Type().In(i))
	}

	instance := s.constructor.Call(deps)[0]

	if s.lifetime == ScopeSingleton {
		s.instance = &instance
	}

	return instance
}

var services = make(map[reflect.Type]*service)

func AddService[TService any](constructor any, lifetime Scope) {
	constructorType := reflect.TypeOf(constructor)

	if constructorType.Kind() != reflect.Func {
		panic(fmt.Sprintf("constructor %v is not a function", constructor))
	}

	if constructorType.NumOut() != 1 {
		panic(fmt.Sprintf("%v constructor must have only one return argument", constructor))
	}

	serviceType := reflect.TypeOf((*TService)(nil)).Elem()

	svc := &service{
		instance:    nil,
		constructor: reflect.ValueOf(constructor),
		typ:         serviceType,
		lifetime:    lifetime,
	}

	services[serviceType] = svc
}

func getServiceUnsafe(serviceType reflect.Type) reflect.Value {
	if serviceType != nil {
		if svc, ok := services[serviceType]; ok {
			return svc.getInstance()
		}
	}

	for _, v := range services {
		instance := v.getInstance()

		if ok := instance.Type().AssignableTo(serviceType); ok {
			return instance
		}
	}

	panic(fmt.Sprintf("failed to resolve service %v", serviceType))
}

func GetService[TService any]() (rservice TService) {
	serviceType := reflect.TypeOf((*TService)(nil)).Elem()

	if serviceType != nil {
		svc := services[serviceType]
		instance, ok := svc.getInstance().Interface().(TService)

		if !ok {
			panic(fmt.Sprintf("gotten service %v type does not satisfy interface %v", svc, serviceType))
		}

		return instance
	}

	panic("no service found for requested type")
}
