package routing

func validatePattern(pattern string) string {
	if pattern[0] != '/' {
		pattern = "/" + pattern
	}

	if pattern[len(pattern)-1] == '/' {
		pattern = pattern[:len(pattern)-1]
	}

	return pattern
}

var globalPrefix = ""

func SetGlobalPrefix(prefix string) {
	globalPrefix = prefix
}

type Router interface {
	Use(middlewares ...Middleware)
	SetPrefix(prefix string)

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

func Use(middlewares ...Middleware) {
	defaultRouter.Use(middlewares...)
}

func GET(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	defaultRouter.GET(pattern, handler, middlewares...)
}

func HEAD(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	defaultRouter.HEAD(pattern, handler, middlewares...)
}

func POST(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	defaultRouter.POST(pattern, handler, middlewares...)
}

func PUT(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	defaultRouter.PUT(pattern, handler, middlewares...)
}

func PATCH(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	defaultRouter.PATCH(pattern, handler, middlewares...)
}

func DELETE(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	defaultRouter.DELETE(pattern, handler, middlewares...)
}

func CONNECT(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	defaultRouter.CONNECT(pattern, handler, middlewares...)
}

func OPTIONS(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	defaultRouter.OPTIONS(pattern, handler, middlewares...)
}

func TRACE(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	defaultRouter.TRACE(pattern, handler, middlewares...)
}
