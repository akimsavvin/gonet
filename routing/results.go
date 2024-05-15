// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package routing

import (
	"net/http"
)

type H = map[string]any

type result struct {
	statusCode int
	payload    any
	header     http.Header
}

func (r *result) StatusCode() int {
	return r.statusCode
}

func (r *result) Payload() any {
	return r.payload
}

func (r *result) HasPayload() bool {
	return r.payload != nil
}

func (r *result) Header() http.Header {
	return r.header
}

func (r *result) WithHeader(name string, value string) Result {
	r.header.Add(name, value)
	return r
}

func (r *result) WithHeaders(headers http.Header) Result {
	for k, v := range headers {
		for _, v := range v {
			r.header.Add(k, v)
		}
	}

	return r
}

func StatusCode(statusCode int, payload any) Result {
	return &result{
		statusCode: statusCode,
		payload:    payload,
		header:     make(http.Header),
	}
}

// Err returns an Result with message in payload equaled to the err.Error()
func Err(statusCode int, err error) Result {
	return StatusCode(statusCode, H{
		"message": err.Error(),
	})
}

func Ok(payload any) Result {
	return StatusCode(StatusOK, payload)
}

func Created(payload any) Result {
	return StatusCode(StatusCreated, payload)
}

func NoContent(headers ...http.Header) Result {
	return StatusCode(StatusNoContent, nil)
}

func BadRequest(payload any) Result {
	return StatusCode(StatusBadRequest, payload)
}

func Unauthorized(payload any) Result {
	return StatusCode(StatusUnauthorized, payload)
}

func Forbidden(payload any) Result {
	return StatusCode(StatusForbidden, payload)
}

func NotFound(payload any) Result {
	return StatusCode(StatusNotFound, payload)
}

func InternalServerError(payload any) Result {
	return StatusCode(StatusInternalServerError, payload)
}
