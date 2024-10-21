// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package env

import "os"

const Key = "ENVIRONMENT"

const (
	Development = "Development"
	Staging     = "Staging"
	Production  = "Production"
)

func Current() (string, bool) {
	env := os.Getenv(Key)
	ok := env != ""

	return env, ok
}

func CurrentOrDefault() string {
	if env, ok := Current(); ok {
		return env
	}

	return Development
}
