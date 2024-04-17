// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package gonet

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
)

type router struct {
	Handler     *http.ServeMux
	prefix      string
	middlewares []Middleware
}

func newRouter() *router {
	return &router{
		Handler: http.NewServeMux(),
	}
}

func applyMiddleware(h HandlerFunc, middlewares ...Middleware) HandlerFunc {
	for _, m := range middlewares {
		h = m(h)
	}

	return h
}

func (r *router) SetPrefix(prefix string) {
	r.prefix = prefix
}

func (r *router) Use(middlewares ...Middleware) {
	r.middlewares = append(r.middlewares, middlewares...)
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

func (r *router) getHandlerPattern(method Method, pattern string) string {
	var patternBuilder strings.Builder
	patternBuilder.WriteString(string(method))
	patternBuilder.WriteString(" ")

	if globalPrefix != "" {
		patternBuilder.WriteString(globalPrefix)
	}

	if r.prefix != "" {
		patternBuilder.WriteString(r.prefix)
	}

	pattern = validatePrefix(pattern)
	patternBuilder.WriteString(pattern)

	return patternBuilder.String()
}

func (r *router) handle(method Method, pattern string, handler HandlerFunc, middlewares ...Middleware) {
	handlerPattern := r.getHandlerPattern(method, pattern)

	middlewares = append(r.middlewares, middlewares...)
	handler = applyMiddleware(handler, middlewares...)

	r.Handler.HandleFunc(handlerPattern, func(w http.ResponseWriter, req *http.Request) {
		ctx := newCtx(w, req)
		result := handler(ctx)

		if ctx.hasErr() {
			return
		}

		if result.HasPayload() {
			if w.Header().Get(HeaderContentType) == "" {
				switch reflect.TypeOf(result.Payload()).Kind() {
				case reflect.String:
					w.Header().Set(HeaderContentType, MIMETextPlainCharsetUTF8)
				case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
					w.Header().Set(HeaderContentType, MIMEApplicationJSONCharsetUTF8)
				default:
					panic("invalid content type")
				}
			}

			w.WriteHeader(result.StatusCode())

			serialized, err := json.Marshal(result.Payload())
			if err != nil {
				ctx.writeErr(err)
				return
			}

			_, err = w.Write(serialized)
			if err != nil {
				ctx.writeErr(err)
				return
			}
		} else {
			w.WriteHeader(result.StatusCode())
		}
	})
}
