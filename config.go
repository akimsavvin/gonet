package gonet

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

func newConfig[T any]() *T {
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

	return cfg
}

func AddConfig[T any]() {
	addProvider[*T](newConfig[T], configPr, LifetimeSingleton)
}
