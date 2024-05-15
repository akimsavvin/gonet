package generic

import "reflect"

func GetType[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}
