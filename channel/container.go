// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package channel

import (
	"fmt"
	"reflect"
)

var defaultContainer = newContainer()

type chCollKey string

type container struct {
	// coll collection of channels
	coll map[chCollKey]*reflect.Value
}

func newContainer() *container {
	return &container{
		coll: make(map[chCollKey]*reflect.Value),
	}
}

func (c *container) getKey(name string, typ reflect.Type) chCollKey {
	return chCollKey(fmt.Sprintf("%s:%s", name, typ.Name()))
}

func (c *container) add(name string, typ reflect.Type, val *reflect.Value) {
	key := c.getKey(name, typ)
	c.coll[key] = val
}

func (c *container) get(name string, typ reflect.Type) *reflect.Value {
	key := c.getKey(name, typ)
	return c.coll[key]
}
