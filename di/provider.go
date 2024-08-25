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

type SingletonServiceAccessor struct {
	Descriptor *ServiceDescriptor
	Instance   *reflect.Value
}

func NewSingletonServiceAccessor(descriptor *ServiceDescriptor) *SingletonServiceAccessor {
	return &SingletonServiceAccessor{
		Descriptor: descriptor,
	}
}

func (accessor *SingletonServiceAccessor) GetInstance(scope *ServiceScope) reflect.Value {
	if accessor.Instance != nil {
		return *accessor.Instance
	}

	if accessor.Descriptor.Instance != nil {
		accessor.Instance = accessor.Descriptor.Instance
		return *accessor.Instance
	}

	instance := scope.CreateServiceInstance(accessor.Descriptor)
	// TODO: save it in the root descriptor to avoid creating it in each scope
	accessor.Instance = &instance

	return instance
}

type ScopedServiceAccessor struct {
	Descriptor *ServiceDescriptor
}

func NewScopedServiceAccessor(descriptor *ServiceDescriptor) *ScopedServiceAccessor {
	return &ScopedServiceAccessor{
		Descriptor: descriptor,
	}
}

func (accessor *ScopedServiceAccessor) GetInstance(scope *ServiceScope) reflect.Value {
	// TODO: implement it
	panic("not implemented")
}

type TransientServiceAccessor struct {
	Descriptor *ServiceDescriptor
}

func NewTransientServiceAccessor(descriptor *ServiceDescriptor) *TransientServiceAccessor {
	return &TransientServiceAccessor{
		Descriptor: descriptor,
	}
}

func (accessor *TransientServiceAccessor) GetInstance(scope *ServiceScope) reflect.Value {
	return scope.CreateServiceInstance(accessor.Descriptor)
}

func NewServiceAccessor(descriptor *ServiceDescriptor) ServiceAccessor {
	switch descriptor.Lifetime {
	case LifetimeSingleton:
		return NewSingletonServiceAccessor(descriptor)
	case LifetimeScoped:
		return NewScopedServiceAccessor(descriptor)
	case LifetimeTransient:
		return NewTransientServiceAccessor(descriptor)
	default:
		log.Panicf("Unknown service lifetime: %d", descriptor.Lifetime)
	}

	return nil
}

type ServiceAccessor interface {
	GetInstance(scope *ServiceScope) reflect.Value
}

type (
	ServiceDependencyTreeNode struct {
		Type   reflect.Type
		Parent *ServiceDependencyTreeNode
		Deps   []*ServiceDependencyTreeNode
	}

	ServiceDependencyTree struct {
		Root *ServiceDependencyTree
	}
)

type (
	ServiceAccessorsListItem struct {
		Value ServiceAccessor
		Next  *ServiceAccessorsListItem
	}

	ServiceAccessorsList struct {
		Root *ServiceAccessorsListItem
		Len  int
	}
)

func NewServiceDescriptorsList() *ServiceAccessorsList {
	return new(ServiceAccessorsList)
}

func (list *ServiceAccessorsList) Append(accessor ServiceAccessor) {
	list.Root = &ServiceAccessorsListItem{
		Value: accessor,
		Next:  list.Root,
	}
	list.Len++
}

func (list *ServiceAccessorsList) First() ServiceAccessor {
	return list.Root.Value
}

func (list *ServiceAccessorsList) Slice() []ServiceAccessor {
	sl := make([]ServiceAccessor, 0, list.Len)

	current := list.Root
	for current != nil {
		sl = append(sl, current.Value)
		current = current.Next
	}

	return sl
}

type ServiceAccessors map[ServiceIdentifier]*ServiceAccessorsList

type ServiceIdentifier struct {
	Type reflect.Type
	Key  string
}

type ServiceScope struct {
	IsRoot       bool
	RootProvider *ServiceProvider

	// Accessors is a map for service identifiers of service descriptors lists
	Accessors ServiceAccessors
}

func NewServiceScope(isRoot bool, rootProvider *ServiceProvider, serviceDescriptors []*ServiceDescriptor) *ServiceScope {
	accessors := make(ServiceAccessors)

	for _, descriptor := range serviceDescriptors {
		id := ServiceIdentifier{
			Type: descriptor.Type,
			Key:  descriptor.Key,
		}

		if accessors[id] == nil {
			accessors[id] = NewServiceDescriptorsList()
		}

		accessor := NewServiceAccessor(descriptor)
		accessors[id].Append(accessor)
	}

	return &ServiceScope{
		IsRoot:       isRoot,
		RootProvider: rootProvider,
		Accessors:    accessors,
	}
}

func (scope *ServiceScope) ResolveFactoryDeps(factory *ServiceFactory) []reflect.Value {
	serviceDeps := make([]reflect.Value, factory.DepsCount)

	for i := 0; i < factory.DepsCount; i++ {
		depType := factory.Type.In(i)
		depID := ServiceIdentifier{
			Type: depType,
		}

		dep, ok := scope.GetService(depID)
		if !ok {
			log.Panicf("[%v]: no service found for factory dependency type \"%v\"\n", factory.Type, depType)
		}

		serviceDeps[i] = dep
	}

	return serviceDeps
}

func (scope *ServiceScope) CreateServiceInstance(descriptor *ServiceDescriptor) reflect.Value {
	deps := scope.ResolveFactoryDeps(descriptor.Factory)
	values := descriptor.Factory.Value.Call(deps)
	if descriptor.Factory.HasErr && !values[1].IsNil() {
		log.Panicf("[%v]: could not create a service instance due to error: %s\n",
			descriptor.ImplementationType, values[1].Interface().(error).Error())
	}

	return values[0]
}

func (scope *ServiceScope) GetService(id ServiceIdentifier) (reflect.Value, bool) {
	isSlice := id.Type.Kind() == reflect.Slice
	if isSlice {
		id.Type = id.Type.Elem()
	}

	accessors, ok := scope.Accessors[id]
	if !ok {
		return reflect.Value{}, false
	}

	if isSlice {
		res := reflect.MakeSlice(reflect.SliceOf(id.Type), accessors.Len, accessors.Len)
		sl := accessors.Slice()

		for i, accessor := range sl {
			instance := accessor.GetInstance(scope)
			res.Index(i).Set(instance)
		}

		return res, ok
	}

	accessor := accessors.First()
	return accessor.GetInstance(scope), true
}

type ServiceProvider struct {
	RootScope *ServiceScope
}

func NewServiceProvider(serviceDescriptors []*ServiceDescriptor) *ServiceProvider {
	sp := new(ServiceProvider)
	sp.RootScope = NewServiceScope(true, sp, serviceDescriptors)

	return sp
}

func (sp *ServiceProvider) GetService(typ reflect.Type) (reflect.Value, bool) {
	return sp.GetKeyedService(typ, "")
}

func (sp *ServiceProvider) GetRequiredService(typ reflect.Type) reflect.Value {
	service, ok := sp.GetService(typ)
	if !ok {
		log.Panicf("[%v]: no service not found for the requested type\n", typ)
	}

	return service
}

func (sp *ServiceProvider) GetKeyedService(typ reflect.Type, key string) (reflect.Value, bool) {
	id := ServiceIdentifier{typ, key}
	return sp.RootScope.GetService(id)
}

func (sp *ServiceProvider) GetRequiredKeyedService(typ reflect.Type, key string) reflect.Value {
	service, ok := sp.GetKeyedService(typ, key)
	if !ok {
		log.Panicf("[%v]: no service not found for type with key \"%s\"\n", typ, key)
	}

	return service
}

var serviceProviderInstance atomic.Pointer[ServiceProvider]

func GetServiceProvider() *ServiceProvider {
	return serviceProviderInstance.Load()
}

func Build() {
	serviceProviderInstance.Store(NewServiceProvider(
		MustGetServiceCollection().Descriptors))
}

func AssertService[T any](service reflect.Value, ok bool) (T, bool) {
	if !ok {
		return generic.Default[T](), false
	}

	return AssertRequiredService[T](service), true
}

func AssertRequiredService[T any](service reflect.Value) T {
	return service.Interface().(T)
}

func GetService[T any]() (T, bool) {
	service, ok := GetServiceProvider().GetService(generic.TypeOf[T]())
	return AssertService[T](service, ok)
}

func GetKeyedService[T any](key string) (T, bool) {
	service, ok := GetServiceProvider().GetKeyedService(generic.TypeOf[T](), key)
	return AssertService[T](service, ok)
}

func GetRequiredService[T any]() T {
	return AssertRequiredService[T](GetServiceProvider().GetRequiredService(generic.TypeOf[T]()))
}

func GetRequiredKeyedService[T any](key string) T {
	return AssertRequiredService[T](GetServiceProvider().GetRequiredKeyedService(generic.TypeOf[T](), key))
}

//func GetScopedService[T](Scope scope) (T, bool)                  {}
//func GetKeyedScopedService[T](key string, Scope scope) (T, bool) {}
//func GetRequiredScopedService[T](Scope scope) T                     {}
//func GetKeyedRequiredScopedService[T](key string, Scope scope) T    {}
