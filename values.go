package gonet

import (
	"reflect"
)

func AddValue[TValue any](value TValue) {
	valTyp := reflect.TypeOf((*TValue)(nil)).Elem()
	val := reflect.ValueOf(value)

	p := &provider{
		typ:    valuePr,
		valTyp: valTyp,
		value:  &val,
		constructor: reflect.ValueOf(func() TValue {
			return value
		}),
		lifetime: singleton,
	}

	providers[valTyp] = p
}
