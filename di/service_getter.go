// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

import "reflect"

// ServiceGetter is an interface for getting a service
type ServiceGetter interface {
	// getService gets a service instance for the provided service identifier
	getService(id serviceIdentifier) (reflect.Value, error)
}
