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

type AuthControllerSuite struct {
	tests.RefreshDatabaseBeforeEachTest // Refresh database before each test
}

func (t *AuthControllerSuite) SetupTest() {
	t.RefreshDatabaseBeforeEachTest.SetupTest()
}

func (suite *AuthControllerSuite) TestCanRegisterUser() {
	t := suite.T()

	payload := strings.NewReader(`{"username": "testuser", "email": "test@example.com", "password": "password123"}`)
	req, err := http.NewRequest("POST", "/api/register", payload)
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
	suite.Equal("User registered successfully", response["message"])

	data := response["data"].(map[string]interface{})
	suite.NotNil(data["user"])
	suite.NotEmpty(data["token"])

	user := data["user"].(map[string]interface{})
	suite.Equal("testuser", user["username"])
	suite.Equal("test@example.com", user["email"])
	suite.NotContains(user, "password") // Password should be hidden
}

func (suite *AuthControllerSuite) TestRegisterValidationErrors() {
	t := suite.T()

	// Test missing username
	payload := strings.NewReader(`{"email": "test@example.com", "password": "password123"}`)
	req, err := http.NewRequest("POST", "/api/register", payload)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.App.Test(req)
	suite.NoError(err)
	suite.Equal(400, resp.StatusCode)

	// Test invalid email
	payload = strings.NewReader(`{"username": "testuser", "email": "invalid-email", "password": "password123"}`)
	req, err = http.NewRequest("POST", "/api/register", payload)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err = suite.App.Test(req)
	suite.NoError(err)
	suite.Equal(400, resp.StatusCode)

	// Test short password
	payload = strings.NewReader(`{"username": "testuser", "email": "test@example.com", "password": "123"}`)
	req, err = http.NewRequest("POST", "/api/register", payload)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err = suite.App.Test(req)
	suite.NoError(err)
	suite.Equal(400, resp.StatusCode)
}

func (suite *AuthControllerSuite) TestDuplicateEmailRegistration() {
	t := suite.T()

	// Register first user
	payload := strings.NewReader(`{"username": "testuser1", "email": "test@example.com", "password": "password123"}`)
	req, err := http.NewRequest("POST", "/api/register", payload)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.App.Test(req)
	suite.NoError(err)
	suite.Equal(201, resp.StatusCode)

	// Try to register second user with same email
	payload = strings.NewReader(`{"username": "testuser2", "email": "test@example.com", "password": "password123"}`)
	req, err = http.NewRequest("POST", "/api/register", payload)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err = suite.App.Test(req)
	suite.NoError(err)
	suite.Equal(409, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	suite.NoError(err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	suite.NoError(err)

	suite.Equal(false, response["success"])
	suite.Equal("User with this email already exists", response["message"])
}

func (suite *AuthControllerSuite) TestDuplicateUsernameRegistration() {
	t := suite.T()

	// Register first user
	payload := strings.NewReader(`{"username": "testuser", "email": "test1@example.com", "password": "password123"}`)
	req, err := http.NewRequest("POST", "/api/register", payload)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.App.Test(req)
	suite.NoError(err)
	suite.Equal(201, resp.StatusCode)

	// Try to register second user with same username
	payload = strings.NewReader(`{"username": "testuser", "email": "test2@example.com", "password": "password123"}`)
	req, err = http.NewRequest("POST", "/api/register", payload)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err = suite.App.Test(req)
	suite.NoError(err)
	suite.Equal(409, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	suite.NoError(err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	suite.NoError(err)

	suite.Equal(false, response["success"])
	suite.Equal("User with this username already exists", response["message"])
}

func (suite *AuthControllerSuite) TestCanLoginWithValidCredentials() {
	t := suite.T()

	// First register a user
	registerPayload := strings.NewReader(`{"username": "testuser", "email": "test@example.com", "password": "password123"}`)
	registerReq, err := http.NewRequest("POST", "/api/register", registerPayload)
	assert.NoError(t, err)
	registerReq.Header.Set("Content-Type", "application/json")

	registerResp, err := suite.App.Test(registerReq)
	suite.NoError(err)
	suite.Equal(201, registerResp.StatusCode)

	// Now try to login
	loginPayload := strings.NewReader(`{"email": "test@example.com", "password": "password123"}`)
	loginReq, err := http.NewRequest("POST", "/api/login", loginPayload)
	assert.NoError(t, err)
	loginReq.Header.Set("Content-Type", "application/json")

	loginResp, err := suite.App.Test(loginReq)
	suite.NoError(err)
	suite.Equal(200, loginResp.StatusCode)

	body, err := io.ReadAll(loginResp.Body)
	suite.NoError(err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	suite.NoError(err)

	suite.Equal(true, response["success"])
	suite.Equal("Login successful", response["message"])

	data := response["data"].(map[string]interface{})
	suite.NotNil(data["user"])
	suite.NotEmpty(data["token"])

	user := data["user"].(map[string]interface{})
	suite.Equal("testuser", user["username"])
	suite.Equal("test@example.com", user["email"])
}

func (suite *AuthControllerSuite) TestLoginWithInvalidCredentials() {
	t := suite.T()

	// First register a user
	registerPayload := strings.NewReader(`{"username": "testuser", "email": "test@example.com", "password": "password123"}`)
	registerReq, err := http.NewRequest("POST", "/api/register", registerPayload)
	assert.NoError(t, err)
	registerReq.Header.Set("Content-Type", "application/json")

	registerResp, err := suite.App.Test(registerReq)
	suite.NoError(err)
	suite.Equal(201, registerResp.StatusCode)

	// Try to login with wrong password
	loginPayload := strings.NewReader(`{"email": "test@example.com", "password": "wrongpassword"}`)
	loginReq, err := http.NewRequest("POST", "/api/login", loginPayload)
	assert.NoError(t, err)
	loginReq.Header.Set("Content-Type", "application/json")

	loginResp, err := suite.App.Test(loginReq)
	suite.NoError(err)
	suite.Equal(401, loginResp.StatusCode)

	body, err := io.ReadAll(loginResp.Body)
	suite.NoError(err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	suite.NoError(err)

	suite.Equal(false, response["success"])
	suite.Equal("Invalid credentials", response["message"])
}

func (suite *AuthControllerSuite) TestLoginWithNonExistentUser() {
	t := suite.T()

	loginPayload := strings.NewReader(`{"email": "nonexistent@example.com", "password": "password123"}`)
	loginReq, err := http.NewRequest("POST", "/api/login", loginPayload)
	assert.NoError(t, err)
	loginReq.Header.Set("Content-Type", "application/json")

	loginResp, err := suite.App.Test(loginReq)
	suite.NoError(err)
	suite.Equal(401, loginResp.StatusCode)

	body, err := io.ReadAll(loginResp.Body)
	suite.NoError(err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	suite.NoError(err)

	suite.Equal(false, response["success"])
	suite.Equal("Invalid credentials", response["message"])
}

func (suite *AuthControllerSuite) TestLoginValidationErrors() {
	t := suite.T()

	// Test missing email
	payload := strings.NewReader(`{"password": "password123"}`)
	req, err := http.NewRequest("POST", "/api/login", payload)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.App.Test(req)
	suite.NoError(err)
	suite.Equal(400, resp.StatusCode)

	// Test invalid email format
	payload = strings.NewReader(`{"email": "invalid-email", "password": "password123"}`)
	req, err = http.NewRequest("POST", "/api/login", payload)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err = suite.App.Test(req)
	suite.NoError(err)
	suite.Equal(400, resp.StatusCode)

	// Test short password
	payload = strings.NewReader(`{"email": "test@example.com", "password": "123"}`)
	req, err = http.NewRequest("POST", "/api/login", payload)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err = suite.App.Test(req)
	suite.NoError(err)
	suite.Equal(400, resp.StatusCode)
}

func TestAuthControllerSuiteRun(t *testing.T) {
	suite.Run(t, new(AuthControllerSuite))
}
