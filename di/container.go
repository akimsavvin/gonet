// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

import "reflect"

var defaultContainer = newContainer()

// di container containing service collection
type container struct {
	sc *servColl
}

// creates a new container
func newContainer() *container {
	return &container{
		sc: newServColl(),
	}
}

// addSD adds service descriptor to the collection
func (c *container) addSD(sd *servDescriptor) {
	c.sc.addSD(sd)
}

func (c *container) getTypVal(typ reflect.Type) reflect.Value {
	return c.sc.getTypVal(typ)
}
