package routing

import (
	"github.com/akimsavvin/gonet/generic"
	"log"
	"reflect"
	"sync/atomic"
)

var AreControllersEnabled = atomic.Bool{}

type contlrColl struct {
	controllers []Controller
}

func processRequest() {
}

func mustValidateContrlConstr(constrVal reflect.Value) {
	contrlTyp := generic.GetType[Controller]()
	constrTyp := constrVal.Type()
	retTyp := constrTyp.Out(0)

	if !retTyp.Implements(contrlTyp) {
		log.Panicf("Controller %d does not implement routing.Controller interface", retTyp)
	}

}

func AddController(constr any) {
	AreControllersEnabled.Store(true)

}
