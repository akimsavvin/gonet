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
