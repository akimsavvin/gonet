// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

import (
	"container/list"
	"iter"
	"reflect"
	"sync"
)

// serviceAccessor is a struct used to get the service instance
type serviceAccessor struct {
	// id is the accessor's serviceIdentifier
	id serviceIdentifier

	// cont is the accessor's Container
	cont *Container

	// once protects the instance from creating multiple times
	once sync.Once

	// factory is the accessor's service factory
	// nil if the service was added with the instance
	factory *serviceFactory

	// instance is the accessor's service instance
	// nil if the service instance has not been requested yet
	// or an error occurred while creating the instance
	instance *reflect.Value

	// instance is the accessor's service factory error
	// nil if the service instance has not been requested yet
	// or an error did not occur while creating the instance
	err error
}

// newServiceAccessor creates a new serviceAccessor
func newServiceAccessor(
	id serviceIdentifier,
	c *Container,
	f *serviceFactory,
	inst *reflect.Value,
) *serviceAccessor {
	return &serviceAccessor{
		id:       id,
		cont:     c,
		factory:  f,
		instance: inst,
	}
}

// createInstance creates the service instance once
func (accessor *serviceAccessor) createInstance() {
	if accessor.instance == nil {
		var instance reflect.Value
		deps, err := accessor.cont.resolveFactoryDeps(accessor.factory)
		if err == nil {
			instance, err = accessor.factory.Call(deps...)
		}

		if err != nil {
			accessor.err = err
			zero := reflect.Zero(accessor.id.Type)
			accessor.instance = &zero
		} else {
			accessor.instance = &instance
		}
	}
}

// Instance returns the accessor service instance
func (accessor *serviceAccessor) Instance() (reflect.Value, error) {
	accessor.once.Do(accessor.createInstance)
	return *accessor.instance, accessor.err
}

// serviceAccessorsList is a singly-linked list of service accessors
type serviceAccessorsList struct {
	// internal is the list.List instance
	internal *list.List
}

// newServiceAccessorsList creates a new serviceAccessorsList
func newServiceAccessorsList(accessors ...*serviceAccessor) *serviceAccessorsList {
	l := &serviceAccessorsList{
		internal: list.New(),
	}

	for _, accessor := range accessors {
		l.Append(accessor)
	}

	return l
}

// Append adds an element to the end of the list
func (list *serviceAccessorsList) Append(accessor *serviceAccessor) {
	list.internal.PushBack(accessor)
}

// Last returns the last element of the list
func (list *serviceAccessorsList) Last() *serviceAccessor {
	accessor, _ := list.internal.Back().Value.(*serviceAccessor)
	return accessor
}

// Len returns the length of the list
func (list *serviceAccessorsList) Len() int {
	return list.internal.Len()
}

// Iter returns an iterator over the elements of the list
func (list *serviceAccessorsList) Iter() iter.Seq2[int, *serviceAccessor] {
	return func(yield func(int, *serviceAccessor) bool) {
		for i, el := 0, list.internal.Front(); el != nil; i, el = i+1, el.Next() {
			if !yield(i, el.Value.(*serviceAccessor)) {
				return
			}
		}
	}
}

// serviceAccessors a map with a list of accessors for the key ID
type serviceAccessors map[serviceIdentifier]*serviceAccessorsList
