// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

import (
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

// NewServiceFactorySuite is a test suite for the newServiceIdentifier function
type NewServiceIdentifierSuite struct {
	suite.Suite
}

// TestNilKey tests the nil key case
func (suite *NewServiceIdentifierSuite) TestNilKey() {
	// Arrange
	typ := reflect.TypeFor[string]()

	// Act
	id := newServiceIdentifier(typ, nil)

	// Assert
	suite.Equal(typ, id.Type)
	suite.Empty(id.Key)
	suite.False(id.HasKey)
}

// TestNilKey tests the not nil key case
func (suite *NewServiceIdentifierSuite) TestNotNilKey() {
	// Arrange
	typ := reflect.TypeFor[string]()
	key := "key"

	// Act
	id := newServiceIdentifier(typ, &key)

	// Assert
	suite.Equal(typ, id.Type)
	suite.Equal(key, id.Key)
	suite.True(id.HasKey)
}

// TestNewServiceIdentifier tests the newServiceIdentifier function
func TestNewServiceIdentifier(t *testing.T) {
	suite.Run(t, new(NewServiceIdentifierSuite))
}
