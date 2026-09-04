// 🔥 GoNet is the first full-fledged framework made for Golang!
// ⚡️ GoNet is inspired by .NET, NestJS and other languages frameworks
// 🤖 GitHub Repository: https://github.com/akimsavvin/gonet

package di

// testRepository is a dependency interface used by the verification tests
type testRepository interface {
	Get() string
}

// testRepositoryImpl is a testRepository implementation
type testRepositoryImpl struct{}

// Get implements testRepository
func (*testRepositoryImpl) Get() string {
	return "value"
}

// testLogger is a second dependency interface used by the verification tests
type testLogger interface {
	Log(string)
}

// testLoggerImpl is a testLogger implementation
type testLoggerImpl struct{}

// Log implements testLogger
func (*testLoggerImpl) Log(string) {}

// testService is a service depending on a testRepository
type testService struct {
	// Repo is the service repository
	Repo testRepository
}

// testController is a service depending on a testService
type testController struct {
	// Service is the controller service
	Service *testService
}

// testCollector is a service depending on every registered testRepository
type testCollector struct {
	// Repos are all the registered repositories
	Repos []testRepository
}

// testServiceA is a service participating in the verification cycle fixtures
type testServiceA struct {
	// B is the next service of the cycle
	B *testServiceB
}

// testServiceB is a service participating in the verification cycle fixtures
type testServiceB struct {
	// C is the next service of the cycle
	C *testServiceC
}

// testServiceC is a service participating in the verification cycle fixtures
type testServiceC struct {
	// A is the next service of the cycle
	A *testServiceA
}

// testSelfDependent is a service depending on itself
type testSelfDependent struct {
	// Self is the service itself
	Self *testSelfDependent
}

// testKeyedConsumer is a service depending on an unkeyed testRepository, used to
// check a keyed registration does not satisfy it
type testKeyedConsumer struct {
	// Repo is the service repository
	Repo testRepository
}

// testMissingLoggerService is a service depending on an unregistered testLogger
type testMissingLoggerService struct {
	// Logger is the service logger
	Logger testLogger
}

// newTestRepository creates a testRepositoryImpl as a testRepository
func newTestRepository() testRepository {
	return &testRepositoryImpl{}
}

// newTestService creates a testService
func newTestService(repo testRepository) *testService {
	return &testService{Repo: repo}
}

// newTestController creates a testController
func newTestController(service *testService) *testController {
	return &testController{Service: service}
}

// newTestCollector creates a testCollector
func newTestCollector(repos []testRepository) *testCollector {
	return &testCollector{Repos: repos}
}

// testSecondController is a second service depending on a testService, used to
// build a diamond shaped graph
type testSecondController struct {
	// Service is the controller service
	Service *testService
}

// newTestSecondController creates a testSecondController
func newTestSecondController(service *testService) *testSecondController {
	return &testSecondController{Service: service}
}

// testCycleEntry is a service outside a cycle, depending on a testCycleZ.
//
// Its type name sorts before both cycle members, so the walk reaches the cycle
// through it and enters the cycle at testCycleZ rather than at the cycle's
// canonical starting type
type testCycleEntry struct {
	// Z is the cycle member the entry depends on
	Z *testCycleZ
}

// testCycleM is a cycle member whose type name sorts before testCycleZ
type testCycleM struct {
	// Z is the next member of the cycle
	Z *testCycleZ
}

// testCycleZ is a cycle member whose type name sorts after testCycleM
type testCycleZ struct {
	// M is the next member of the cycle
	M *testCycleM
}
