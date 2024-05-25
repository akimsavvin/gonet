package routing

import (
	"strings"
)

var defaultRouter = newRouter()

type router struct {
	prefix      string
	middlewares []Middleware
}

func newRouter() Router {
	return &router{}
}

func (r *router) handlePattern(method Method, pattern string) string {
	b := strings.Builder{}

	b.WriteString(string(method))

	b.WriteString(r.prefix)
	b.WriteString(globalPrefix)
	b.WriteString(validatePattern(pattern))

	return b.String()
}

func (r *router) handle(method Method, pattern string, handler HandlerFunc, middlewares ...Middleware) {
	pattern = r.handlePattern(method, pattern)
	// TODO: handle pattern somehow
	panic("not implemented")
}

func (r *router) Use(middlewares ...Middleware) {
	r.middlewares = append(r.middlewares, middlewares...)
}

func (r *router) SetPrefix(prefix string) {
	r.prefix = prefix
}

func (r *router) GET(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	r.handle(MethodGet, pattern, handler, middlewares...)
}

func (r *router) HEAD(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	r.handle(MethodHead, pattern, handler, middlewares...)
}

func (r *router) POST(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	r.handle(MethodPost, pattern, handler, middlewares...)
}

func (r *router) PUT(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	r.handle(MethodPut, pattern, handler, middlewares...)
}

func (r *router) PATCH(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	r.handle(MethodPatch, pattern, handler, middlewares...)
}

func (r *router) DELETE(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	r.handle(MethodDelete, pattern, handler, middlewares...)
}

func (r *router) CONNECT(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	r.handle(MethodConnect, pattern, handler, middlewares...)
}

func (r *router) OPTIONS(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	r.handle(MethodOptions, pattern, handler, middlewares...)
}

func (r *router) TRACE(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	r.handle(MethodTrace, pattern, handler, middlewares...)
}
