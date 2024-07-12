// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test_mustValidateConstrVal_Success(t *testing.T) {
	// Arrange
	type testStruct struct{}

	data := []struct {
		Type      reflect.Type
		ConstrVal reflect.Value
	}{
		{
			reflect.TypeOf(new(testStruct)),
			reflect.ValueOf(func() *testStruct {
				return new(testStruct)
			}),
		},
		{
			reflect.TypeOf(new(testStruct)),
			reflect.ValueOf(func() (*testStruct, error) {
				return new(testStruct), nil
			}),
		},
	}

	// Act & Assert
	for _, d := range data {
		assert.NotPanicsf(t, func() {
			mustValidateConstrVal(d.Type, d.ConstrVal)
		}, "Validation paniced with valid return type")
	}
}

func Test_mustValidateConstrVal_InvalidReturnType(t *testing.T) {
	// Arrange
	type testStruct struct{}
	type otherTestStruct struct{}

	typ := reflect.TypeOf(new(testStruct))
	constrVal := reflect.ValueOf(func() *otherTestStruct {
		return new(otherTestStruct)
	})

	// Act & Assert
	assert.Panicsf(t, func() {
		mustValidateConstrVal(typ, constrVal)
	}, "Validation did not panic with invalid return type")
}

func Test_mustValidateConstrVal_TooMuchReturnValues(t *testing.T) {
	// Arrange
	type testStruct struct{}
	type otherTestStruct struct{}

	typ := reflect.TypeOf(new(testStruct))
	constrVal := reflect.ValueOf(func() (*testStruct, *otherTestStruct, error) {
		return new(testStruct), new(otherTestStruct), nil
	})

	// Act & Assert
	assert.Panicsf(t, func() {
		mustValidateConstrVal(typ, constrVal)
	}, "Validation did not panic with too much return values")
}

func Test_mustValidateConstrVal_InvalidSecondReturnValue(t *testing.T) {
	// Arrange
	type testStruct struct{}
	type otherTestStruct struct{}

	typ := reflect.TypeOf(new(testStruct))
	constrVal := reflect.ValueOf(func() (*testStruct, *otherTestStruct) {
		return new(testStruct), new(otherTestStruct)
	})

	// Act & Assert
	assert.Panicsf(t, func() {
		mustValidateConstrVal(typ, constrVal)
	}, "Validation did not panic with not error second return value")
}
