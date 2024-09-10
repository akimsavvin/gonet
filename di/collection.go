// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

import (
	"github.com/akimsavvin/gonet/generic"
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

type (
	// serviceFactory is a service factory function description
	serviceFactory struct {
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
)

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

	// Implementation is the optional service implementation instance
	Instance *reflect.Value

	// Factory is service serviceFactory
	Factory *serviceFactory
}

// serviceCollection contains a list of service descriptors
type serviceCollection struct {
	// Descriptors is a slice of service descriptors added to the collection
	Descriptors []*serviceDescriptor
	// mx is a mutex to protect the Descriptors
	mx sync.RWMutex
}

// newServiceCollection creates a new serviceCollection
func newServiceCollection() *serviceCollection {
	return &serviceCollection{}
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
	coll.Descriptors = append(coll.Descriptors, sd)
}

// descriptors returns a slice of added service descriptors
func (coll *serviceCollection) descriptors() []*serviceDescriptor {
	coll.mx.RLock()
	defer coll.mx.RUnlock()

	descriptors := make([]*serviceDescriptor, len(coll.Descriptors))
	copy(descriptors, coll.Descriptors)

	return descriptors
}

// serviceCollectionInstance the instance of the service collection created with init function
var serviceCollectionInstance atomic.Pointer[serviceCollection]

// GetServiceCollection returns an instance of ServiceCollection
func GetServiceCollection() ServiceCollection {
	return serviceCollectionInstance.Load()
}

// init creates a new default service collection instance
func init() {
	serviceCollectionInstance.Store(newServiceCollection())
}

// AddService adds a new singleton service to the default service collection with the provided factory or instance
func AddService[T any](factoryOrValue any) {
	GetServiceCollection().AddService(generic.TypeOf[T](), factoryOrValue)
}

// AddKeyedService adds a new keyed singleton service to the default service collection with the provided factory or instance
func AddKeyedService[T any](key string, factoryOrValue any) {
	GetServiceCollection().AddKeyedService(generic.TypeOf[T](), key, factoryOrValue)
}
