package tests

// WithRefreshDatabase is a helper struct that can be embedded in test suites
type WithRefreshDatabase struct {
	TestCase
}

// SetupTest overrides the base SetupTest to enable database refresh
func (w *WithRefreshDatabase) SetupTest() {
	w.enableRefreshDatabase()
	w.TestCase.SetupTest()
}

// SetupSuite can be called to refresh database once per test suite
func (w *WithRefreshDatabase) SetupSuite() {
	w.enableRefreshDatabase()
}

// RefreshDatabaseBeforeEachTest enables refreshing database before each test method
type RefreshDatabaseBeforeEachTest struct {
	TestCase
}

// SetupTest refreshes database before each test method
func (r *RefreshDatabaseBeforeEachTest) SetupTest() {
	r.enableRefreshDatabase()
	r.refreshDatabaseBetweenTests() // Reset flag for each test
	r.TestCase.SetupTest()
}
