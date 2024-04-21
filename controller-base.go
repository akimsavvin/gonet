// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package gonet

import (
	"net/http"
)

type ControllerBase struct{}

type actionResult struct {
	statusCode int
	payload    any
	headers    http.Header
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
	return r.headers
}

func (c *ControllerBase) StatusCode(statusCode int, payload any, headers ...http.Header) ActionResult {
	result := &actionResult{
		statusCode: statusCode,
		payload:    payload,
	}

	if len(headers) > 0 {
		result.headers = headers[0]
	}

	return result
}

// Err returns an ActionResult with message in payload equaled to the err.Error()
func (c *ControllerBase) Err(statusCode int, err error, headers ...http.Header) ActionResult {
	return c.StatusCode(statusCode, H{
		"message": err.Error(),
	}, headers...)
}

func (c *ControllerBase) Ok(payload any, headers ...http.Header) ActionResult {
	return c.StatusCode(StatusOK, payload, headers...)
}

func (c *ControllerBase) Created(payload any, headers ...http.Header) ActionResult {
	return c.StatusCode(StatusCreated, payload, headers...)
}

func (c *ControllerBase) NoContent(headers ...http.Header) ActionResult {
	return c.StatusCode(StatusNoContent, nil, headers...)
}

func (c *ControllerBase) BadRequest(payload any, headers ...http.Header) ActionResult {
	return c.StatusCode(StatusBadRequest, payload, headers...)
}

func (c *ControllerBase) Unauthorized(payload any, headers ...http.Header) ActionResult {
	return c.StatusCode(StatusUnauthorized, payload, headers...)
}

func (c *ControllerBase) Forbidden(payload any, headers ...http.Header) ActionResult {
	return c.StatusCode(StatusForbidden, payload, headers...)
}

func (c *ControllerBase) NotFound(payload any, headers ...http.Header) ActionResult {
	return c.StatusCode(StatusNotFound, payload, headers...)
}

func (c *ControllerBase) InternalServerError(payload any, headers ...http.Header) ActionResult {
	return c.StatusCode(StatusInternalServerError, payload, headers...)
}
