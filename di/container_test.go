// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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
	c, err := NewContainer(opt)
	suite.Require().NoError(err)

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
	c, err := NewContainer(opt)
	suite.Require().NoError(err)

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
	c, err := NewContainer(opt)
	suite.Require().NoError(err)

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
	c, err := NewContainer(opt)
	suite.Require().NoError(err)

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
	c, err := NewContainer(opt)
	require.NoError(t, err)

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
	c, err := NewContainer(opt)
	require.NoError(t, err)

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
	c, err := NewContainer(opt)
	require.NoError(t, err)

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
	c, err := NewContainer(opt)
	require.NoError(t, err)

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
	c, err := NewContainer(opt1, opt2)
	require.NoError(t, err)

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

// NewContainerSuite is the suite for testing the NewContainer function
type NewContainerSuite struct {
	suite.Suite
}

// TestSoundGraph tests a sound graph yields a usable Container and no error
func (suite *NewContainerSuite) TestSoundGraph() {
	// Arrange
	opts := []Option{
		WithService[testRepository](newTestRepository),
		WithFactory(newTestService),
		WithFactory(newTestController),
	}

	// Act
	c, err := NewContainer(opts...)

	// Assert
	suite.Require().NoError(err)
	suite.Require().NotNil(c)

	controller, err := c.GetService[*testController]()
	suite.NoError(err)
	suite.NotNil(controller.Service.Repo)
}

// TestEmptyContainer tests a Container with no registrations is sound
func (suite *NewContainerSuite) TestEmptyContainer() {
	// Act
	c, err := NewContainer()

	// Assert
	suite.NoError(err)
	suite.NotNil(c)
}

// TestFaultyGraphYieldsNilContainer tests a faulty graph yields a nil Container
// and a describing error
func (suite *NewContainerSuite) TestFaultyGraphYieldsNilContainer() {
	// Act
	c, err := NewContainer(WithFactory(newTestService))

	// Assert
	suite.Nil(c)
	suite.Require().Error(err)

	var verErr *VerificationError
	suite.Require().ErrorAs(err, &verErr)
	suite.Len(verErr.Faults, 1)
	suite.Contains(err.Error(), "di: container verification failed:")
	suite.ErrorIs(err, ErrServiceNotFound)
}

// TestCycleDoesNotDeadlock tests a cycle is returned as an error instead of
// deadlocking the process
func (suite *NewContainerSuite) TestCycleDoesNotDeadlock() {
	// Act
	c, err := NewContainer(
		WithFactory(func(b *testServiceB) *testServiceA { return &testServiceA{b} }),
		WithFactory(func(a *testServiceA) *testServiceB { return &testServiceB{} }),
	)

	// Assert
	suite.Nil(c)
	suite.Require().Error(err)

	var cycleErr *CircularDependencyError
	suite.Require().ErrorAs(err, &cycleErr)
	suite.Equal(
		[]reflect.Type{
			reflect.TypeFor[*testServiceA](),
			reflect.TypeFor[*testServiceB](),
		},
		cycleErr.Cycle,
	)
}

// TestFactoriesNotCalled tests no factory is invoked while the Container is
// constructed
func (suite *NewContainerSuite) TestFactoriesNotCalled() {
	// Arrange
	repoCalls, serviceCalls, controllerCalls := 0, 0, 0

	// Act
	c, err := NewContainer(
		WithService[testRepository](func() testRepository {
			repoCalls++
			return &testRepositoryImpl{}
		}),
		WithFactory(func(repo testRepository) *testService {
			serviceCalls++
			return &testService{repo}
		}),
		WithFactory(func(service *testService) *testController {
			controllerCalls++
			return &testController{service}
		}),
	)

	// Assert
	suite.Require().NoError(err)
	suite.Require().NotNil(c)
	suite.Equal(0, repoCalls)
	suite.Equal(0, serviceCalls)
	suite.Equal(0, controllerCalls)
}

// TestLazySingleton tests a service is created on first resolution and reused
// afterwards
func (suite *NewContainerSuite) TestLazySingleton() {
	// Arrange
	calls := 0
	c, err := NewContainer(
		WithService[testRepository](func() testRepository {
			calls++
			return &testRepositoryImpl{}
		}),
	)
	suite.Require().NoError(err)

	// Act
	first, err1 := c.GetService[testRepository]()
	second, err2 := c.GetService[testRepository]()

	// Assert
	suite.Equal(1, calls)
	suite.NoError(err1)
	suite.NoError(err2)
	suite.Same(first, second)
}

// TestFactoryErrorSurfacesOnResolution tests a factory error does not fail the
// construction and is returned on resolution
func (suite *NewContainerSuite) TestFactoryErrorSurfacesOnResolution() {
	// Arrange
	calls := 0

	// Act
	c, err := NewContainer(
		WithService[testRepository](func() (testRepository, error) {
			calls++
			return nil, errors.ErrUnsupported
		}),
		WithFactory(newTestService),
	)

	// Assert
	suite.Require().NoError(err)
	suite.Require().NotNil(c)
	suite.Equal(0, calls)

	service, resErr := c.GetService[*testService]()
	suite.Nil(service)
	suite.Require().Error(resErr)
	suite.ErrorIs(resErr, errors.ErrUnsupported)
	suite.Equal(1, calls)
}

// TestNewContainer tests the NewContainer function
func TestNewContainer(t *testing.T) {
	suite.Run(t, new(NewContainerSuite))
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
	c, err := NewContainer(WithFactory(f))
	suite.Require().NoError(err)

	// Act
	res, err := c.GetService[string]()

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
	c, err := NewContainer(WithFactory(f))
	suite.Require().NoError(err)

	// Act
	res, err := c.GetService[string]()

	// Assert
	suite.Empty(res)
	suite.Error(err)
}

// TestInstance tests the instance service
func (suite *GetServiceSuite) TestInstance() {
	// Arrange
	inst := "test"
	c, err := NewContainer(WithValue(inst))
	suite.Require().NoError(err)

	// Act
	res, err := c.GetService[string]()

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
	c, err := NewContainer(WithKeyedFactory(key, f))
	suite.Require().NoError(err)

	// Act
	res, err := c.GetKeyedService[string](key)

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
	c, err := NewContainer(WithKeyedFactory(key, f))
	suite.Require().NoError(err)

	// Act
	res, err := c.GetKeyedService[string](key)

	// Assert
	suite.Empty(res)
	suite.Error(err)
}

// TestInstance tests the instance service
func (suite *GetKeyedServiceSuite) TestInstance() {
	// Arrange
	key := "key"
	inst := "test"
	c, err := NewContainer(WithKeyedValue(key, inst))
	suite.Require().NoError(err)

	// Act
	res, err := c.GetKeyedService[string](key)

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
	c, err := NewContainer(WithFactory(f))
	suite.Require().NoError(err)

	// Act & Assert
	var res string
	suite.NotPanics(func() {
		res = c.MustGetService[string]()
	})
	suite.Equal(inst, res)
}

// TestFactoryError tests the factory service error
func (suite *MustGetServiceSuite) TestFactoryError() {
	// Arrange
	f := func() (string, error) {
		return "", errors.ErrUnsupported
	}
	c, err := NewContainer(WithFactory(f))
	suite.Require().NoError(err)

	// Act & Assert
	suite.Panics(func() {
		c.MustGetService[string]()
	})
}

// TestInstance tests the instance service
func (suite *MustGetServiceSuite) TestInstance() {
	// Arrange
	inst := "test"
	c, err := NewContainer(WithValue(inst))
	suite.Require().NoError(err)

	// Act & Assert
	var res string
	suite.NotPanics(func() {
		res = c.MustGetService[string]()
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
	c, err := NewContainer(WithKeyedFactory(key, f))
	suite.Require().NoError(err)

	// Act & Assert
	var res string
	suite.NotPanics(func() {
		res = c.MustGetKeyedService[string](key)
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
	c, err := NewContainer(WithKeyedFactory(key, f))
	suite.Require().NoError(err)

	// Act & Assert
	suite.Panics(func() {
		c.MustGetKeyedService[string](key)
	})
}

// TestInstance tests the instance service
func (suite *MustGetKeyedServiceSuite) TestInstance() {
	// Arrange
	key := "key"
	inst := "test"
	c, err := NewContainer(WithKeyedValue(key, inst))
	suite.Require().NoError(err)

	// Act & Assert
	var res string
	suite.NotPanics(func() {
		res = c.MustGetKeyedService[string](key)
	})
	suite.Equal(inst, res)
}

// TestMustGetKeyedService tests the MustGetKeyedService function
func TestMustGetKeyedService(t *testing.T) {
	suite.Run(t, new(MustGetKeyedServiceSuite))
}

// VerificationAgreementSuite is the suite for testing the verify and getService
// methods reach the same verdict, so the two paths cannot drift
type VerificationAgreementSuite struct {
	suite.Suite
}

// agreementCase is a fixture graph with the dependency type under test
type agreementCase struct {
	// Name is the fixture name
	Name string

	// Opts are the fixture registrations
	Opts []Option

	// DepType is the dependency type both paths are asked about
	DepType reflect.Type
}

// reportsMissing reports whether verification found the dependency missing
func (suite *VerificationAgreementSuite) reportsMissing(c *Container, depType reflect.Type) bool {
	err := c.verify()
	if err == nil {
		return false
	}

	var verErr *VerificationError
	suite.Require().ErrorAs(err, &verErr)

	for _, fault := range verErr.Faults {
		if missing, ok := fault.(*MissingDependencyError); ok && missing.DependencyType == depType {
			return true
		}
	}

	return false
}

// TestAgreement tests both paths reach the same verdict for every fixture
func (suite *VerificationAgreementSuite) TestAgreement() {
	// Arrange
	cases := []agreementCase{
		{
			Name: "registered dependency",
			Opts: []Option{
				WithService[testRepository](newTestRepository),
				WithFactory(newTestService),
			},
			DepType: reflect.TypeFor[testRepository](),
		},
		{
			Name:    "unregistered dependency",
			Opts:    []Option{WithFactory(newTestService)},
			DepType: reflect.TypeFor[testRepository](),
		},
		{
			Name: "keyed only registration",
			Opts: []Option{
				WithKeyedService[testRepository]("db", newTestRepository),
				WithFactory(newTestService),
			},
			DepType: reflect.TypeFor[testRepository](),
		},
		{
			Name: "instance registration",
			Opts: []Option{
				WithValue[testRepository](&testRepositoryImpl{}),
				WithFactory(newTestService),
			},
			DepType: reflect.TypeFor[testRepository](),
		},
		{
			Name: "slice dependency with registrations",
			Opts: []Option{
				WithService[testRepository](newTestRepository),
				WithFactory(newTestCollector),
			},
			DepType: reflect.TypeFor[[]testRepository](),
		},
		{
			Name:    "slice dependency without registrations",
			Opts:    []Option{WithFactory(newTestCollector)},
			DepType: reflect.TypeFor[[]testRepository](),
		},
		{
			Name: "container dependency",
			Opts: []Option{
				WithFactory(func(cont *Container) *testService { return &testService{} }),
			},
			DepType: reflect.TypeFor[*Container](),
		},
	}

	for _, tc := range cases {
		suite.Run(tc.Name, func() {
			// Arrange
			// newContainer skips verification, so the resolution attempt runs
			// even on a graph NewContainer would reject
			c := newContainer(tc.Opts...)

			// Act
			verificationMissing := suite.reportsMissing(c, tc.DepType)

			_, resErr := c.getService(newDependencyIdentifier(tc.DepType))
			resolutionMissing := errors.Is(resErr, ErrServiceNotFound)

			// Assert
			suite.Equal(
				resolutionMissing,
				verificationMissing,
				"verification and resolution disagree on %v",
				tc.DepType,
			)
		})
	}
}

// TestVerificationAgreement tests the verify and getService methods agree
func TestVerificationAgreement(t *testing.T) {
	suite.Run(t, new(VerificationAgreementSuite))
}
