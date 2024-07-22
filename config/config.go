// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package config

import (
	"github.com/akimsavvin/gonet/di"
	"github.com/akimsavvin/gonet/env"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"strings"
)

type Type int

const (
	JSON Type = iota + 1
	YAML
	ENV
)

func buildCfgFileName(typ Type) string {
	cfgExts := map[Type]string{
		JSON: "json",
		YAML: "yaml",
		ENV:  "env",
	}

	var cfgNameBuilder strings.Builder
	cfgNameBuilder.WriteString("config.")

	if env, ok := env.Current(); ok {
		cfgNameBuilder.WriteString(env + ".")
	}

	cfgNameBuilder.WriteString(cfgExts[typ])

	return cfgNameBuilder.String()
}

func newConfig[T any](typ Type) T {
	fileName := buildCfgFileName(typ)

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		log.Panicf("Config file %s does not exist\n", fileName)
	}

	cfg := new(T)

	if err := cleanenv.ReadConfig(fileName, cfg); err != nil {
		log.Panicf("Can not read config file %s, due to error: %s\n", fileName, err.Error())
	}

	return *cfg
}

func Add[T any](cfgTyp Type) {
	cfg := newConfig[T](cfgTyp)
	di.AddValue(cfg)
}
