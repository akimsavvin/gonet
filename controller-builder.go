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

func newCBuilder(mux *http.ServeMux) *cBuilder {
	return &cBuilder{
		mux: mux,
	}
}

// cBuilder is a controller builder passed as the first argument in the controller Register method
type cBuilder struct {
	mux         *http.ServeMux
	prefix      string
	middlewares []Middleware
}

func applyMiddleware(h HandlerFunc, middlewares ...Middleware) HandlerFunc {
	for _, m := range middlewares {
		h = m(h)
	}

	return h
}

func (cb *cBuilder) SetPrefix(prefix string) {
	cb.prefix = prefix
}

func (cb *cBuilder) Use(middlewares ...Middleware) {
	cb.middlewares = append(cb.middlewares, middlewares...)
}

func (cb *cBuilder) GET(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	cb.handle(MethodGet, pattern, handler, middlewares...)
}

func (cb *cBuilder) HEAD(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	cb.handle(MethodHead, pattern, handler, middlewares...)
}

func (cb *cBuilder) POST(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	cb.handle(MethodPost, pattern, handler, middlewares...)
}

func (cb *cBuilder) PUT(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	cb.handle(MethodPut, pattern, handler, middlewares...)
}

func (cb *cBuilder) PATCH(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	cb.handle(MethodPatch, pattern, handler, middlewares...)
}

func (cb *cBuilder) DELETE(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	cb.handle(MethodDelete, pattern, handler, middlewares...)
}

func (cb *cBuilder) CONNECT(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	cb.handle(MethodConnect, pattern, handler, middlewares...)
}

func (cb *cBuilder) OPTIONS(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	cb.handle(MethodOptions, pattern, handler, middlewares...)
}

func (cb *cBuilder) TRACE(pattern string, handler HandlerFunc, middlewares ...Middleware) {
	cb.handle(MethodTrace, pattern, handler, middlewares...)
}

func (cb *cBuilder) buildHandlerPattern(method Method, pattern string) string {
	var patternBuilder strings.Builder
	patternBuilder.WriteString(string(method))
	patternBuilder.WriteString(" ")

	if globalPrefix != "" {
		patternBuilder.WriteString(globalPrefix)
	}

	if cb.prefix != "" {
		patternBuilder.WriteString(cb.prefix)
	}

	pattern = validatePrefix(pattern)
	patternBuilder.WriteString(pattern)

	return patternBuilder.String()
}

func (cb *cBuilder) handle(method Method, pattern string, handler HandlerFunc, middlewares ...Middleware) {
	handlerPattern := cb.buildHandlerPattern(method, pattern)

	middlewares = append(cb.middlewares, middlewares...)
	handler = applyMiddleware(handler, middlewares...)

	cb.mux.HandleFunc(handlerPattern, func(w http.ResponseWriter, req *http.Request) {
		ctx := newCtx(w, req)
		res := handler(ctx)

		if ctx.hasErr() {
			return
		}

		if res.HasPayload() {
			for k, v := range res.Header() {
				for _, v := range v {
					w.Header().Add(k, v)
				}
			}

			if w.Header().Get(HeaderContentType) == "" {
				switch reflect.TypeOf(res.Payload()).Kind() {
				case reflect.String:
					w.Header().Set(HeaderContentType, MIMETextPlainCharsetUTF8)
				case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
					w.Header().Set(HeaderContentType, MIMEApplicationJSONCharsetUTF8)
				default:
					panic("invalid content type")
				}
			}

			w.WriteHeader(res.StatusCode())

			serialized, err := json.Marshal(res.Payload())
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
			w.WriteHeader(res.StatusCode())
		}
	})
}
