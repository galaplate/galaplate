package tests

import (
	coretesting "github.com/galaplate/core/testing"
	"github.com/galaplate/galaplate/router"
)

type TestCase struct {
	coretesting.TestCase
}

func (tc *TestCase) SetupSuite() {
	tc.Config = coretesting.DefaultTestConfig()
	tc.Config.SetupRoutes = router.SetupRouter
}

type WithRefreshDatabase struct {
	coretesting.WithRefreshDatabase
}

func (w *WithRefreshDatabase) SetupSuite() {
	w.Config = coretesting.DefaultTestConfig()
	w.Config.SetupRoutes = router.SetupRouter
	w.Config.RefreshDatabase = true
}

type RefreshDatabaseBeforeEachTest struct {
	coretesting.RefreshDatabaseBeforeEachTest
}

func (r *RefreshDatabaseBeforeEachTest) SetupSuite() {
	r.Config = coretesting.DefaultTestConfig()
	r.Config.SetupRoutes = router.SetupRouter
	r.Config.RefreshDatabase = true
}
