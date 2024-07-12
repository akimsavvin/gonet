// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

// Lifetime represents a service lifetime in the di container
type Lifetime int

const (
	// LifetimeTransient a service is created for each request
	LifetimeTransient Lifetime = iota + 1

	// LifetimeScoped a service is created once per scope
	LifetimeScoped

	// LifetimeSingleton a service is created once
	LifetimeSingleton
)
