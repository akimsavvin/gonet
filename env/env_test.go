// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package env

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCurrent(t *testing.T) {
	// Arrange
	cases := []struct {
		Before    func()
		Expected0 string
		Expected1 bool
	}{
		{
			Before:    func() {},
			Expected0: "",
			Expected1: false,
		},
		{
			Before: func() {
				os.Setenv("GONET_ENVIRONMENT", "Staging")
			},
			Expected0: "Staging",
			Expected1: true,
		},
	}

	for _, c := range cases {
		// Act
		c.Before()
		res0, res1 := Current()

		// Assert
		assert.Equal(t, c.Expected0, res0)
		assert.Equal(t, c.Expected1, res1)
	}
}

func TestCurrentOrDefault(t *testing.T) {
	// Arrange
	cases := []struct {
		Before   func()
		Expected string
	}{
		{
			Before:   func() {},
			Expected: "Development",
		},
		{
			Before: func() {
				os.Setenv("GONET_ENVIRONMENT", "Staging")
			},
			Expected: "Staging",
		},
	}

	for _, c := range cases {
		// Act
		c.Before()
		res := CurrentOrDefault()

		// Assert
		assert.Equal(t, c.Expected, res)
	}
}
