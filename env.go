package gonet

import "os"

const (
	EnvDefault = "production"
)

func GetCurrentEnv() (string, bool) {
	env := os.Getenv("ENVIRONMENT")
	ok := env != ""

	return env, ok
}
