// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

import (
	"github.com/akimsavvin/gonet/generic"
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

// TestServiceCollection_AddDescriptor tests the serviceCollection's AddDescriptor method
func TestServiceCollection_AddDescriptor(t *testing.T) {
	// Arrange
	sd := new(serviceDescriptor)
	coll := newServiceCollection()

	// Act
	coll.AddDescriptor(sd)

	// Assert
	assert.Contains(t, coll.Descriptors, sd)
}

// TestServiceCollection_AddServiceFactory tests the serviceCollection's AddServiceFactory method
func TestServiceCollection_AddServiceFactory(t *testing.T) {
	// Arrange
	typ := generic.TypeOf[string]()
	f := func() string {
		return gofakeit.Sentence(5)
	}

	coll := newServiceCollection()

	// Act
	coll.AddServiceFactory(typ, nil, f)

	// Assert
	assert.Equal(t, coll.Descriptors[0].ImplementationType, typ)
}

// TestServiceCollection_AddServiceInstance tests the serviceCollection's AddServiceInstance method
func TestServiceCollection_AddServiceInstance(t *testing.T) {
	// Arrange
	typ := generic.TypeOf[string]()
	inst := gofakeit.Sentence(5)

	coll := newServiceCollection()

	// Act
	coll.AddServiceInstance(typ, nil, inst)

	// Assert
	assert.Equal(t, coll.Descriptors[0].ImplementationType, typ)
}

// TestServiceCollection_AddServiceKey tests the serviceCollection's AddServiceKey method
func TestServiceCollection_AddServiceKey(t *testing.T) {
	// Arrange
	typ := generic.TypeOf[string]()
	key := gofakeit.BuzzWord()
	inst := gofakeit.Sentence(5)

	coll := newServiceCollection()

	// Act
	coll.AddServiceKey(typ, &key, inst)

	// Assert
	assert.True(t, coll.Descriptors[0].HasKey)
	assert.Equal(t, coll.Descriptors[0].Key, key)
	assert.Equal(t, coll.Descriptors[0].ImplementationType, reflect.TypeOf(inst))
}

// TestServiceCollection_AddService tests the serviceCollection's AddService method
func TestServiceCollection_AddService(t *testing.T) {
	// Arrange
	typ := generic.TypeOf[string]()
	inst := gofakeit.Sentence(5)

	coll := newServiceCollection()

	// Act
	coll.AddService(typ, inst)

	// Assert
	assert.False(t, coll.Descriptors[0].HasKey)
	assert.Empty(t, coll.Descriptors[0].Key)
	assert.Equal(t, coll.Descriptors[0].ImplementationType, reflect.TypeOf(inst))
}

// TestServiceCollection_AddKeyedService tests the serviceCollection's AddKeyedService method
func TestServiceCollection_AddKeyedService(t *testing.T) {
	// Arrange
	typ := generic.TypeOf[string]()
	key := gofakeit.BuzzWord()
	f := func() string {
		return gofakeit.Sentence(5)
	}

	coll := newServiceCollection()

	// Act
	coll.AddKeyedService(typ, key, f)

	// Assert
	assert.True(t, coll.Descriptors[0].HasKey)
	assert.Equal(t, coll.Descriptors[0].Key, key)
	assert.Equal(t, coll.Descriptors[0].ImplementationType,
		reflect.TypeOf(f).Out(0))
}

// TestAddService tests the AddService functions
func TestAddService(t *testing.T) {
	// Arrange
	inst := gofakeit.Sentence(5)

	coll := GetServiceCollection()

	// Act
	AddService[string](inst)

	// Assert
	descriptors := coll.descriptors()
	assert.False(t, descriptors[0].HasKey)
	assert.Empty(t, descriptors[0].Key)
	assert.Equal(t, descriptors[0].ImplementationType, reflect.TypeOf(inst))
}

// TestAddKeyedService tests the AddKeyedService functions
func TestAddKeyedService(t *testing.T) {
	// Arrange
	key := gofakeit.BuzzWord()
	f := func() string {
		return gofakeit.Sentence(5)
	}

	coll := GetServiceCollection()

	// Act
	AddKeyedService[string](key, f)

	// Assert
	descriptors := coll.descriptors()
	assert.True(t, descriptors[1].HasKey)
	assert.Equal(t, descriptors[1].Key, key)
	assert.Equal(t, descriptors[1].ImplementationType,
		reflect.TypeOf(f).Out(0))
}

// TestGetServiceCollection tests the GetServiceCollection function
func TestGetServiceCollection(t *testing.T) {
	// Act
	coll := GetServiceCollection()

	// Assert
	assert.IsType(t, &serviceCollection{}, coll)
}
