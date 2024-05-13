// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

type Lifetime int

const (
	LifetimeTransient Lifetime = iota + 1
	LifetimeScoped             = iota + 1
	LifetimeSingleton          = iota + 1
)
