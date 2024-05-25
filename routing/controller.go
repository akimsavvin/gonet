// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package routing

type HandlerFunc func(ctx Context) Result

type ControllerBuilder interface {
	SetPrefix(prefix string)
	Router
}

type Controller interface {
	Register(cb ControllerBuilder)
}
