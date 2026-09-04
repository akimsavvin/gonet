// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

// VerifySuite is the suite for testing the Container.verify method
type VerifySuite struct {
	suite.Suite
}

// faults returns the aggregated faults of the provided verification error
func (suite *VerifySuite) faults(err error) []error {
	suite.Require().Error(err)

	var verErr *VerificationError
	suite.Require().ErrorAs(err, &verErr)

	return verErr.Faults
}

// missing returns the missing dependency faults of the provided error
func (suite *VerifySuite) missing(err error) []*MissingDependencyError {
	var res []*MissingDependencyError
	for _, fault := range suite.faults(err) {
		if missing, ok := fault.(*MissingDependencyError); ok {
			res = append(res, missing)
		}
	}

	return res
}

// cycles returns the circular dependency faults of the provided error
func (suite *VerifySuite) cycles(err error) []*CircularDependencyError {
	var res []*CircularDependencyError
	for _, fault := range suite.faults(err) {
		if cycle, ok := fault.(*CircularDependencyError); ok {
			res = append(res, cycle)
		}
	}

	return res
}

// TestSoundGraph tests a graph with every dependency satisfied has no faults
func (suite *VerifySuite) TestSoundGraph() {
	// Arrange
	c := newContainer(
		WithService[testRepository](newTestRepository),
		WithFactory(newTestService),
		WithFactory(newTestController),
	)

	// Act
	err := c.verify()

	// Assert
	suite.NoError(err)
}

// TestEmptyContainer tests a container with no registrations has no faults
func (suite *VerifySuite) TestEmptyContainer() {
	// Arrange
	c := newContainer()

	// Act
	err := c.verify()

	// Assert
	suite.NoError(err)
}

// TestSharedDependencyExpandedOnce tests a dependency reached from two parents
// is expanded once, so its faults are reported once.
//
// Both controllers depend on the same testService, whose own testRepository
// dependency is not registered. A second expansion of the testService node would
// report the same missing dependency twice.
func (suite *VerifySuite) TestSharedDependencyExpandedOnce() {
	// Arrange
	c := newContainer(
		WithFactory(newTestController),
		WithFactory(newTestSecondController),
		WithFactory(newTestService),
	)

	// Act
	missing := suite.missing(c.verify())

	// Assert
	suite.Len(missing, 1)
	suite.Equal(reflect.TypeFor[*testService](), missing[0].RequestingType)
	suite.Equal(reflect.TypeFor[testRepository](), missing[0].DependencyType)
}

// TestSharedDependencyExpandedOnceAcrossRoots tests the single expansion holds
// whichever parent the walk reaches the shared node from first
func (suite *VerifySuite) TestSharedDependencyExpandedOnceAcrossRoots() {
	// Arrange
	opts := []Option{
		WithFactory(newTestController),
		WithFactory(newTestSecondController),
		WithFactory(newTestService),
	}

	// Act
	for i := range opts {
		rotated := make([]Option, 0, len(opts))
		rotated = append(rotated, opts[i:]...)
		rotated = append(rotated, opts[:i]...)

		missing := suite.missing(newContainer(rotated...).verify())

		// Assert
		suite.Len(missing, 1)
	}
}

// TestDeepChainTerminates tests a long dependency chain is walked without
// re-expanding a node
func (suite *VerifySuite) TestDeepChainTerminates() {
	// Arrange
	c := newContainer(
		WithFactory(newTestController),
		WithFactory(newTestSecondController),
		WithFactory(newTestService),
		WithService[testRepository](newTestRepository),
	)

	// Act
	err := c.verify()

	// Assert
	suite.NoError(err)
}

// TestDirectMissingDependency tests an unregistered factory parameter is
// reported with the requesting and dependency types
func (suite *VerifySuite) TestDirectMissingDependency() {
	// Arrange
	c := newContainer(WithFactory(newTestService))

	// Act
	err := c.verify()
	missing := suite.missing(err)

	// Assert
	suite.Len(missing, 1)
	suite.Equal(reflect.TypeFor[*testService](), missing[0].RequestingType)
	suite.Equal(reflect.TypeFor[testRepository](), missing[0].DependencyType)
	suite.ErrorIs(err, ErrServiceNotFound)
}

// TestTransitiveMissingDependency tests the fault is reported at the level where
// the dependency is actually missing
func (suite *VerifySuite) TestTransitiveMissingDependency() {
	// Arrange
	c := newContainer(
		WithFactory(newTestController),
		WithFactory(newTestService),
	)

	// Act
	missing := suite.missing(c.verify())

	// Assert
	suite.Len(missing, 1)
	suite.Equal(reflect.TypeFor[*testService](), missing[0].RequestingType)
	suite.Equal(reflect.TypeFor[testRepository](), missing[0].DependencyType)
}

// TestRegisteredInstanceSatisfiesDependency tests a value registration satisfies
// a factory parameter
func (suite *VerifySuite) TestRegisteredInstanceSatisfiesDependency() {
	// Arrange
	c := newContainer(
		WithValue[testRepository](&testRepositoryImpl{}),
		WithFactory(newTestService),
	)

	// Act
	err := c.verify()

	// Assert
	suite.NoError(err)
}

// TestContainerDependencyIsSatisfied tests a factory may depend on the Container
func (suite *VerifySuite) TestContainerDependencyIsSatisfied() {
	// Arrange
	c := newContainer(
		WithFactory(func(cont *Container) *testService { return &testService{} }),
	)

	// Act
	err := c.verify()

	// Assert
	suite.NoError(err)
}

// TestKeyedRegistrationDoesNotSatisfyUnkeyedDependency tests a keyed only
// registration is reported as missing, matching runtime resolution
func (suite *VerifySuite) TestKeyedRegistrationDoesNotSatisfyUnkeyedDependency() {
	// Arrange
	c := newContainer(
		WithKeyedService[testRepository]("db", newTestRepository),
		WithFactory(func(repo testRepository) *testKeyedConsumer {
			return &testKeyedConsumer{repo}
		}),
	)

	// Act
	missing := suite.missing(c.verify())

	// Assert
	suite.Len(missing, 1)
	suite.Equal(reflect.TypeFor[*testKeyedConsumer](), missing[0].RequestingType)
	suite.Equal(reflect.TypeFor[testRepository](), missing[0].DependencyType)
}

// TestSliceDependencyWithRegistrations tests a slice dependency with at least
// one registration is satisfied
func (suite *VerifySuite) TestSliceDependencyWithRegistrations() {
	// Arrange
	c := newContainer(
		WithService[testRepository](newTestRepository),
		WithFactory(newTestCollector),
	)

	// Act
	err := c.verify()

	// Assert
	suite.NoError(err)
}

// TestSliceDependencyWithoutRegistrations tests a slice dependency with no
// registration reports the element type as missing
func (suite *VerifySuite) TestSliceDependencyWithoutRegistrations() {
	// Arrange
	c := newContainer(WithFactory(newTestCollector))

	// Act
	missing := suite.missing(c.verify())

	// Assert
	suite.Len(missing, 1)
	suite.Equal(reflect.TypeFor[*testCollector](), missing[0].RequestingType)
	suite.Equal(reflect.TypeFor[[]testRepository](), missing[0].DependencyType)
}

// TestTwoServiceCycle tests a two service cycle is reported and does not deadlock
func (suite *VerifySuite) TestTwoServiceCycle() {
	// Arrange
	c := newContainer(
		WithFactory(func(b *testServiceB) *testServiceA { return &testServiceA{b} }),
		WithFactory(func(a *testServiceA) *testServiceB { return &testServiceB{} }),
	)

	// Act
	cycles := suite.cycles(c.verify())

	// Assert
	suite.Len(cycles, 1)
	suite.Equal(
		[]reflect.Type{
			reflect.TypeFor[*testServiceA](),
			reflect.TypeFor[*testServiceB](),
		},
		cycles[0].Cycle,
	)
}

// TestLongerCycle tests a three service cycle names the ordered sequence
func (suite *VerifySuite) TestLongerCycle() {
	// Arrange
	c := newContainer(
		WithFactory(func(b *testServiceB) *testServiceA { return &testServiceA{b} }),
		WithFactory(func(cs *testServiceC) *testServiceB { return &testServiceB{cs} }),
		WithFactory(func(a *testServiceA) *testServiceC { return &testServiceC{a} }),
	)

	// Act
	cycles := suite.cycles(c.verify())

	// Assert
	suite.Len(cycles, 1)
	suite.Equal(
		[]reflect.Type{
			reflect.TypeFor[*testServiceA](),
			reflect.TypeFor[*testServiceB](),
			reflect.TypeFor[*testServiceC](),
		},
		cycles[0].Cycle,
	)
}

// TestSelfDependency tests a factory depending on its own return type is a cycle
func (suite *VerifySuite) TestSelfDependency() {
	// Arrange
	c := newContainer(
		WithFactory(func(self *testSelfDependent) *testSelfDependent {
			return &testSelfDependent{self}
		}),
	)

	// Act
	cycles := suite.cycles(c.verify())

	// Assert
	suite.Len(cycles, 1)
	suite.Equal([]reflect.Type{reflect.TypeFor[*testSelfDependent]()}, cycles[0].Cycle)
}

// TestSharedDependencyIsNotACycle tests two services sharing a dependency report
// no cycle
func (suite *VerifySuite) TestSharedDependencyIsNotACycle() {
	// Arrange
	c := newContainer(
		WithService[testRepository](newTestRepository),
		WithFactory(newTestService),
		WithFactory(func(repo testRepository) *testCollector { return &testCollector{} }),
	)

	// Act
	err := c.verify()

	// Assert
	suite.NoError(err)
}

// TestCycleCanonicalStart tests a cycle is reported from its canonical starting
// type, not from the node the walk happened to enter it through.
//
// testCycleEntry sorts first, so the walk reaches the cycle through it and
// enters at testCycleZ, while the canonical start is testCycleM.
func (suite *VerifySuite) TestCycleCanonicalStart() {
	// Arrange
	c := newContainer(
		WithFactory(func(z *testCycleZ) *testCycleEntry { return &testCycleEntry{z} }),
		WithFactory(func(m *testCycleM) *testCycleZ { return &testCycleZ{m} }),
		WithFactory(func(z *testCycleZ) *testCycleM { return &testCycleM{z} }),
	)

	// Act
	cycles := suite.cycles(c.verify())

	// Assert
	suite.Len(cycles, 1)
	suite.Equal(
		[]reflect.Type{
			reflect.TypeFor[*testCycleM](),
			reflect.TypeFor[*testCycleZ](),
		},
		cycles[0].Cycle,
	)
}

// TestDuplicateRegistrationCycleReportedOnce tests two accessor cycles with the
// same type sequence are deduplicated into a single fault.
//
// testCycleZ is registered twice, so the testCycleM edge fans out to both
// accessors and the walk closes the same cycle of types through each of them.
func (suite *VerifySuite) TestDuplicateRegistrationCycleReportedOnce() {
	// Arrange
	c := newContainer(
		WithService[*testCycleZ](func(m *testCycleM) *testCycleZ { return &testCycleZ{m} }),
		WithService[*testCycleZ](func(m *testCycleM) *testCycleZ { return &testCycleZ{m} }),
		WithFactory(func(z *testCycleZ) *testCycleM { return &testCycleM{z} }),
	)

	// Act
	cycles := suite.cycles(c.verify())

	// Assert
	suite.Len(cycles, 1)
	suite.Equal(
		[]reflect.Type{
			reflect.TypeFor[*testCycleM](),
			reflect.TypeFor[*testCycleZ](),
		},
		cycles[0].Cycle,
	)
}

// TestCycleReportedOncePerRegistrationOrder tests a three service cycle yields
// exactly one fault whatever order its factories are registered in
func (suite *VerifySuite) TestCycleReportedOncePerRegistrationOrder() {
	// Arrange
	opts := []Option{
		WithFactory(func(b *testServiceB) *testServiceA { return &testServiceA{b} }),
		WithFactory(func(cs *testServiceC) *testServiceB { return &testServiceB{cs} }),
		WithFactory(func(a *testServiceA) *testServiceC { return &testServiceC{a} }),
	}

	// Act
	for i := range opts {
		rotated := make([]Option, 0, len(opts))
		rotated = append(rotated, opts[i:]...)
		rotated = append(rotated, opts[:i]...)

		cycles := suite.cycles(newContainer(rotated...).verify())

		// Assert
		suite.Len(cycles, 1)
		suite.Equal(
			[]reflect.Type{
				reflect.TypeFor[*testServiceA](),
				reflect.TypeFor[*testServiceB](),
				reflect.TypeFor[*testServiceC](),
			},
			cycles[0].Cycle,
		)
	}
}

// TestMultipleMissingDependencies tests one verification reports every missing
// dependency
func (suite *VerifySuite) TestMultipleMissingDependencies() {
	// Arrange
	c := newContainer(
		WithFactory(newTestService),
		WithFactory(newTestCollector),
		WithFactory(func(logger testLogger) *testMissingLoggerService {
			return &testMissingLoggerService{logger}
		}),
	)

	// Act
	missing := suite.missing(c.verify())

	// Assert
	suite.Len(missing, 3)
}

// TestMixedFaultKinds tests one verification reports a missing dependency and a
// cycle together
func (suite *VerifySuite) TestMixedFaultKinds() {
	// Arrange
	c := newContainer(
		WithFactory(newTestService),
		WithFactory(func(b *testServiceB) *testServiceA { return &testServiceA{b} }),
		WithFactory(func(a *testServiceA) *testServiceB { return &testServiceB{} }),
	)

	// Act
	err := c.verify()

	// Assert
	suite.Len(suite.missing(err), 1)
	suite.Len(suite.cycles(err), 1)
	suite.Len(suite.faults(err), 2)
}

// TestVerify tests the Container.verify method
func TestVerify(t *testing.T) {
	suite.Run(t, new(VerifySuite))
}
