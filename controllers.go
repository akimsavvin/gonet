// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package gonet

import (
	"fmt"
	"github.com/akimsavvin/gonet/generic"
	"reflect"
)

func AddController(constructor any) {
	typ := reflect.TypeOf(constructor).Out(0)
	ctrlTyp := generic.GetType[Controller]()
	if !typ.Implements(ctrlTyp) {
		panic(fmt.Sprintf("controller %s does not implement gonet.Controller interface", typ))
	}

	//addTypeProvider(singleton, controllerPr, constructor, typ)
}

func getControllers() []Controller {
	ctrls := make([]Controller, 0)

	//for _, p := range providers {
	//	c, ok := p.getInstance().Interface().(Controller)
	//
	//	if p.typ == controllerPr && ok {
	//		ctrls = append(ctrls, c)
	//	}
	//}

	return ctrls
}
