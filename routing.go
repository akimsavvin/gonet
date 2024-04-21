// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package gonet

import (
	"context"
	"net/http"
)

type H = map[string]any

type Controller interface {
	Register(b ControllerBuilder)
}

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

type HandlerFunc func(ctx Context) ActionResult

type ActionResult interface {
	StatusCode() int
	Payload() any
	HasPayload() bool
	Header() http.Header
}

type NextFunc = HandlerFunc

type Middleware func(next NextFunc) HandlerFunc

type Context interface {
	/* ======================== */
	/* ==== Setters Getters === */
	/* ======================== */

	Context() context.Context

	Request() *http.Request
	SetRequest(req *http.Request)

	Header() http.Header

	Writer() http.ResponseWriter

	/* ======================== */
	/* ======= Methods ======== */
	/* ======================== */

	Param(key string) string
	Query(key string) string
	Queries(key string) []string

	QueryString() string
	Bind(ptr any) bool

	Set(key string, value any)
	Get(key string) any
}
