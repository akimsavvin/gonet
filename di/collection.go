// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

import (
	"log"
	"reflect"
	"sync"
	"sync/atomic"
)

// ServiceCollection is a collection of service descriptors
type ServiceCollection interface {
	// AddService adds a service descriptor to the collection for the provided type
	// the second argument must be either a factory creating a service instance or an instance itself
	AddService(typ reflect.Type, factoryOrInstance any)

	// AddKeyedService adds a service descriptor to the collection for the provided type and key
	// the third argument must be either a factory creating a service instance or an instance itself
	AddKeyedService(typ reflect.Type, key string, factoryOrInstance any)

	// descriptors returns a slice of added service descriptors
	descriptors() []*serviceDescriptor
}

// serviceDescriptor describes a service in the service collection
type serviceDescriptor struct {
	// Type is the service type
	Type reflect.Type

	// HasKey is true if the service is keyed
	HasKey bool

	// Key is the service key.
	// Empty if the service is not keyed
	Key string

	// ImplementationType is the service implementation type
	ImplementationType reflect.Type

	// Instance is the optional service implementation instance
	Instance *reflect.Value

	// Factory is service serviceFactory
	Factory *serviceFactory
}

// serviceCollection implements the ServiceCollection interface
// contains a list of service descriptors
type serviceCollection struct {
	// mx is a mutex to protect the sds
	mx sync.RWMutex
	// sds is a slice of service descriptors added to the collection
	sds []*serviceDescriptor
}

// newServiceCollection creates a new serviceCollection
func newServiceCollection() *serviceCollection {
	return &serviceCollection{
		sds: make([]*serviceDescriptor, 0),
	}
}

// AddService adds a new service to the service collection
func (coll *serviceCollection) AddService(typ reflect.Type, factoryOrInstance any) {
	coll.AddServiceKey(typ, nil, factoryOrInstance)
}

// AddKeyedService adds a new keyed service to the service collection
func (coll *serviceCollection) AddKeyedService(typ reflect.Type, key string, factoryOrInstance any) {
	coll.AddServiceKey(typ, &key, factoryOrInstance)
}

// AddServiceKey adds service to the service collection with the provided key
func (coll *serviceCollection) AddServiceKey(typ reflect.Type, key *string, factoryOrInstance any) {
	if reflect.TypeOf(factoryOrInstance).Kind() == reflect.Func {
		coll.AddServiceFactory(typ, key, factoryOrInstance)
	} else {
		coll.AddServiceInstance(typ, key, factoryOrInstance)
	}
}

// AddServiceFactory adds a new keyed service with a factory to the service collection
func (coll *serviceCollection) AddServiceFactory(typ reflect.Type, key *string, factory any) {
	sd := &serviceDescriptor{
		Type:   typ,
		HasKey: key != nil,
	}

	if sd.HasKey {
		sd.Key = *key
	}

	f := newServiceFactory(factory)

	if !f.ReturnType.AssignableTo(typ) {
		log.Panicf("[%v, %v]: service factory return type must be assignable to service type\n",
			typ, reflect.TypeOf(factory))
	}

	sd.Factory = f
	sd.ImplementationType = f.ReturnType

	coll.AddDescriptor(sd)
}

// AddServiceInstance adds a new keyed service with an instance to the service collection
func (coll *serviceCollection) AddServiceInstance(typ reflect.Type, key *string, instance any) {
	sd := &serviceDescriptor{
		Type:   typ,
		HasKey: key != nil,
	}

	if sd.HasKey {
		sd.Key = *key
	}

	instVal := reflect.ValueOf(instance)
	instTyp := instVal.Type()

	if !instTyp.AssignableTo(typ) {
		log.Panicf("[%v, %t]: service instance must be assignable to service type\n", typ, instance)
	}

	sd.Instance = &instVal
	sd.ImplementationType = instTyp

	coll.AddDescriptor(sd)
}

// AddDescriptor adds a new service descriptor to the service collection
func (coll *serviceCollection) AddDescriptor(sd *serviceDescriptor) {
	coll.mx.Lock()
	defer coll.mx.Unlock()
	coll.sds = append(coll.sds, sd)
}

// descriptors returns a slice of added service descriptors
func (coll *serviceCollection) descriptors() []*serviceDescriptor {
	coll.mx.RLock()
	defer coll.mx.RUnlock()

	descriptors := make([]*serviceDescriptor, len(coll.sds))
	copy(descriptors, coll.sds)

	return descriptors
}

// serviceCollectionInstance the instance of the service collection created with init function
var serviceCollectionInstance atomic.Pointer[serviceCollection]

// ServiceCollectionInst returns an instance of the default ServiceCollection
func ServiceCollectionInst() ServiceCollection {
	return serviceCollectionInstance.Load()
}

// init creates a new default service collection instance
func init() {
	serviceCollectionInstance.Store(newServiceCollection())
}

// AddService adds a new singleton service to the default service collection with the provided factory or instance
func AddService[T any](factoryOrValue any) {
	ServiceCollectionInst().AddService(reflect.TypeFor[T](), factoryOrValue)
}

// AddValue adds a new singleton value the default service collection with the provided value
// Same as the AddService[T](value), but typed
func AddValue[T any](value T) {
	ServiceCollectionInst().AddService(reflect.TypeFor[T](), value)
}

// AddKeyedService adds a new keyed singleton service to the default service collection with the provided factory or instance
func AddKeyedService[T any](key string, factoryOrValue any) {
	ServiceCollectionInst().AddKeyedService(reflect.TypeFor[T](), key, factoryOrValue)
}

// AddKeyedValue adds a new keyed value to the default service collection with the provided value
// Same as the AddKeyedService[T](value), but typed
func AddKeyedValue[T any](key string, value T) {
	ServiceCollectionInst().AddKeyedService(reflect.TypeFor[T](), key, value)
}
