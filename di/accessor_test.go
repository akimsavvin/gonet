// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

import (
	"errors"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

// ServiceAccessorInstanceSuite is the suite for testing the serviceAccessor.Instance method
type ServiceAccessorInstanceSuite struct {
	suite.Suite
}

// TestFactoryError tests the serviceAccessor.Instance method
// with the instance set
func (suite *ServiceAccessorInstanceSuite) TestInstance() {
	// Arrange
	id := serviceIdentifier{
		Type:   reflect.TypeFor[string](),
		Key:    "key",
		HasKey: true,
	}
	inst := reflect.ValueOf("instance")

	accessor := newServiceAccessor(id, nil, nil, &inst)

	// Act
	val1, err1 := accessor.Instance()
	val2, err2 := accessor.Instance()

	// Assert
	suite.Equal(inst, val1)
	suite.Equal(val1, val2)
	suite.NoError(err1)
	suite.NoError(err2)
}

// TestFactoryError tests the serviceAccessor.Instance method
// with the factory called once
func (suite *ServiceAccessorInstanceSuite) TestFactory() {
	// Arrange
	id := serviceIdentifier{
		Type:   reflect.TypeFor[string](),
		Key:    "key",
		HasKey: true,
	}

	const value = "instance"
	timesCalled := 0
	f := newServiceFactory(func() string {
		timesCalled++
		return value
	})

	accessor := newServiceAccessor(id, nil, f, nil)

	// Act
	val1, err1 := accessor.Instance()
	val2, err2 := accessor.Instance()

	// Assert
	suite.Equal(1, timesCalled)
	suite.Equal(value, val1.Interface())
	suite.Equal(val1, val2)
	suite.NoError(err1)
	suite.NoError(err2)
}

// TestFactoryError tests the serviceAccessor.Instance method
// with the factory returning an error
func (suite *ServiceAccessorInstanceSuite) TestFactoryError() {
	// Arrange
	id := serviceIdentifier{
		Type:   reflect.TypeFor[string](),
		Key:    "key",
		HasKey: true,
	}

	timesCalled := 0
	f := newServiceFactory(func() (string, error) {
		timesCalled++
		return "", errors.ErrUnsupported
	})

	accessor := newServiceAccessor(id, nil, f, nil)

	// Act
	val1, err1 := accessor.Instance()
	val2, err2 := accessor.Instance()

	// Assert
	suite.Equal(1, timesCalled)
	suite.Empty(val1.Interface())
	suite.Equal(val1, val2)
	suite.Error(err1)
	suite.Error(err2)
}

// TestServiceAccessor_Instance tests the serviceAccessor.Instance method
func TestServiceAccessor_Instance(t *testing.T) {
	suite.Run(t, new(ServiceAccessorInstanceSuite))
}

// TestServiceAccessorsList_Append tests the serviceAccessorsList.Append method
func TestServiceAccessorsList_Append(t *testing.T) {
	// Arrange
	list := newServiceAccessorsList()
	accessor := &serviceAccessor{}
	_ = gofakeit.Struct(accessor)

	// Act
	list.Append(accessor)

	// Assert
	assert.Equal(t, accessor, list.Last())
}

// TestServiceAccessorsList_Last tests the serviceAccessorsList.Last method
func TestServiceAccessorsList_Last(t *testing.T) {
	// Arrange
	accessor := &serviceAccessor{}
	_ = gofakeit.Struct(accessor)
	list := newServiceAccessorsList(accessor)

	// Act
	res := list.Last()

	// Assert
	assert.Equal(t, accessor, res)
}

// TestServiceAccessorsList_Len tests the serviceAccessorsList.Len method
func TestServiceAccessorsList_Len(t *testing.T) {
	// Arrange
	accessor0, accessor1 := &serviceAccessor{}, &serviceAccessor{}
	_ = gofakeit.Struct(accessor0)
	_ = gofakeit.Struct(accessor1)
	list := newServiceAccessorsList(accessor0, accessor1)

	// Act
	l := list.Len()

	// Assert
	assert.Equal(t, 2, l)
}

// TestServiceAccessorsList_Iter tests the serviceAccessorsList.Iter method
func TestServiceAccessorsList_Iter(t *testing.T) {
	// Arrange
	accessor0, accessor1 := &serviceAccessor{}, &serviceAccessor{}
	_ = gofakeit.Struct(accessor0)
	_ = gofakeit.Struct(accessor1)
	list := newServiceAccessorsList(accessor0, accessor1)

	// Act & Assert
	for i, a := range list.Iter() {
		switch i {
		case 0:
			assert.Equal(t, accessor0, a)
		case 1:
			assert.Equal(t, accessor1, a)
		}
	}
}
