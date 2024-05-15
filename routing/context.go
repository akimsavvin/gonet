// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package routing

import (
	"context"
	"net/http"
)

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
