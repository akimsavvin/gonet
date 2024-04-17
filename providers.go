// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package gonet

import (
	"fmt"
	"reflect"
)

type Lifetime int

const (
	singleton Lifetime = iota + 1
	scoped             = iota + 1
	transient          = iota + 1
)

type providerType int

const (
	servicePr    providerType = iota + 1
	controllerPr              = iota + 1
	valuePr                   = iota + 1
	configPr                  = iota + 1
)

func resolveValueDeps(val reflect.Value) []reflect.Value {
	typ := val.Type()

	if typ.Kind() != reflect.Func {
		panic(fmt.Sprintf("factory %s is not a function", typ))
	}

	count := typ.NumIn()
	deps := make([]reflect.Value, count)

	for i := 0; i < count; i++ {
		ptyp := typ.In(i)
		if ptyp.Kind() == reflect.Pointer {
			ptyp = ptyp.Elem()
		}

		p := getTypeProvider(ptyp)

		if p.typ == controllerPr {
			panic(fmt.Sprintf("no provider found for type %s", p.valTyp))
		}

		deps[i] = p.getInstance()
	}

	return deps
}

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

	deps := resolveValueDeps(p.constructor)
	instance := p.constructor.Call(deps)[0]

	if p.lifetime == singleton {
		p.value = &instance
	}

	return instance
}

type Providers map[reflect.Type]*provider

var providers Providers = make(map[reflect.Type]*provider)

// GetProvidersUnsafe returns a slice of internal providers
// Use it only when it's really needed
func GetProvidersUnsafe() Providers {
	return providers
}

func addTypeProvider(lifetime Lifetime, typ providerType, constructor any, valTyp reflect.Type) {
	if valTyp.Kind() == reflect.Pointer {
		valTyp = valTyp.Elem()
	}

	pvdr := &provider{
		typ:         typ,
		valTyp:      valTyp,
		value:       nil,
		constructor: reflect.ValueOf(constructor),
		lifetime:    lifetime,
	}

	providers[valTyp] = pvdr
}

func addProvider[TType any](typ providerType, lifetime Lifetime, constructor any) {
	valTyp := reflect.TypeOf((*TType)(nil)).Elem()
	constrTyp := reflect.TypeOf(constructor)
	retTyp := constrTyp.Out(0)

	if (valTyp.Kind() == reflect.Interface &&
		retTyp.Kind() != reflect.Pointer) ||
		valTyp.Kind() == reflect.Struct {
		retTyp = retTyp.Elem()
	}

	if !retTyp.AssignableTo(valTyp) {
		panic(fmt.Sprintf("constructor '%v' return type '%v' is not assignable to '%v'", constrTyp, retTyp, valTyp))
	}

	addTypeProvider(lifetime, typ, constructor, valTyp)
}

func getTypeProvider(typ reflect.Type) *provider {
	p, ok := providers[typ]
	if !ok || p == nil {
		panic(fmt.Sprintf("no provider found for type %s", typ))
	}

	return p
}

func getProvider[TProvider any]() *provider {
	typ := reflect.TypeOf((*TProvider)(nil)).Elem()
	return getTypeProvider(typ)
}
