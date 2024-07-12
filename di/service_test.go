// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type Service interface {
	Test() int
}

type MyService struct{}

func NewMyService() *MyService {
	return new(MyService)
}

func (s *MyService) Test() int {
	return 5
}

func TestAddService(t *testing.T) {
	// Arrange & Act & Assert
	assert.NotPanics(t, func() {
		// Arrange
		AddService[Service](NewMyService)

		// Act
		s := GetService[Service]()
		res := s.Test()

		// Assert
		assert.Equal(t, res, 5)
	})
}

func TestAddSingleton(t *testing.T) {
	// Arrange & Act & Assert
	assert.NotPanics(t, func() {
		// Arrange
		AddSingleton[Service](NewMyService)

		// Act
		s1 := GetService[Service]()
		s2 := GetService[Service]()
		res := s1.Test()

		// Assert
		assert.Equal(t, res, 5)
		assert.Equal(t, s1, s2)
	})
}

func TestAddScoped(t *testing.T) {
	// Arrange & Act & Assert
	assert.NotPanics(t, func() {
		// Arrange
		AddScoped[Service](NewMyService)
		scope := NewScope()

		// Act
		s1 := GetScopedService[Service](scope)
		s2 := GetScopedService[Service](scope)
		res := s1.Test()

		// Assert
		assert.Equal(t, res, 5)
		assert.Same(t, s1, s2)
	})
}
