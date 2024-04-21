package gonet

import "reflect"

func getGenericType[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}
