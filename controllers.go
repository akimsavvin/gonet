package gonet

import (
	"fmt"
	"reflect"
)

type Controller interface {
	Register(g *ApiGroup)
}

func AddController(constructor any) {
	validateConstructor(constructor)
	typ := reflect.TypeOf(constructor).Out(0)
	ctrlTyp := reflect.TypeOf((*Controller)(nil)).Elem()
	if !typ.Implements(ctrlTyp) {
		panic(fmt.Sprintf("controller %s does not implement Constructor interface", typ.Name()))
	}

	addTypeProvider(typ, singleton, controllerPr, constructor)
}

func getControllers() []Controller {
	ctrls := make([]Controller, 0)

	for _, p := range providers {
		c, ok := p.getInstance().Interface().(Controller)

		if p.typ == controllerPr && ok {
			ctrls = append(ctrls, c)
		}
	}

	return ctrls
}
