package tests

import (
	coretesting "github.com/galaplate/core/testing"
	_ "github.com/galaplate/galaplate/db/migrations"
	_ "github.com/galaplate/galaplate/providers"
)

type TestCase struct {
	coretesting.TestCase
}

func (tc *TestCase) SetupSuite() {
	tc.Config = coretesting.DefaultTestConfig()
}

type WithRefreshDatabase struct {
	coretesting.WithRefreshDatabase
}

func (w *WithRefreshDatabase) SetupSuite() {
	w.Config = coretesting.DefaultTestConfig()
	w.Config.RefreshDatabase = true
}

type RefreshDatabaseBeforeEachTest struct {
	coretesting.RefreshDatabaseBeforeEachTest
}

func (r *RefreshDatabaseBeforeEachTest) SetupSuite() {
	r.Config = coretesting.DefaultTestConfig()
	r.Config.RefreshDatabase = true
}
