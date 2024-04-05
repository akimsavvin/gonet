package gonet

import (
	"fmt"
	"reflect"
)

type Lifetime int

const (
	LifetimeSingleton          Lifetime = iota + 1
	LifetimeRequest_DO_NOT_USE          = iota + 1
	LifetimeTransient                   = iota + 1
)

type providerKind int

const (
	servicePr    providerKind = iota + 1
	controllerPr              = iota + 1
	valuePr                   = iota + 1
	configPr                  = iota + 1
)

type provider struct {
	kind providerKind
	// typ: all providers have
	typ reflect.Type
	// value: singleton providers have
	value       *reflect.Value
	constructor reflect.Value
	lifetime    Lifetime
}

func (p *provider) getInstance() reflect.Value {
	if p.lifetime == LifetimeSingleton && p.value != nil {
		return *p.value
	}

	depsCount := p.constructor.Type().NumIn()
	deps := make([]reflect.Value, depsCount)

	for i := 0; i < depsCount; i++ {
		depP := getTypeProvider(p.constructor.Type().In(i))
		if depP.kind == controllerPr {
			panic(fmt.Sprintf("no provider found for type %s", depP.typ))
		}

		deps[i] = depP.getInstance()
	}

	instance := p.constructor.Call(deps)[0]

	if p.lifetime == LifetimeSingleton {
		p.value = &instance
	}

	return instance
}

var providers = make(map[reflect.Type]*provider)

func addTypeProvider(typ reflect.Type, constructor any, kind providerKind, lifetime Lifetime) {
	pvdr := &provider{
		kind:        kind,
		typ:         typ,
		value:       nil,
		constructor: reflect.ValueOf(constructor),
		lifetime:    lifetime,
	}

	providers[typ] = pvdr
}

func addProvider[TType any](constructor any, kind providerKind, lifetime Lifetime) {
	typ := reflect.TypeOf((*TType)(nil)).Elem()
	addTypeProvider(typ, constructor, kind, lifetime)
}

func getTypeProvider(typ reflect.Type) *provider {
	p := providers[typ]
	if p == nil {
		panic(fmt.Sprintf("not provider found for type %s", typ))
	}

	return p
}

func getProvider[TProvider any]() *provider {
	typ := reflect.TypeOf((*TProvider)(nil)).Elem()
	return getTypeProvider(typ)
}
