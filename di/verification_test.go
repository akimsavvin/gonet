// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

// MissingDependencyErrorSuite is the suite for testing the MissingDependencyError type
type MissingDependencyErrorSuite struct {
	suite.Suite
}

// TestFields tests both types are readable from the fault value
func (suite *MissingDependencyErrorSuite) TestFields() {
	// Arrange
	requesting := reflect.TypeFor[string]()
	dependency := reflect.TypeFor[int]()

	// Act
	fault := &MissingDependencyError{
		RequestingType: requesting,
		DependencyType: dependency,
	}

	// Assert
	suite.Equal(requesting, fault.RequestingType)
	suite.Equal(dependency, fault.DependencyType)
}

// TestErrorsIs tests the fault wraps the ErrServiceNotFound sentinel
func (suite *MissingDependencyErrorSuite) TestErrorsIs() {
	// Arrange
	fault := &MissingDependencyError{
		RequestingType: reflect.TypeFor[string](),
		DependencyType: reflect.TypeFor[int](),
	}

	// Act & Assert
	suite.ErrorIs(fault, ErrServiceNotFound)
}

// TestMessage tests the message names both types
func (suite *MissingDependencyErrorSuite) TestMessage() {
	// Arrange
	fault := &MissingDependencyError{
		RequestingType: reflect.TypeFor[*testService](),
		DependencyType: reflect.TypeFor[testRepository](),
	}

	// Act
	msg := fault.Error()

	// Assert
	suite.Contains(msg, "di.testService")
	suite.Contains(msg, "di.testRepository")
}

// TestMissingDependencyError tests the MissingDependencyError type
func TestMissingDependencyError(t *testing.T) {
	suite.Run(t, new(MissingDependencyErrorSuite))
}

// CircularDependencyErrorSuite is the suite for testing the CircularDependencyError type
type CircularDependencyErrorSuite struct {
	suite.Suite
}

// TestCycleRoundTrip tests the ordered types round-trip through the fault
func (suite *CircularDependencyErrorSuite) TestCycleRoundTrip() {
	// Arrange
	cycle := []reflect.Type{
		reflect.TypeFor[*testServiceA](),
		reflect.TypeFor[*testServiceB](),
		reflect.TypeFor[*testServiceC](),
	}

	// Act
	fault := &CircularDependencyError{Cycle: cycle}

	// Assert
	suite.Equal(cycle, fault.Cycle)
}

// TestMessage tests the message names every participant of the cycle
func (suite *CircularDependencyErrorSuite) TestMessage() {
	// Arrange
	fault := &CircularDependencyError{
		Cycle: []reflect.Type{
			reflect.TypeFor[*testServiceA](),
			reflect.TypeFor[*testServiceB](),
			reflect.TypeFor[*testServiceC](),
		},
	}

	// Act
	msg := fault.Error()

	// Assert
	suite.Contains(msg, "di.testServiceA")
	suite.Contains(msg, "di.testServiceB")
	suite.Contains(msg, "di.testServiceC")
	suite.Contains(msg, "circular dependency")
}

// TestMessageClosesCycle tests the message returns to the starting type
func (suite *CircularDependencyErrorSuite) TestMessageClosesCycle() {
	// Arrange
	fault := &CircularDependencyError{
		Cycle: []reflect.Type{
			reflect.TypeFor[*testServiceA](),
			reflect.TypeFor[*testServiceB](),
		},
	}

	// Act
	msg := fault.Error()

	// Assert
	suite.Equal(
		"circular dependency: *di.testServiceA -> *di.testServiceB -> *di.testServiceA",
		msg,
	)
}

// TestCircularDependencyError tests the CircularDependencyError type
func TestCircularDependencyError(t *testing.T) {
	suite.Run(t, new(CircularDependencyErrorSuite))
}

// VerificationErrorSuite is the suite for testing the VerificationError type
type VerificationErrorSuite struct {
	suite.Suite
}

// TestErrorsAs tests errors.As finds both fault kinds inside the aggregate
func (suite *VerificationErrorSuite) TestErrorsAs() {
	// Arrange
	missing := &MissingDependencyError{
		RequestingType: reflect.TypeFor[*testService](),
		DependencyType: reflect.TypeFor[testRepository](),
	}
	circular := &CircularDependencyError{
		Cycle: []reflect.Type{
			reflect.TypeFor[*testServiceA](),
			reflect.TypeFor[*testServiceB](),
		},
	}

	var err error = &VerificationError{Faults: []error{missing, circular}}

	// Act & Assert
	var gotMissing *MissingDependencyError
	var gotCircular *CircularDependencyError

	suite.ErrorAs(err, &gotMissing)
	suite.ErrorAs(err, &gotCircular)
	suite.ErrorIs(err, ErrServiceNotFound)
	suite.Equal(missing, gotMissing)
	suite.Equal(circular, gotCircular)
}

// TestMessage tests the message lists every fault under a header
func (suite *VerificationErrorSuite) TestMessage() {
	// Arrange
	err := &VerificationError{
		Faults: []error{
			&MissingDependencyError{
				RequestingType: reflect.TypeFor[*testService](),
				DependencyType: reflect.TypeFor[testRepository](),
			},
			&CircularDependencyError{
				Cycle: []reflect.Type{
					reflect.TypeFor[*testServiceA](),
					reflect.TypeFor[*testServiceB](),
				},
			},
		},
	}

	// Act
	msg := err.Error()

	// Assert
	suite.Contains(msg, "di: container verification failed:")
	suite.Contains(msg, "\n  - service \"*di.testService\" requires \"di.testRepository\", which is not registered")
	suite.Contains(msg, "\n  - circular dependency: *di.testServiceA -> *di.testServiceB -> *di.testServiceA")
}

// TestUnwrap tests every fault is exposed for programmatic inspection
func (suite *VerificationErrorSuite) TestUnwrap() {
	// Arrange
	faults := []error{
		&MissingDependencyError{
			RequestingType: reflect.TypeFor[*testService](),
			DependencyType: reflect.TypeFor[testRepository](),
		},
		&CircularDependencyError{Cycle: []reflect.Type{reflect.TypeFor[*testServiceA]()}},
	}
	err := &VerificationError{Faults: faults}

	// Act
	unwrapped := err.Unwrap()

	// Assert
	suite.Equal(faults, unwrapped)
}

// TestVerificationError tests the VerificationError type
func TestVerificationError(t *testing.T) {
	suite.Run(t, new(VerificationErrorSuite))
}
