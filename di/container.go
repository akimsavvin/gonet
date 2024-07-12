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
		sc: newServColl(LifetimeSingleton),
	}
}

// addSD adds service descriptor to the collection
func (c *container) addSD(sd *servDescriptor) {
	c.sc.addSD(sd)
}

// getTypSD gets service descriptor by type
func (c *container) getTypSD(typ reflect.Type) *servDescriptor {
	return c.sc.getTypSD(typ)
}

// getTypSD gets service descriptor by type
func (c *container) getScopedColl() *servColl {
	return c.sc.getScopedColl()
}

func (c *container) getTypVal(typ reflect.Type) reflect.Value {
	return c.sc.getTypVal(typ)
}

// newScope creates a new Scope for container
func (c *container) newScope() *Scope {
	return newScope(c)
}

// NewScope creates a new Scope
func NewScope() *Scope {
	return defaultContainer.newScope()
}
