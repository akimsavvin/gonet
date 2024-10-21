// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package env

import (
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type TestCurrentSuite struct {
	suite.Suite
}

// Clear all the environment variables
func (suite *TestCurrentSuite) SetupTest() {
	os.Clearenv()
}

// TestNonEmpty get an empty environment variable and return false
func (suite *TestCurrentSuite) TestEmpty() {
	// Act
	env, ok := Current()

	// Assert
	suite.Empty(env)
	suite.False(ok)
}

// TestNonEmpty get the current environment variable and return it
func (suite *TestCurrentSuite) TestNonEmpty() {
	// Arrange
	os.Setenv("ENVIRONMENT", Staging)

	// Act
	env, ok := Current()

	// Assert
	suite.Equal("Staging", env)
	suite.True(ok)
}

func TestCurrent(t *testing.T) {
	suite.Run(t, new(TestCurrentSuite))
}

type TestCurrentOrDefaultSuite struct {
	suite.Suite
}

// Clear all the environment variables
func (suite *TestCurrentOrDefaultSuite) SetupTest() {
	os.Clearenv()
}

// TestNonEmpty get an empty environment variable and return default value
func (suite *TestCurrentOrDefaultSuite) TestEmpty() {
	// Act
	env := CurrentOrDefault()

	// Assert
	suite.Equal("Development", env)
}

// TestNonEmpty get the current environment variable and return it
func (suite *TestCurrentOrDefaultSuite) TestNonEmpty() {
	// Arrange
	os.Setenv("ENVIRONMENT", Staging)

	// Act
	env := CurrentOrDefault()

	// Assert
	suite.Equal("Staging", env)
}

func TestCurrentOrDefault(t *testing.T) {
	suite.Run(t, new(TestCurrentOrDefaultSuite))
}
