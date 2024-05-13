// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package gonet

import "os"

const (
	EnvDefault = "production"
)

func GetCurrentEnv() (string, bool) {
	env := os.Getenv("GONET_ENVIRONMENT")
	ok := env != ""

	return env, ok
}
