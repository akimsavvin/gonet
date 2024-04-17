// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package gonet

import (
	"reflect"
)

func AddValue(value any) {
	val := reflect.ValueOf(value)
	valTyp := val.Type()

	p := &provider{
		typ:    valuePr,
		valTyp: valTyp,
		value:  &val,
		constructor: reflect.ValueOf(func() any {
			return value
		}),
		lifetime: singleton,
	}

	providers[valTyp] = p
}
