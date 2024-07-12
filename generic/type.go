// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package generic

import "reflect"

func GetType[T any]() reflect.Type {
	return reflect.TypeOf(new(T)).Elem()
}
