// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package routing

type HandlerFunc func(ctx Context) Result

type ControllerBuilder interface {
	SetPrefix(prefix string)
	Use(middlewares ...Middleware)

	GET(pattern string, handler HandlerFunc, middlewares ...Middleware)
	HEAD(pattern string, handler HandlerFunc, middlewares ...Middleware)
	POST(pattern string, handler HandlerFunc, middlewares ...Middleware)
	PUT(pattern string, handler HandlerFunc, middlewares ...Middleware)
	PATCH(pattern string, handler HandlerFunc, middlewares ...Middleware)
	DELETE(pattern string, handler HandlerFunc, middlewares ...Middleware)
	CONNECT(pattern string, handler HandlerFunc, middlewares ...Middleware)
	OPTIONS(pattern string, handler HandlerFunc, middlewares ...Middleware)
	TRACE(pattern string, handler HandlerFunc, middlewares ...Middleware)
}

type Controller interface {
	Register(cb ControllerBuilder)
}
