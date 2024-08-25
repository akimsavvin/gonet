// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

import (
	"github.com/akimsavvin/gonet/generic"
	"log"
	"reflect"
	"sync/atomic"
)

// ServiceLifetime lifetime of a service
type ServiceLifetime int

const (
	// LifetimeTransient is created every time it is requested
	LifetimeTransient ServiceLifetime = iota

	// LifetimeScoped service is created once for a scope and re-used every time it is requested in the scope
	LifetimeScoped

	// LifetimeSingleton service is created only once and then re-used every time it is requested
	LifetimeSingleton
)

// ServiceFactory is a service factory function description
type ServiceFactory struct {
	Type       reflect.Type
	Value      reflect.Value
	DepsCount  int
	ReturnType reflect.Type
	HasErr     bool
}

// NewServiceFactory creates a new ServiceFactory for the provided factory function
func NewServiceFactory(factory any) *ServiceFactory {
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

	return &ServiceFactory{
		Type:       typ,
		Value:      val,
		DepsCount:  typ.NumIn(),
		ReturnType: typ.Out(0),
		HasErr:     numOut == 2,
	}
}

// ServiceDescriptor describes a service in the service collection
type ServiceDescriptor struct {
	Type               reflect.Type
	Key                string
	ImplementationType reflect.Type
	Instance           *reflect.Value
	Factory            *ServiceFactory
	Lifetime           ServiceLifetime
}

// ServiceCollection contains a list of service descriptors
type ServiceCollection struct {
	Descriptors []*ServiceDescriptor
}

// NewServiceCollection creates a new ServiceCollection
func NewServiceCollection() *ServiceCollection {
	return &ServiceCollection{
		Descriptors: make([]*ServiceDescriptor, 0),
	}
}

// AddService adds a new service to the service collection with the provided lifetime
// Either factory or instance must be provided if there are both, the factory will be ignored
// instance can be provided only for singleton services
func (coll *ServiceCollection) AddService(typ reflect.Type, factory, instance any, lifetime ServiceLifetime) {
	coll.AddKeyedService(typ, "", factory, instance, lifetime)
}

// AddKeyedService adds a new keyed service to the service collection with the provided lifetime
// Either factory or instance must be provided if there are both, the factory will be ignored
// instance can be provided only for singleton services
func (coll *ServiceCollection) AddKeyedService(typ reflect.Type, key string, factory, instance any, lifetime ServiceLifetime) {
	sd := &ServiceDescriptor{
		Type:     typ,
		Key:      key,
		Lifetime: lifetime,
	}

	switch {
	case instance != nil:
		if lifetime != LifetimeSingleton {
			log.Panicf("[%v, %t]: non singleton service can not be initialized with instance\n", typ, instance)
		}

		instVal := reflect.ValueOf(instance)
		instTyp := instVal.Type()

		if !instTyp.AssignableTo(typ) {
			log.Panicf("[%v, %t]: service instance must be assignable to service type\n", typ, instance)
		}

		sd.Instance = &instVal
		sd.ImplementationType = instTyp
	case factory != nil:
		f := NewServiceFactory(factory)

		if !f.ReturnType.AssignableTo(typ) {
			log.Panicf("[%v, %v]: service factory return type must be assignable to service type\n",
				typ, reflect.TypeOf(factory))
		}

		sd.Factory = f
		sd.ImplementationType = f.ReturnType
	default:
		log.Panicf("[%v]: no factory or instance provided\n", typ)
	}

	coll.Descriptors = append(coll.Descriptors, sd)
}

// AddSingleton adds a new singleton service to the service collection with the provided factory or instance
func (coll *ServiceCollection) AddSingleton(typ reflect.Type, factoryOrInstance any) {
	coll.AddKeyedSingleton(typ, "", factoryOrInstance)
}

// AddKeyedSingleton adds a new keyed singleton service to the service collection with the provided factory or instance
func (coll *ServiceCollection) AddKeyedSingleton(typ reflect.Type, key string, factoryOrInstance any) {
	if reflect.TypeOf(factoryOrInstance).Kind() == reflect.Func {
		coll.AddKeyedService(typ, key, factoryOrInstance, nil, LifetimeSingleton)
	} else {
		coll.AddKeyedService(typ, key, nil, factoryOrInstance, LifetimeSingleton)
	}
}

// AddScoped adds a new scoped service to the service collection with the provided factory
func (coll *ServiceCollection) AddScoped(typ reflect.Type, factory any) {
	coll.AddService(typ, factory, nil, LifetimeScoped)
}

// AddKeyedScoped adds a new keyed scoped service to the service collection with the provided factory
func (coll *ServiceCollection) AddKeyedScoped(typ reflect.Type, key string, factory any) {
	coll.AddKeyedService(typ, key, factory, nil, LifetimeScoped)
}

// AddTransient adds a new transient service to the service collection with the provided factory
func (coll *ServiceCollection) AddTransient(typ reflect.Type, factory any) {
	coll.AddService(typ, factory, nil, LifetimeTransient)
}

// AddKeyedTransient adds a new keyed transient service to the service collection with the provided factory
func (coll *ServiceCollection) AddKeyedTransient(typ reflect.Type, key string, factory any) {
	coll.AddKeyedService(typ, key, factory, nil, LifetimeTransient)
}

// serviceCollectionInstance the instance of the service collection created with Init function
var serviceCollectionInstance atomic.Pointer[ServiceCollection]

// MustGetServiceCollection returns the default instance of the service collection
// panics if the service collection is not initialized
func MustGetServiceCollection() *ServiceCollection {
	serviceCollection := serviceCollectionInstance.Load()
	if serviceCollection == nil {
		log.Panicln("service collection instance is not initialized")
	}

	return serviceCollection
}

// init creates a new default service collection instance
func init() {
	serviceCollectionInstance.Store(NewServiceCollection())
}

// AddSingleton adds a new singleton service to the default service collection with the provided factory or instance
func AddSingleton[T any](factoryOrValue any) {
	MustGetServiceCollection().AddSingleton(generic.TypeOf[T](), factoryOrValue)
}

// AddKeyedSingleton adds a new keyed singleton service to the default service collection with the provided factory or instance
func AddKeyedSingleton[T any](key string, factoryOrValue any) {
	MustGetServiceCollection().AddKeyedSingleton(generic.TypeOf[T](), key, factoryOrValue)
}

// AddScoped adds a new scoped service to the default service collection with the provided factory
func AddScoped[T any](factory any) {
	MustGetServiceCollection().AddScoped(generic.TypeOf[T](), factory)
}

// AddKeyedScoped adds a new keyed scoped service to the default service collection with the provided factory
func AddKeyedScoped[T any](key string, factory any) {
	MustGetServiceCollection().AddKeyedScoped(generic.TypeOf[T](), key, factory)
}

// AddTransient adds a new transient service to the default service collection with the provided factory
func AddTransient[T any](factory any) {
	MustGetServiceCollection().AddTransient(generic.TypeOf[T](), factory)
}

// AddKeyedTransient adds a new keyed transient service to the default service collection with the provided factory
func AddKeyedTransient[T any](key string, factory any) {
	MustGetServiceCollection().AddKeyedTransient(generic.TypeOf[T](), key, factory)
}
