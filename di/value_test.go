// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddValue(t *testing.T) {
	// Arrange
	type MyValue struct {
		Value int
	}

	AddValue(MyValue{5})

	// Act
	v := GetService[MyValue]()

	// Assert
	assert.Equal(t, 5, v.Value)
}
