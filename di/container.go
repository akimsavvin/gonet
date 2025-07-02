// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

import (
	"errors"
	"fmt"
	"log"
	"reflect"
)

var (
	// ErrServiceNotFound is the error returned when a requested service is not found is service provider
	ErrServiceNotFound = errors.New("di: requested service not found")
)

// DependencyError is a custom error type for dependency injection failures.
type DependencyError struct {
	// DependencyType is the type of the dependency that failed to be created
	DependencyType reflect.Type

	// RequestingType is the type of the service requesting the dependency
	RequestingType reflect.Type

	// Err is the underlying error that occurred
	Err error
}

// Error implements the error interface for DependencyError.
func (e *DependencyError) Error() string {
	return fmt.Sprintf(
		"di: failed to create dependency %q for service %q: %v",
		e.DependencyType,
		e.RequestingType,
		e.Err,
	)
}

// Unwrap returns the underlying error
func (e *DependencyError) Unwrap() error {
	return e.Err
}

// Option is a struct adding a new service to the Container
type Option interface {
	apply(*Container)
}

// serviceFactoryOption adds a new keyed service with a factory to the Container
type serviceFactoryOption struct {
	typ     reflect.Type
	key     *string
	factory any
}

// apply applies the Option
func (opt *serviceFactoryOption) apply(c *Container) {
	id := newServiceIdentifier(opt.typ, opt.key)

	f := newServiceFactory(opt.factory)
	if !f.ReturnType.AssignableTo(opt.typ) {
		log.Panicf("[%v, %v]: service factory return type must be assignable to service type\n",
			opt.typ, reflect.TypeOf(opt.factory))
	}

	accessor := newServiceAccessor(id, c, f, nil)
	c.appendAccessor(id, accessor)
}

// withServiceFactory returns a new instance of serviceFactoryOption
func withServiceFactory[T any](key *string, factory any) Option {
	return &serviceFactoryOption{
		typ:     reflect.TypeFor[T](),
		key:     key,
		factory: factory,
	}
}

// serviceInstanceOption adds a new keyed service with an instance to the Container
type serviceInstanceOption struct {
	typ      reflect.Type
	key      *string
	instance any
}

// apply applies the Option
func (opt *serviceInstanceOption) apply(c *Container) {
	id := newServiceIdentifier(opt.typ, opt.key)

	instVal := reflect.ValueOf(opt.instance)
	instTyp := instVal.Type()

	if !instTyp.AssignableTo(opt.typ) {
		log.Panicf("[%v, %t]: service instance must be assignable to service type\n", opt.typ, opt.instance)
	}

	accessor := newServiceAccessor(id, c, nil, &instVal)
	c.appendAccessor(id, accessor)
}

// withServiceInstance adds a new keyed service with an instance to the Container
func withServiceInstance[T any](key *string, instance any) Option {
	return &serviceInstanceOption{
		typ:      reflect.TypeFor[T](),
		key:      key,
		instance: instance,
	}
}

// withServiceKey adds service to the Container with the provided key
func withServiceKey[T any](key *string, factoryOrInstance any) Option {
	if reflect.TypeOf(factoryOrInstance).Kind() == reflect.Func {
		return withServiceFactory[T](key, factoryOrInstance)
	} else {
		return withServiceInstance[T](key, factoryOrInstance)
	}
}

// WithService adds a new service to the Container with the provided factory or instance
func WithService[T any](factoryOrInstance any) Option {
	return withServiceKey[T](nil, factoryOrInstance)
}

// WithKeyedService adds a new keyed service to the Container with the provided factory or instance
func WithKeyedService[T any](key string, factoryOrInstance any) Option {
	return withServiceKey[T](&key, factoryOrInstance)
}

// WithValue adds a new value to the Container with the provided value
// Same as the WithService[T](value), but typed
func WithValue[T any](value T) Option {
	return withServiceKey[T](nil, value)
}

// WithKeyedValue adds a new keyed value to the Container with the provided value
// Same as the WithKeyedService[T](key, value), but typed
func WithKeyedValue[T any](key string, value T) Option {
	return withServiceKey[T](&key, value)
}

// factoryOption adds a new service factory to the Container
type factoryOption struct {
	factory any
	key     *string
}

// apply applies the Option
func (opt *factoryOption) apply(c *Container) {
	f := newServiceFactory(opt.factory)
	id := newServiceIdentifier(f.ReturnType, opt.key)
	accessor := newServiceAccessor(id, c, f, nil)
	c.appendAccessor(id, accessor)
}

// WithFactory adds a new service factory to the Container
func WithFactory(factory any) Option {
	return &factoryOption{
		factory: factory,
	}
}

// WithKeyedFactory adds a new keyed service factory to the Container
func WithKeyedFactory(key string, factory any) Option {
	return &factoryOption{
		factory: factory,
		key:     &key,
	}
}

// Container is a service container
type Container struct {
	// accessors is a map for service identifiers of service descriptors lists
	accessors serviceAccessors
}

// NewContainer creates a new Container
func NewContainer(opts ...Option) *Container {
	c := &Container{
		accessors: make(serviceAccessors),
	}

	l := len(opts)
	extOpts := make([]Option, l, l+1)
	copy(extOpts, opts)

	// Extend options with the default accessors
	extOpts = append(
		extOpts,
		withServiceInstance[ServiceGetter](nil, c),
	)

	for _, opt := range extOpts {
		opt.apply(c)
	}

	return c
}

// appendAccessor appends a service accessor to the container to the provided id
// creates a new serviceAccessorsList if the id does not exist
func (c *Container) appendAccessor(id serviceIdentifier, accessor *serviceAccessor) {
	if l, ok := c.accessors[id]; ok {
		l.Append(accessor)
	} else {
		c.accessors[id] = newServiceAccessorsList(accessor)
	}
}

// resolveFactoryDeps returns a slice of service dependency for the provided factory
func (c *Container) resolveFactoryDeps(factory *serviceFactory) ([]reflect.Value, error) {
	serviceDeps := make([]reflect.Value, factory.DepsCount)

	for i := 0; i < factory.DepsCount; i++ {
		depType := factory.Type.In(i)
		depID := serviceIdentifier{
			Type: depType,
		}

		dep, err := c.getService(depID)
		if err != nil {
			return nil, &DependencyError{
				RequestingType: factory.ReturnType,
				DependencyType: depType,
				Err:            err,
			}
		}

		serviceDeps[i] = dep
	}

	return serviceDeps, nil
}

// getService gets a service instance for the provided service identifier
func (c *Container) getService(id serviceIdentifier) (reflect.Value, error) {
	isSlice := id.Type.Kind() == reflect.Slice
	if isSlice {
		id.Type = id.Type.Elem()
	}

	accessors, ok := c.accessors[id]
	if !ok {
		if isSlice {
			return reflect.Zero(reflect.SliceOf(id.Type)), ErrServiceNotFound
		}
		return reflect.Zero(id.Type), ErrServiceNotFound
	}

	if !isSlice {
		accessor := accessors.Last()
		return accessor.Instance()
	}

	slTyp := reflect.SliceOf(id.Type)
	res := reflect.MakeSlice(slTyp, accessors.Len(), accessors.Len())
	for i, accessor := range accessors.Iter() {
		instance, err := accessor.Instance()
		if err != nil {
			return reflect.Zero(slTyp), err
		}

		res.Index(i).Set(instance)
	}

	return res, nil
}

// getServiceKey returns an asserted service instance
// for the provided type and key
// from the provided ServiceGetter
func getServiceKey[T any](sc ServiceGetter, key *string) (T, error) {
	id := newServiceIdentifier(reflect.TypeFor[T](), key)
	service, err := sc.getService(id)
	return service.Interface().(T), err
}

// GetService returns the asserted service instance for the provided type
// from the provided ServiceGetter
func GetService[T any](sc ServiceGetter) (T, error) {
	return getServiceKey[T](sc, nil)
}

// MustGetService returns the asserted service instance for the provided type
// from the provided ServiceGetter.
//
// Panics if no service is found or any error occurred while creating the instance
func MustGetService[T any](sc ServiceGetter) T {
	service, err := GetService[T](sc)
	if err != nil {
		log.Panicf(
			"[%v]: could not get the service instance, due to error: %s\n",
			reflect.TypeFor[T](),
			err.Error(),
		)
	}

	return service
}

// GetKeyedService returns the asserted service instance for the provided type and key
// from the provided ServiceGetter
func GetKeyedService[T any](sc ServiceGetter, key string) (T, error) {
	return getServiceKey[T](sc, &key)
}

// MustGetKeyedService returns the asserted service instance for the provided type and key
// from the provided ServiceGetter.
//
// Panics if no service is found or any error occurred while creating the instance
func MustGetKeyedService[T any](sc ServiceGetter, key string) T {
	service, err := GetKeyedService[T](sc, key)
	if err != nil {
		log.Panicf(
			"[%v:%s]: could not get the service instance, due to error: %s\n",
			reflect.TypeFor[T](),
			key,
			err.Error(),
		)
	}

	return service
}
