// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

// NewServiceFactorySuite test suite for newServiceFactory
type NewServiceFactorySuite struct {
	suite.Suite
}

// TestUsual tests the usual case
func (suite *NewServiceFactorySuite) TestUsual() {
	// Arrange
	type MyStruct struct{}

	// Act & Assert
	suite.NotPanics(func() {
		f := newServiceFactory(func(int, string) (s MyStruct) {
			return
		})

		if suite.NotNil(f) {
			suite.Equal(2, f.DepsCount)
		}
	})
}

// TestWithError tests the constructor with an error returned
func (suite *NewServiceFactorySuite) TestWithError() {
	// Arrange
	type MyStruct struct{}

	// Act & Assert
	suite.NotPanics(func() {
		f := newServiceFactory(func() (s MyStruct, err error) {
			return
		})

		if suite.NotNil(f) {
			suite.True(f.HasErr)
		}
	})
}

// TestNonErrorSecondReturnArgument tests the constructor with an error returned
func (suite *NewServiceFactorySuite) TestNonErrorSecondReturnArgument() {
	// Arrange
	type (
		MyStruct  struct{}
		MyStruct2 struct{}
	)

	// Act & Assert
	suite.Panics(func() {
		newServiceFactory(func() (s MyStruct, s2 MyStruct2) {
			return
		})
	})
}

// TestNonFunctionArgument tests the case when client provides a non-function argument
func (suite *NewServiceFactorySuite) TestNonFunctionArgument() {
	// Act & Assert
	suite.Panics(func() {
		newServiceFactory(0.1)
		newServiceFactory(1)
		newServiceFactory("string")
		newServiceFactory(struct{}{})
	})
}

// TestNoReturnArguments tests the case the provided function has no return arguments
func (suite *NewServiceFactorySuite) TestNoReturnArguments() {
	// Act & Assert
	suite.Panics(func() {
		newServiceFactory(func() {})
	})
}

// TestTooMuchReturnArguments tests the case when the function has too many return arguments
func (suite *NewServiceFactorySuite) TestTooMuchReturnArguments() {
	// Act & Assert
	suite.Panics(func() {
		newServiceFactory(func() (v1 int, err error, v2 int) {
			return
		})
	})
}

// TestNewServiceFactory tests the newServiceFactory function
func TestNewServiceFactory(t *testing.T) {
	suite.Run(t, new(NewServiceFactorySuite))
}
