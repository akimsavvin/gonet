// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

import "reflect"

var defaultContainer = newContainer()

type container struct {
	sc *servColl
}

func newContainer() *container {
	return &container{
		sc: newServColl(LifetimeSingleton),
	}
}

func (c *container) addSD(sd *servDescriptor) {
	c.sc.addSD(sd)
}

func (c *container) getTypSD(typ reflect.Type) *servDescriptor {
	return c.sc.getTypSD(typ)
}

func (c *container) getScopedColl() *servColl {
	return c.sc.getScopedColl()
}

func (c *container) getTypVal(typ reflect.Type) reflect.Value {
	return c.sc.getTypVal(typ)
}

func (c *container) newScope() *scope {
	return newScope(c)
}
