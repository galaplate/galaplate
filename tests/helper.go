package tests

import coretesting "github.com/galaplate/core/testing"

func NewHTTPTestHelper(tc *TestCase) *coretesting.HTTPTestHelper {
	return coretesting.NewHTTPTestHelper(&tc.TestCase)
}

func NewAssertHelper(tc *TestCase) *coretesting.AssertHelper {
	return coretesting.NewAssertHelper(&tc.TestCase)
}

func NewDatabaseHelper(tc *TestCase) *coretesting.DatabaseHelper {
	return coretesting.NewDatabaseHelper(&tc.TestCase)
}
