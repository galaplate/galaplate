package tests

// DatabaseRefreshingTest is an interface for tests that need database refresh
type DatabaseRefreshingTest interface {
	RefreshDatabase() error
	EnableRefreshDatabase()
}

// WithRefreshDatabase is a helper struct that can be embedded in test suites
// to provide Laravel-style RefreshDatabase functionality
type WithRefreshDatabase struct {
	TestCase
}

// SetupTest overrides the base SetupTest to enable database refresh
func (w *WithRefreshDatabase) SetupTest() {
	w.EnableRefreshDatabase()
	w.TestCase.SetupTest()
}

// SetupSuite can be called to refresh database once per test suite
func (w *WithRefreshDatabase) SetupSuite() {
	w.EnableRefreshDatabase()
}

// RefreshDatabaseBeforeEachTest enables refreshing database before each test method
type RefreshDatabaseBeforeEachTest struct {
	TestCase
}

// SetupTest refreshes database before each test method
func (r *RefreshDatabaseBeforeEachTest) SetupTest() {
	r.EnableRefreshDatabase()
	r.RefreshDatabaseBetweenTests() // Reset flag for each test
	r.TestCase.SetupTest()
}
