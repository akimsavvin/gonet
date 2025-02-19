// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

// WithServiceSuite is the suite for testing the WithService function
type WithServiceSuite struct {
	suite.Suite
}

// TestInstance tests the instance service
func (suite *WithServiceSuite) TestInstance() {
	// Arrange
	inst := "test"

	opt := WithService[string](inst)
	c := NewContainer(opt)

	id := serviceIdentifier{
		Type: reflect.TypeOf(inst),
	}

	// Act
	lastAccessor := c.accessors[id].Last()

	// Assert
	suite.Equal(id, lastAccessor.id)
	suite.Equal(c, lastAccessor.cont)
	suite.Nil(lastAccessor.factory)
	suite.Equal(inst, lastAccessor.instance.Interface())
	suite.NoError(lastAccessor.err)
}

// TestInstance tests the instance service
func (suite *WithServiceSuite) TestFactory() {
	// Arrange
	inst := "test"
	f := func() string {
		return inst
	}

	opt := WithService[string](f)
	c := NewContainer(opt)

	id := serviceIdentifier{
		Type: reflect.TypeOf(inst),
	}

	// Act
	lastAccessor := c.accessors[id].Last()

	// Assert
	suite.Equal(id, lastAccessor.id)
	suite.Equal(c, lastAccessor.cont)
	suite.NotNil(lastAccessor.factory)
	suite.Nil(lastAccessor.instance)
	suite.NoError(lastAccessor.err)
}

// TestWithService tests the WithService function
func TestWithService(t *testing.T) {
	suite.Run(t, new(WithServiceSuite))
}

// WithKeyedServiceSuite is the suite for testing the WithKeyedService function
type WithKeyedServiceSuite struct {
	suite.Suite
}

// TestInstance tests the instance service
func (suite *WithKeyedServiceSuite) TestInstance() {
	// Arrange
	key := "key"
	inst := "test"

	opt := WithKeyedService[string](key, inst)
	c := NewContainer(opt)

	id := serviceIdentifier{
		Type:   reflect.TypeOf(inst),
		Key:    key,
		HasKey: true,
	}

	// Act
	lastAccessor := c.accessors[id].Last()

	// Assert
	suite.Equal(id, lastAccessor.id)
	suite.Equal(c, lastAccessor.cont)
	suite.Nil(lastAccessor.factory)
	suite.Equal(inst, lastAccessor.instance.Interface())
	suite.NoError(lastAccessor.err)
}

// TestInstance tests the instance service
func (suite *WithKeyedServiceSuite) TestFactory() {
	// Arrange
	key := "key"
	inst := "test"
	f := func() string {
		return inst
	}

	opt := WithKeyedService[string](key, f)
	c := NewContainer(opt)

	id := serviceIdentifier{
		Type:   reflect.TypeOf(inst),
		Key:    key,
		HasKey: true,
	}

	// Act
	lastAccessor := c.accessors[id].Last()

	// Assert
	suite.Equal(id, lastAccessor.id)
	suite.Equal(c, lastAccessor.cont)
	suite.NotNil(lastAccessor.factory)
	suite.Nil(lastAccessor.instance)
	suite.NoError(lastAccessor.err)
}

// TestWithKeyedService tests the WithKeyedService function
func TestWithKeyedService(t *testing.T) {
	suite.Run(t, new(WithKeyedServiceSuite))
}

// TestWithValue tests the WithValue function
func TestWithValue(t *testing.T) {
	// Arrange
	inst := "test"

	opt := WithValue(inst)
	c := NewContainer(opt)

	id := serviceIdentifier{
		Type: reflect.TypeOf(inst),
	}

	// Act
	lastAccessor := c.accessors[id].Last()

	// Assert
	assert.Equal(t, id, lastAccessor.id)
	assert.Equal(t, c, lastAccessor.cont)
	assert.Nil(t, lastAccessor.factory)
	assert.Equal(t, inst, lastAccessor.instance.Interface())
	assert.NoError(t, lastAccessor.err)
}

// TestWithKeyedValue tests the WithKeyedValue function
func TestWithKeyedValue(t *testing.T) {
	// Arrange
	key := "key"
	inst := "test"

	opt := WithKeyedValue[string](key, inst)
	c := NewContainer(opt)

	id := serviceIdentifier{
		Type:   reflect.TypeOf(inst),
		Key:    key,
		HasKey: true,
	}

	// Act
	lastAccessor := c.accessors[id].Last()

	// Assert
	assert.Equal(t, id, lastAccessor.id)
	assert.Equal(t, c, lastAccessor.cont)
	assert.Nil(t, lastAccessor.factory)
	assert.Equal(t, inst, lastAccessor.instance.Interface())
	assert.NoError(t, lastAccessor.err)
}

// TestWithFactory tests the WithFactory function
func TestWithFactory(t *testing.T) {
	// Arrange
	inst := "test"
	f := func() string {
		return inst
	}

	opt := WithFactory(f)
	c := NewContainer(opt)

	id := serviceIdentifier{
		Type: reflect.TypeOf(inst),
	}

	// Act
	lastAccessor := c.accessors[id].Last()
	res, err := lastAccessor.factory.Call()

	// Assert
	assert.Equal(t, id, lastAccessor.id)
	assert.Equal(t, c, lastAccessor.cont)
	assert.Nil(t, lastAccessor.instance)
	assert.NotNil(t, lastAccessor.factory)
	assert.NoError(t, lastAccessor.err)

	assert.Equal(t, inst, res.Interface())
	assert.NoError(t, err)
}

// TestWithKeyedFactory tests the WithKeyedFactory function
func TestWithKeyedFactory(t *testing.T) {
	// Arrange
	key := "key"
	inst := "test"
	f := func() string {
		return inst
	}

	opt := WithKeyedFactory(key, f)
	c := NewContainer(opt)

	id := serviceIdentifier{
		Type:   reflect.TypeOf(inst),
		Key:    key,
		HasKey: true,
	}

	// Act
	lastAccessor := c.accessors[id].Last()
	res, err := lastAccessor.factory.Call()

	// Assert
	assert.Equal(t, id, lastAccessor.id)
	assert.Equal(t, c, lastAccessor.cont)
	assert.Nil(t, lastAccessor.instance)
	assert.NotNil(t, lastAccessor.factory)
	assert.NoError(t, lastAccessor.err)

	assert.Equal(t, inst, res.Interface())
	assert.NoError(t, err)
}

// TestMultiple tests the adding multiple services with the same identifier
func TestMultiple(t *testing.T) {
	// Arrange
	inst1, inst2 := "test1", "test2"
	opt1, opt2 := WithValue(inst1), WithValue(inst2)

	id := serviceIdentifier{
		Type: reflect.TypeFor[string](),
	}

	// Act
	c := NewContainer(opt1, opt2)
	res := make([]*serviceAccessor, 0, 2)
	for _, a := range c.accessors[id].Iter() {
		res = append(res, a)
	}

	// Assert
	if assert.Equal(t, 2, len(res)) {
		assert.Equal(t, inst2, res[1].instance.Interface())
		assert.Equal(t, inst1, res[0].instance.Interface())
	}
}

// GetServiceSuite is the suite for testing the GetService function
type GetServiceSuite struct {
	suite.Suite
}

// TestFactory tests the factory service
func (suite *GetServiceSuite) TestFactory() {
	// Arrange
	inst := "test"
	f := func() string {
		return inst
	}
	c := NewContainer(WithFactory(f))

	// Act
	res, err := GetService[string](c)

	// Assert
	suite.Equal(inst, res)
	suite.NoError(err)
}

// TestFactoryError tests the factory service error
func (suite *GetServiceSuite) TestFactoryError() {
	// Arrange
	f := func() (string, error) {
		return "", errors.ErrUnsupported
	}
	c := NewContainer(WithFactory(f))

	// Act
	res, err := GetService[string](c)

	// Assert
	suite.Empty(res)
	suite.Error(err)
}

// TestInstance tests the instance service
func (suite *GetServiceSuite) TestInstance() {
	// Arrange
	inst := "test"
	c := NewContainer(WithValue(inst))

	// Act
	res, err := GetService[string](c)

	// Assert
	suite.Equal(inst, res)
	suite.NoError(err)
}

// TestGetService tests the GetService function
func TestGetService(t *testing.T) {
	suite.Run(t, new(GetServiceSuite))
}

// GetKeyedServiceSuite is the suite for testing the GetKeyedService function
type GetKeyedServiceSuite struct {
	suite.Suite
}

// TestInstance tests the factory service
func (suite *GetKeyedServiceSuite) TestFactory() {
	// Arrange
	key := "key"
	inst := "test"
	f := func() string {
		return inst
	}
	c := NewContainer(WithKeyedFactory(key, f))

	// Act
	res, err := GetKeyedService[string](c, key)

	// Assert
	suite.Equal(inst, res)
	suite.NoError(err)
}

// TestInstance tests the factory service error
func (suite *GetKeyedServiceSuite) TestFactoryError() {
	// Arrange
	key := "key"
	f := func() (string, error) {
		return "", errors.ErrUnsupported
	}
	c := NewContainer(WithKeyedFactory(key, f))

	// Act
	res, err := GetKeyedService[string](c, key)

	// Assert
	suite.Empty(res)
	suite.Error(err)
}

// TestInstance tests the instance service
func (suite *GetKeyedServiceSuite) TestInstance() {
	// Arrange
	key := "key"
	inst := "test"
	c := NewContainer(WithKeyedValue(key, inst))

	// Act
	res, err := GetKeyedService[string](c, key)

	// Assert
	suite.Equal(inst, res)
	suite.NoError(err)
}

// TestGetKeyedService tests the GetKeyedService function
func TestGetKeyedService(t *testing.T) {
	suite.Run(t, new(GetKeyedServiceSuite))
}

// MustGetServiceSuite is the suite for testing the MustGetService function
type MustGetServiceSuite struct {
	suite.Suite
}

// TestFactory tests the factory service
func (suite *MustGetServiceSuite) TestFactory() {
	// Arrange
	inst := "test"
	f := func() string {
		return inst
	}
	c := NewContainer(WithFactory(f))

	// Act & Assert
	var res string
	suite.NotPanics(func() {
		res = MustGetService[string](c)
	})
	suite.Equal(inst, res)
}

// TestFactoryError tests the factory service error
func (suite *MustGetServiceSuite) TestFactoryError() {
	// Arrange
	f := func() (string, error) {
		return "", errors.ErrUnsupported
	}
	c := NewContainer(WithFactory(f))

	// Act & Assert
	suite.Panics(func() {
		MustGetService[string](c)
	})
}

// TestInstance tests the instance service
func (suite *MustGetServiceSuite) TestInstance() {
	// Arrange
	inst := "test"
	c := NewContainer(WithValue(inst))

	// Act & Assert
	var res string
	suite.NotPanics(func() {
		res = MustGetService[string](c)
	})
	suite.Equal(inst, res)
}

// TestMustGetService tests the MustGetService function
func TestMustGetService(t *testing.T) {
	suite.Run(t, new(MustGetServiceSuite))
}

// MustGetKeyedServiceSuite is the suite for testing the MustGetKeyedService function
type MustGetKeyedServiceSuite struct {
	suite.Suite
}

// TestFactory tests the factory service
func (suite *MustGetKeyedServiceSuite) TestFactory() {
	// Arrange
	key := "key"
	inst := "test"
	f := func() string {
		return inst
	}
	c := NewContainer(WithKeyedFactory(key, f))

	// Act & Assert
	var res string
	suite.NotPanics(func() {
		res = MustGetKeyedService[string](c, key)
	})
	suite.Equal(inst, res)
}

// TestFactoryError tests the factory service error
func (suite *MustGetKeyedServiceSuite) TestFactoryError() {
	// Arrange
	key := "key"
	f := func() (string, error) {
		return "", errors.ErrUnsupported
	}
	c := NewContainer(WithKeyedFactory(key, f))

	// Act & Assert
	suite.Panics(func() {
		MustGetKeyedService[string](c, key)
	})
}

// TestInstance tests the instance service
func (suite *MustGetKeyedServiceSuite) TestInstance() {
	// Arrange
	key := "key"
	inst := "test"
	c := NewContainer(WithKeyedValue(key, inst))

	// Act & Assert
	var res string
	suite.NotPanics(func() {
		res = MustGetKeyedService[string](c, key)
	})
	suite.Equal(inst, res)
}

// TestMustGetKeyedService tests the MustGetKeyedService function
func TestMustGetKeyedService(t *testing.T) {
	suite.Run(t, new(MustGetKeyedServiceSuite))
}
