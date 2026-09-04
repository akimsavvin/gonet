// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

import "reflect"

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

// newServiceIdentifier creates a new serviceIdentifier
func newServiceIdentifier(typ reflect.Type, key *string) serviceIdentifier {
	id := serviceIdentifier{
		Type:   typ,
		HasKey: key != nil,
	}

	if id.HasKey {
		id.Key = *key
	}

	return id
}

// newDependencyIdentifier creates a serviceIdentifier for a factory dependency
// of the provided type.
//
// Dependencies are always resolved without a key, so a keyed registration does
// not satisfy a factory parameter.
func newDependencyIdentifier(depType reflect.Type) serviceIdentifier {
	return newServiceIdentifier(depType, nil)
}

// resolveLookup normalizes the serviceIdentifier into the identifier actually
// used as the Container accessors map key.
//
// A slice identifier is reduced to its element type, keeping the key, since a
// slice request is served by every accessor registered for the element type.
// The second return value reports whether the identifier was a slice.
func (id serviceIdentifier) resolveLookup() (serviceIdentifier, bool) {
	if id.Type.Kind() != reflect.Slice {
		return id, false
	}

	id.Type = id.Type.Elem()
	return id, true
}

// newDependencyLookup creates the lookup serviceIdentifier for a factory
// dependency of the provided type and reports whether the dependency was a
// slice.
//
// It performs the same two steps the resolution path performs, so a dependency
// this function resolves to an existing identifier is resolvable at runtime,
// and one it does not is not.
func newDependencyLookup(depType reflect.Type) (serviceIdentifier, bool) {
	return newDependencyIdentifier(depType).resolveLookup()
}
