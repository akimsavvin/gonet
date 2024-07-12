// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

import "reflect"

// Scope represents a scope for dependency injection
type Scope struct {
	cont *container
	sc   *servColl
}

// creates a new Scope for the container
func newScope(cont *container) *Scope {
	return &Scope{
		cont: cont,
		sc:   cont.getScopedColl(),
	}
}

// getTypVal gets a value for provided type
func (s *Scope) getTypVal(typ reflect.Type) reflect.Value {
	scopedSD := s.sc.getTypSD(typ)
	if scopedSD == nil {
		return s.cont.getTypVal(typ)
	}

	return s.sc.getTypVal(typ)
}
