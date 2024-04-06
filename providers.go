package gonet

import (
	"fmt"
	"reflect"
)

type Lifetime int

const (
	singleton          Lifetime = iota + 1
	Request_DO_NOT_USE          = iota + 1
	transient                   = iota + 1
)

type providerType int

const (
	servicePr    providerType = iota + 1
	controllerPr              = iota + 1
	valuePr                   = iota + 1
	configPr                  = iota + 1
)

type provider struct {
	typ providerType
	// valTyp: all providers have
	valTyp reflect.Type
	// value: singleton providers have
	value       *reflect.Value
	constructor reflect.Value
	lifetime    Lifetime
}

func (p *provider) getInstance() reflect.Value {
	if p.lifetime == singleton && p.value != nil {
		return *p.value
	}

	depsCount := p.constructor.Type().NumIn()
	deps := make([]reflect.Value, depsCount)

	for i := 0; i < depsCount; i++ {
		depP := getTypeProvider(p.constructor.Type().In(i))
		if depP.typ == controllerPr {
			panic(fmt.Sprintf("no provider found for type %s", depP.valTyp))
		}

		deps[i] = depP.getInstance()
	}

	instance := p.constructor.Call(deps)[0]

	if p.lifetime == singleton {
		p.value = &instance
	}

	return instance
}

var providers = make(map[reflect.Type]*provider)

func addTypeProvider(valTyp reflect.Type, lifetime Lifetime, typ providerType, constructor any) {
	pvdr := &provider{
		typ:         typ,
		valTyp:      valTyp,
		value:       nil,
		constructor: reflect.ValueOf(constructor),
		lifetime:    lifetime,
	}

	providers[valTyp] = pvdr
}

func addProvider[TType any](lifetime Lifetime, constructor any, typ providerType) {
	valTyp := reflect.TypeOf((*TType)(nil)).Elem()
	constrTyp := reflect.TypeOf(constructor)
	retTyp := constrTyp.Out(0)

	if !retTyp.AssignableTo(valTyp) {
		panic(fmt.Sprintf("constructor %v return type %v is not assignable to %v", constrTyp.Name(), retTyp, valTyp))
	}

	addTypeProvider(valTyp, lifetime, typ, constructor)
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
