// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

import "reflect"

type scope struct {
	cont *container
	sc   *servColl
}

func newScope(cont *container) *scope {
	return &scope{
		cont: cont,
		sc:   cont.getScopedColl(),
	}
}

func (s *scope) getTypVal(typ reflect.Type) reflect.Value {
	scopedSD := s.sc.getTypSD(typ)
	if scopedSD == nil {
		return s.cont.getTypVal(typ)
	}

	return s.sc.getTypVal(typ)
}
