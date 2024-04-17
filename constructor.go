// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package gonet

import (
	"fmt"
	"reflect"
)

func validateConstructor(constructor any) {
	typ := reflect.TypeOf(constructor)

	if typ.Kind() != reflect.Func {
		panic(fmt.Sprintf("constructor %v is not a function", constructor))
	}

	if typ.NumOut() != 1 {
		panic(fmt.Sprintf("%v constructor must have only one return argument", constructor))
	}

	const prefix = "New"
	name := typ.Name()
	if len(name) != 0 && name[:3] != prefix {
		panic(fmt.Sprintf("constructor's %v name does not start with '%s' prefix", name, prefix))
	}
}
