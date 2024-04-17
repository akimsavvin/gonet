// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package gonet

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"reflect"
)

func newConfig[T any]() T {
	fileName := "config.yaml"

	if env, ok := GetCurrentEnv(); ok {
		fileName = fmt.Sprintf("config.%s.yaml", env)
	}

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		panic(fmt.Sprintf("config file %s does not exist", fileName))
	}

	cfg := new(T)

	if err := cleanenv.ReadConfig(fileName, cfg); err != nil {
		panic(fmt.Sprintf("could not read config file %s, due to error: %s", fileName, err.Error()))
	}

	return *cfg
}

func AddConfig[T any]() {
	cfg := newConfig[T]()

	val := reflect.ValueOf(cfg)
	valTyp := val.Type()

	p := &provider{
		typ:    configPr,
		valTyp: valTyp,
		value:  &val,
		constructor: reflect.ValueOf(func() T {
			return cfg
		}),
		lifetime: singleton,
	}

	providers[valTyp] = p
}

func GetConfig[T any]() T {
	p := getProvider[T]()
	if p.typ != configPr {
		panic(fmt.Sprintf("no config found for type: %v", p.valTyp))
	}

	return p.value.Interface().(T)
}
