// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package gonet

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"sync"
	"unicode/utf8"
)

type ctx struct {
	req     *http.Request
	rwriter http.ResponseWriter
	// Mutex for keys
	mu   sync.RWMutex
	keys map[string]any
	err  error
}

func newCtx(rwriter http.ResponseWriter, req *http.Request) *ctx {
	return &ctx{
		req:     req,
		rwriter: rwriter,
		mu:      sync.RWMutex{},
		keys:    make(map[string]any),
	}
}

/* ======================== */
/* ==== Setters Getters === */
/* ======================== */

func (ctx *ctx) Context() context.Context {
	return ctx.req.Context()
}

func (ctx *ctx) Request() *http.Request {
	return ctx.req
}

func (ctx *ctx) SetRequest(req *http.Request) {
	ctx.req = req
}

func (ctx *ctx) Header() http.Header {
	return ctx.req.Header
}

func (ctx *ctx) Writer() http.ResponseWriter {
	return ctx.rwriter
}

/* ======================== */
/* ======= Methods ======== */
/* ======================== */

func (ctx *ctx) Param(key string) string {
	return ctx.req.PathValue(key)
}

func (ctx *ctx) Query(key string) string {
	queries := ctx.Queries(key)

	if len(queries) == 0 {
		return ""
	}

	return queries[0]
}

func (ctx *ctx) Queries(key string) []string {
	return ctx.req.URL.Query()[key]
}

func (ctx *ctx) QueryString() string {
	return ctx.req.URL.RawQuery
}

type ErrValidationError struct {
	Message string `json:"message"`
}

func (err ErrValidationError) Error() string {
	return err.Message
}

type ErrInvalidBody struct {
	Message string               `json:"message"`
	Errors  []ErrValidationError `json:"errors" xml:"errors" yaml:"errors" bson:"errors"`
}

func (err ErrInvalidBody) Error() string {
	return err.Message
}

func getTagIntValue(structTag reflect.StructTag, name string) (int, bool) {
	tag, ok := structTag.Lookup(name)

	if !ok {
		return 0, false
	}

	num, err := strconv.ParseInt(tag, 10, 64)
	if err != nil {
		return 0, false
	}

	return int(num), true
}

func (ctx *ctx) hasErr() bool {
	return ctx.err != nil
}

func (ctx *ctx) writeErr(err error) {
	ctx.err = err

	ctx.rwriter.Header().Set(HeaderContentType, MIMEApplicationJSONCharsetUTF8)
	ctx.rwriter.WriteHeader(http.StatusInternalServerError)

	serialized, _ := json.Marshal(H{
		"message": err.Error(),
	})

	_, err = ctx.rwriter.Write(serialized)
	if err != nil {
		return
	}
}

func (ctx *ctx) Bind(ptr any) bool {
	defer ctx.req.Body.Close()

	// TODO: validate for empty body
	err := json.NewDecoder(ctx.req.Body).Decode(ptr)
	if err != nil {
		ctx.writeErr(err)
		return false
	}

	val := reflect.ValueOf(ptr).Elem()
	typ := val.Type()
	validErrs := make([]ErrValidationError, 0)

	for i := 0; i < typ.NumField(); i++ {
		structField := typ.Field(i)
		fieldVal := val.Field(i)

		minVal, minExists := getTagIntValue(structField.Tag, "min")
		maxVal, maxExists := getTagIntValue(structField.Tag, "max")

		switch structField.Type.Kind() {
		case reflect.String:
			length := utf8.RuneCountInString(fieldVal.String())

			if minExists && length < minVal {
				validErrs = append(validErrs, ErrValidationError{
					fmt.Sprintf(
						"%s length must be at least %d symbols",
						structField.Name, minVal,
					),
				})
			}

			if maxExists && length > maxVal {
				validErrs = append(validErrs, ErrValidationError{
					fmt.Sprintf(
						"%s length must be not more than %d symbols",
						structField.Name, maxVal,
					),
				})
			}
		case reflect.Array, reflect.Slice:
			if minExists && fieldVal.Len() < minVal {
				validErrs = append(validErrs, ErrValidationError{
					fmt.Sprintf(
						"%s must have at least %d elements",
						structField.Name, minVal,
					),
				})
			}

			if maxExists && fieldVal.Len() < maxVal {
				validErrs = append(validErrs, ErrValidationError{
					fmt.Sprintf(
						"%s must have not more than %d elements",
						structField.Name, maxVal,
					),
				})
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if minExists && int(fieldVal.Int()) < minVal {

				validErrs = append(validErrs, ErrValidationError{
					fmt.Sprintf(
						"%s value must be at least %d",
						structField.Name, minVal,
					),
				})
			}

			if maxExists && int(fieldVal.Int()) > maxVal {
				validErrs = append(validErrs, ErrValidationError{
					fmt.Sprintf(
						"%s value must be not more than %d",
						structField.Name, maxVal,
					),
				})
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			if minExists && int(fieldVal.Uint()) < minVal {
				validErrs = append(validErrs, ErrValidationError{
					fmt.Sprintf(
						"%s value must be at least %d",
						structField.Name, minVal,
					),
				})
			}

			if maxExists && int(fieldVal.Uint()) > maxVal {
				validErrs = append(validErrs, ErrValidationError{
					fmt.Sprintf(
						"%s value must be not more than %d",
						structField.Name, maxVal,
					),
				})
			}
		default:
		}
	}

	if len(validErrs) != 0 {
		err = ErrInvalidBody{
			Message: "Some errors occurred while validating your request",
			Errors:  validErrs,
		}

		ctx.err = err

		serialized, _ := json.Marshal(err)

		ctx.rwriter.Header().Set(HeaderContentType, MIMEApplicationJSONCharsetUTF8)
		ctx.rwriter.WriteHeader(StatusBadRequest)
		_, err = ctx.rwriter.Write(serialized)
		if err != nil {
			http.Error(ctx.rwriter, err.Error(), StatusInternalServerError)
		}

		return false
	}

	return true
}

func (ctx *ctx) Set(key string, value any) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.keys[key] = value
}

func (ctx *ctx) Get(key string) any {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	return ctx.keys[key]
}
