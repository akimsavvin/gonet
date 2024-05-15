// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package routing

import "net/http"

type Result interface {
	/* ======================== */
	/* ======= Getters ======== */
	/* ======================== */
	StatusCode() int
	Payload() any
	HasPayload() bool
	Header() http.Header

	/* ======================== */
	/* ======= Headers ======== */
	/* ======================== */
	WithHeader(name string, value string) Result
	WithHeaders(headers http.Header) Result
}
