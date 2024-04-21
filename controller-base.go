// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package gonet

import (
	"net/http"
)

type ControllerBase struct{}

type Headers = http.Header

type actionResult struct {
	statusCode int
	payload    any
	header     http.Header
}

func (r *actionResult) StatusCode() int {
	return r.statusCode
}

func (r *actionResult) Payload() any {
	return r.payload
}

func (r *actionResult) HasPayload() bool {
	return r.payload != nil
}

func (r *actionResult) Header() http.Header {
	return r.header
}

func (r *actionResult) WithHeader(name string, value string) ActionResult {
	r.header.Add(name, value)
	return r
}

func (r *actionResult) WithHeaders(headers http.Header) ActionResult {
	for k, v := range headers {
		for _, v := range v {
			r.header.Add(k, v)
		}
	}

	return r
}

func (c *ControllerBase) StatusCode(statusCode int, payload any) ActionResult {
	return &actionResult{
		statusCode: statusCode,
		payload:    payload,
		header:     make(http.Header),
	}
}

// Err returns an ActionResult with message in payload equaled to the err.Error()
func (c *ControllerBase) Err(statusCode int, err error) ActionResult {
	return c.StatusCode(statusCode, H{
		"message": err.Error(),
	})
}

func (c *ControllerBase) Ok(payload any) ActionResult {
	return c.StatusCode(StatusOK, payload)
}

func (c *ControllerBase) Created(payload any) ActionResult {
	return c.StatusCode(StatusCreated, payload)
}

func (c *ControllerBase) NoContent(headers ...http.Header) ActionResult {
	return c.StatusCode(StatusNoContent, nil)
}

func (c *ControllerBase) BadRequest(payload any) ActionResult {
	return c.StatusCode(StatusBadRequest, payload)
}

func (c *ControllerBase) Unauthorized(payload any) ActionResult {
	return c.StatusCode(StatusUnauthorized, payload)
}

func (c *ControllerBase) Forbidden(payload any) ActionResult {
	return c.StatusCode(StatusForbidden, payload)
}

func (c *ControllerBase) NotFound(payload any) ActionResult {
	return c.StatusCode(StatusNotFound, payload)
}

func (c *ControllerBase) InternalServerError(payload any) ActionResult {
	return c.StatusCode(StatusInternalServerError, payload)
}
