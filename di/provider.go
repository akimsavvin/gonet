// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

import (
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
	// descriptor is the accessor's serviceDescriptor
	descriptor *serviceDescriptor

	// once protects the instance from creating multiple times
	once sync.Once
	// instance is the accessor's service instance
	// nil if the service instance has not been requested yet
	instance *reflect.Value
}

// newServiceAccessor creates a new serviceAccessor
func newServiceAccessor(descriptor *serviceDescriptor) *serviceAccessor {
	return &serviceAccessor{
		descriptor: descriptor,
		instance:   descriptor.Instance,
	}
}

// createInstance creates an instance of the accessor descriptor service
func (accessor *serviceAccessor) createInstance(sp *serviceProvider) reflect.Value {
	deps := sp.resolveFactoryDeps(accessor.descriptor.Factory)
	instance, err := accessor.descriptor.Factory.Call(deps...)
	if err != nil {
		log.Panicf("[%v]: could not create a service instance due to error: %s\n",
			accessor.descriptor.Factory, err.Error())
	}

	return instance
}

// Instance returns the accessor service instance
func (accessor *serviceAccessor) Instance(sp *serviceProvider) reflect.Value {
	accessor.once.Do(func() {
		if accessor.instance == nil {
			instance := accessor.createInstance(sp)
			accessor.instance = &instance
		}
	})

	return *accessor.instance
}

type (
	// serviceAccessorsListItem is an item in serviceAccessorsList
	serviceAccessorsListItem struct {
		// Value is the serviceAccessor instance
		Value *serviceAccessor

		// Prev is the previous element in the list
		Prev *serviceAccessorsListItem
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
		Prev:  list.Tail,
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

	for current := list.Tail; current != nil; current = current.Prev {
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

// serviceProvider implements the ServiceProvider interface
type serviceProvider struct {
	// accessors is a map for service identifiers of service descriptors lists
	accessors serviceAccessors
}

// newServiceProvider creates a new serviceProvider
func newServiceProvider(serviceDescriptors []*serviceDescriptor) *serviceProvider {
	sp := new(serviceProvider)
	accessors := make(serviceAccessors)

	l := len(serviceDescriptors)
	extendedServiceDescriptors := make([]*serviceDescriptor, l, l+1)
	copy(extendedServiceDescriptors, serviceDescriptors)

	spValue := reflect.ValueOf(sp)
	spDescriptor := &serviceDescriptor{
		Type:               reflect.TypeFor[ServiceProvider](),
		ImplementationType: spValue.Type(),
		Instance:           &spValue,
	}

	extendedServiceDescriptors = append(extendedServiceDescriptors, spDescriptor)

	for _, descriptor := range extendedServiceDescriptors {
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

	sp.accessors = accessors
	return sp
}

// resolveFactoryDeps returns a slice of service dependencies for the provided factory
func (sp *serviceProvider) resolveFactoryDeps(factory *serviceFactory) []reflect.Value {
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

// GetServiceID gets an instance for the provided service identifier
func (sp *serviceProvider) GetServiceID(id serviceIdentifier) (reflect.Value, bool) {
	isSlice := id.Type.Kind() == reflect.Slice
	if isSlice {
		id.Type = id.Type.Elem()
	}

	accessors, ok := sp.accessors[id]
	if !ok {
		if isSlice {
			return reflect.Zero(reflect.SliceOf(id.Type)), false
		}
		return reflect.Zero(id.Type), false
	}

	if !isSlice {
		accessor := accessors.Last()
		return accessor.Instance(sp), true
	}

	res := reflect.MakeSlice(reflect.SliceOf(id.Type), accessors.Len, accessors.Len)
	for i, accessor := range accessors.Slice() {
		instance := accessor.Instance(sp)
		res.Index(i).Set(instance)
	}

	return res, true
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

// ServiceProviderInstance returns an instance of the default ServiceProvider
func ServiceProviderInstance() ServiceProvider {
	if instance := serviceProviderInstance.Load(); instance != nil {
		return instance
	}

	log.Panicln("[ServiceProvider]: provider was not built")
	return nil
}

// BuildServiceProvider builds a default ServiceProvider instance from the default ServiceCollection
func BuildServiceProvider() {
	serviceProviderInstance.Store(newServiceProvider(ServiceCollectionInstance().descriptors()))
}

// AssertService is used to assert returned value from the GetService method to provided generic type
func AssertService[T any](service reflect.Value, ok bool) (T, bool) {
	return AssertRequiredService[T](service), ok
}

// AssertRequiredService is used to assert returned value from the GetRequiredService method to provided generic type
func AssertRequiredService[T any](service reflect.Value) T {
	return service.Interface().(T)
}

// GetServiceSP returns an asserted service instance from the provided ServiceProvider instance
func GetServiceSP[T any](sp ServiceProvider) (T, bool) {
	return AssertService[T](sp.GetService(reflect.TypeFor[T]()))
}

// GetService returns an asserted service instance from the default ServiceProvider instance
func GetService[T any]() (T, bool) {
	return GetServiceSP[T](ServiceProviderInstance())
}

// GetKeyedServiceSP returns an asserted keyed service instance from the provided ServiceProvider instance
func GetKeyedServiceSP[T any](sp ServiceProvider, key string) (T, bool) {
	return AssertService[T](sp.GetKeyedService(reflect.TypeFor[T](), key))
}

// GetKeyedService returns an asserted keyed service instance from the default ServiceProvider instance
func GetKeyedService[T any](key string) (T, bool) {
	return GetKeyedServiceSP[T](ServiceProviderInstance(), key)
}

// GetRequiredServiceSP returns an asserted required service instance from the provided ServiceProvider instance
func GetRequiredServiceSP[T any](sp ServiceProvider) T {
	return AssertRequiredService[T](sp.GetRequiredService(reflect.TypeFor[T]()))
}

// GetRequiredService returns an asserted required service instance from the default ServiceProvider instance
func GetRequiredService[T any]() T {
	return GetRequiredServiceSP[T](ServiceProviderInstance())
}

// GetRequiredKeyedServiceSP returns an asserted keyed required service instance from the provided ServiceProvider instance
func GetRequiredKeyedServiceSP[T any](sp ServiceProvider, key string) T {
	return AssertRequiredService[T](sp.GetRequiredKeyedService(reflect.TypeFor[T](), key))
}

// GetRequiredKeyedService returns an asserted keyed required service instance from the default ServiceProvider instance
func GetRequiredKeyedService[T any](key string) T {
	return GetRequiredKeyedServiceSP[T](ServiceProviderInstance(), key)
}
