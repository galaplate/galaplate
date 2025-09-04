package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/galaplate/galaplate/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestControllerSuite struct {
	tests.TestCase
}

func (t *TestControllerSuite) SetupTest() {
	t.TestCase.SetupTest()
}

func (suite *TestControllerSuite) TestCanAccessHealtCheck() {
	t := suite.T()

	req, err := http.NewRequest("GET", "/api/health", nil)
	assert.NoError(t, err)

	resp, err := suite.App.Test(req)
	suite.NoError(err)
	suite.Equal(200, resp.StatusCode)
}

func (suite *TestControllerSuite) TestErrorFieldNameRequired() {
	t := suite.T()

	payload := strings.NewReader(`{"name": ""}`)
	req, err := http.NewRequest("POST", "/api/test", payload)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.App.Test(req)
	suite.NoError(err)
	suite.Equal(422, resp.StatusCode)
}

func (suite *TestControllerSuite) TestCanCreateTestData() {
	t := suite.T()

	payload := strings.NewReader(`{"name": "this is a test"}`)
	req, err := http.NewRequest("POST", "/api/test", payload)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.App.Test(req)
	suite.NoError(err)
	suite.Equal(201, resp.StatusCode)
}

func (suite *TestControllerSuite) TestGetTestData() {
	t := suite.T()

	req, err := http.NewRequest("GET", "/api/test/123", nil)
	assert.NoError(t, err)

	resp, err := suite.App.Test(req)
	suite.NoError(err)
	suite.Equal(200, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	suite.NoError(err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	suite.NoError(err)

	suite.Equal(true, response["success"])
	data := response["data"].(map[string]interface{})
	suite.Equal("123", data["id"])
	suite.Equal("Test Item", data["name"])
	suite.Equal("This is a test item", data["description"])
}

func (suite *TestControllerSuite) TestCreateTestDataWithDescription() {
	t := suite.T()

	payload := strings.NewReader(`{"name": "test with description", "description": "this is a detailed test"}`)
	req, err := http.NewRequest("POST", "/api/test", payload)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.App.Test(req)
	suite.NoError(err)
	suite.Equal(201, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	suite.NoError(err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	suite.NoError(err)

	suite.Equal(true, response["success"])
	suite.Equal("Test data created successfully", response["message"])
	data := response["data"].(map[string]interface{})
	suite.Equal(float64(1), data["id"])
	suite.Equal("test with description", data["name"])
	suite.Equal("this is a detailed test", data["description"])
}

func (suite *TestControllerSuite) TestInvalidJSONPayload() {
	t := suite.T()

	payload := strings.NewReader(`{"name": "test", "invalid": json}`)
	req, err := http.NewRequest("POST", "/api/test", payload)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.App.Test(req)
	suite.NoError(err)
	suite.Equal(500, resp.StatusCode)
}

func (suite *TestControllerSuite) TestMissingContentType() {
	t := suite.T()

	payload := strings.NewReader(`{"name": "test without content type"}`)
	req, err := http.NewRequest("POST", "/api/test", payload)
	assert.NoError(t, err)

	resp, err := suite.App.Test(req)
	suite.NoError(err)
	suite.Equal(422, resp.StatusCode)
}

func TestControllerSuiteRun(t *testing.T) {
	suite.Run(t, new(TestControllerSuite))
}
