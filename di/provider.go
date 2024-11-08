// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

import (
	"github.com/akimsavvin/gonet/v2/generic"
	"log"
	"reflect"
	"slices"
	"sync"
	"sync/atomic"
)

// ServiceProvider a provider to access the services
type ServiceProvider interface {
	// GetService returns the service instance by the type
	// the second return argument indicates whether a service for the requested type was found
	GetService(typ reflect.Type) (reflect.Value, bool)

	// GetKeyedService returns the service instance by the type and key
	// the second return argument indicates whether a service for the requested type was found
	GetKeyedService(typ reflect.Type, key string) (reflect.Value, bool)

	// GetRequiredService returns the service instance by the type
	// panics if no service found for the requested type
	GetRequiredService(typ reflect.Type) reflect.Value

	// GetRequiredKeyedService returns the service instance by the type and key
	// panics if no service found for the requested type
	GetRequiredKeyedService(typ reflect.Type, key string) reflect.Value
}

// serviceAccessor is a struct used to get the service instance
type serviceAccessor struct {
	// Descriptor is the accessor's serviceDescriptor
	Descriptor *serviceDescriptor

	// mx protects the service instance
	mx sync.RWMutex
	// Instance is the accessor's service instance
	// nil if the service instance has not been requested yet
	Instance *reflect.Value
}

// newServiceAccessor creates a new serviceAccessor
func newServiceAccessor(descriptor *serviceDescriptor) *serviceAccessor {
	return &serviceAccessor{
		Descriptor: descriptor,
		Instance:   descriptor.Instance,
	}
}

// GetInstance returns the accessor service instance
func (accessor *serviceAccessor) GetInstance(sp *serviceProvider) reflect.Value {
	accessor.mx.Lock()
	instance := accessor.Instance
	accessor.mx.Unlock()

	if instance != nil {
		return *accessor.Instance
	}

	accessor.mx.Lock()
	defer accessor.mx.Unlock()

	// if the goroutine was interrupted before the write lock was acquired,
	// need to check if the instance has been already set
	if accessor.Instance != nil {
		return *accessor.Instance
	}

	newInstance := sp.GetServiceInstance(accessor.Descriptor)
	accessor.Instance = &newInstance

	return newInstance
}

type (
	// serviceAccessorsListItem is an item in serviceAccessorsList
	serviceAccessorsListItem struct {
		// Value is the serviceAccessor instance
		Value *serviceAccessor

		// Value is the next serviceAccessorsListItem in the list
		Next *serviceAccessorsListItem
	}

	// serviceAccessorsList is a singly-linked list of service accessors
	serviceAccessorsList struct {
		// Tail is the last element of the list
		Tail *serviceAccessorsListItem

		// Len is the number of elements in the list
		Len int
	}
)

// newServiceAccessorsList creates a new serviceAccessorsList
func newServiceAccessorsList(accessors ...*serviceAccessor) *serviceAccessorsList {
	list := new(serviceAccessorsList)

	for _, accessor := range accessors {
		list.Append(accessor)
	}

	return list
}

// Append adds an element to the end of the list
func (list *serviceAccessorsList) Append(accessor *serviceAccessor) {
	list.Tail = &serviceAccessorsListItem{
		Value: accessor,
		Next:  list.Tail,
	}
	list.Len++
}

// Last returns the last element of the list
func (list *serviceAccessorsList) Last() *serviceAccessor {
	return list.Tail.Value
}

// Slice converts the list to a slice
func (list *serviceAccessorsList) Slice() []*serviceAccessor {
	sl := make([]*serviceAccessor, 0, list.Len)

	for current := list.Tail; current != nil; current = current.Next {
		sl = append(sl, current.Value)
	}

	slices.Reverse(sl)

	return sl
}

// serviceAccessors a map with a list of accessors for the key ID
type serviceAccessors map[serviceIdentifier]*serviceAccessorsList

// serviceIdentifier stores the service type and key
type serviceIdentifier struct {
	// Type is the service type
	Type reflect.Type

	// Key is the service key.
	// Empty if the service is not keyed
	Key string

	// HasKey is true if the service is keyed
	HasKey bool
}

// ResolveFactoryDeps returns a slice of service dependencies for the provided factory
func (sp *serviceProvider) ResolveFactoryDeps(factory *serviceFactory) []reflect.Value {
	serviceDeps := make([]reflect.Value, factory.DepsCount)

	for i := 0; i < factory.DepsCount; i++ {
		depType := factory.Type.In(i)
		depID := serviceIdentifier{
			Type: depType,
		}

		dep, ok := sp.GetServiceID(depID)
		if !ok {
			log.Panicf("[%v]: no service found for factory dependency type \"%v\"\n", factory.Type, depType)
		}

		serviceDeps[i] = dep
	}

	return serviceDeps
}

// GetServiceInstance gets a descriptor's existing service instance or creates a new one
func (sp *serviceProvider) GetServiceInstance(descriptor *serviceDescriptor) reflect.Value {
	deps := sp.ResolveFactoryDeps(descriptor.Factory)
	values := descriptor.Factory.Value.Call(deps)
	if descriptor.Factory.HasErr && !values[1].IsNil() {
		log.Panicf("[%v]: could not create a service instance due to error: %s\n",
			descriptor.ImplementationType, values[1].Interface().(error).Error())
	}

	return values[0]
}

// serviceProvider implements the ServiceProvider interface
type serviceProvider struct {
	// accessors is a map for service identifiers of service descriptors lists
	accessors serviceAccessors
}

// newServiceProvider creates a new serviceProvider
func newServiceProvider(serviceDescriptors []*serviceDescriptor) *serviceProvider {
	accessors := make(serviceAccessors)

	for _, descriptor := range serviceDescriptors {
		id := serviceIdentifier{
			Type:   descriptor.Type,
			Key:    descriptor.Key,
			HasKey: descriptor.HasKey,
		}

		if _, ok := accessors[id]; !ok {
			accessors[id] = newServiceAccessorsList()
		}

		accessor := newServiceAccessor(descriptor)
		accessors[id].Append(accessor)
	}

	return &serviceProvider{
		accessors: accessors,
	}
}

// GetServiceID gets an instance for the provided service identifier
func (sp *serviceProvider) GetServiceID(id serviceIdentifier) (reflect.Value, bool) {
	isSlice := id.Type.Kind() == reflect.Slice
	if isSlice {
		id.Type = id.Type.Elem()
	}

	accessors, ok := sp.accessors[id]
	if !ok {
		return reflect.Value{}, false
	}

	if isSlice {
		res := reflect.MakeSlice(reflect.SliceOf(id.Type), accessors.Len, accessors.Len)
		for i, accessor := range accessors.Slice() {
			instance := accessor.GetInstance(sp)
			res.Index(i).Set(instance)
		}

		return res, true
	}

	accessor := accessors.Last()
	return accessor.GetInstance(sp), true
}

// GetService gets an instance for the provided service type
func (sp *serviceProvider) GetService(typ reflect.Type) (reflect.Value, bool) {
	id := serviceIdentifier{typ, "", false}
	return sp.GetServiceID(id)
}

// GetKeyedService gets an instance for the provided service type and key
func (sp *serviceProvider) GetKeyedService(typ reflect.Type, key string) (reflect.Value, bool) {
	id := serviceIdentifier{typ, key, true}
	return sp.GetServiceID(id)
}

// GetRequiredService gets an instance for the provided service type and panics if no service found
func (sp *serviceProvider) GetRequiredService(typ reflect.Type) reflect.Value {
	service, ok := sp.GetService(typ)
	if !ok {
		log.Panicf("[%v]: no service not found for the requested type\n", typ)
	}

	return service
}

// GetRequiredKeyedService gets an instance for the provided service type and key and panics if no service found
func (sp *serviceProvider) GetRequiredKeyedService(typ reflect.Type, key string) reflect.Value {
	service, ok := sp.GetKeyedService(typ, key)
	if !ok {
		log.Panicf("[%v]: no service not found for type with key \"%s\"\n", typ, key)
	}

	return service
}

// serviceProviderInstance a default serviceProvider instance
var serviceProviderInstance atomic.Pointer[serviceProvider]

// GetServiceProvider returns an instance of ServiceProvider
func GetServiceProvider() ServiceProvider {
	if instance := serviceProviderInstance.Load(); instance != nil {
		return instance
	}

	log.Panicln("[ServiceProvider]: provider was not built")
	return nil
}

// Build builds a default serviceProvider instance from the current service collection
func Build() {
	serviceProviderInstance.Store(newServiceProvider(
		GetServiceCollection().descriptors()))
}

// AssertService is used to assert returned value from the GetService method to provided generic type
func AssertService[T any](service reflect.Value, ok bool) (T, bool) {
	if !ok {
		return generic.Default[T](), false
	}

	return AssertRequiredService[T](service), true
}

// AssertRequiredService is used to assert returned value from the GetRequiredService method to provided generic type
func AssertRequiredService[T any](service reflect.Value) T {
	return service.Interface().(T)
}

// GetServiceSP returns an asserted service instance from the provided ServiceProvider instance
func GetServiceSP[T any](sp ServiceProvider) (T, bool) {
	return AssertService[T](sp.GetService(generic.TypeOf[T]()))
}

// GetService returns an asserted service instance from the default ServiceProvider instance
func GetService[T any]() (T, bool) {
	return GetServiceSP[T](GetServiceProvider())
}

// GetKeyedServiceSP returns an asserted keyed service instance from the provided ServiceProvider instance
func GetKeyedServiceSP[T any](sp ServiceProvider, key string) (T, bool) {
	return AssertService[T](sp.GetKeyedService(generic.TypeOf[T](), key))
}

// GetKeyedService returns an asserted keyed service instance from the default ServiceProvider instance
func GetKeyedService[T any](key string) (T, bool) {
	return GetKeyedServiceSP[T](GetServiceProvider(), key)
}

// GetRequiredServiceSP returns an asserted required service instance from the provided ServiceProvider instance
func GetRequiredServiceSP[T any](sp ServiceProvider) T {
	return AssertRequiredService[T](sp.GetRequiredService(generic.TypeOf[T]()))
}

// GetRequiredService returns an asserted required service instance from the default ServiceProvider instance
func GetRequiredService[T any]() T {
	return GetRequiredServiceSP[T](GetServiceProvider())
}

// GetRequiredKeyedServiceSP returns an asserted keyed required service instance from the provided ServiceProvider instance
func GetRequiredKeyedServiceSP[T any](sp ServiceProvider, key string) T {
	return AssertRequiredService[T](sp.GetRequiredKeyedService(generic.TypeOf[T](), key))
}

// GetRequiredKeyedService returns an asserted keyed required service instance from the default ServiceProvider instance
func GetRequiredKeyedService[T any](key string) T {
	return GetRequiredKeyedServiceSP[T](GetServiceProvider(), key)
}
