// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package env

import "os"

const (
	Development = "Development"
	Staging     = "Staging"
	Production  = "Production"
)

func Current() (string, bool) {
	env := os.Getenv("ENVIRONMENT")
	ok := env != ""

	return env, ok
}

func CurrentOrDefault() string {
	env, ok := Current()

	if !ok {
		return Development
	}

	return env
}
