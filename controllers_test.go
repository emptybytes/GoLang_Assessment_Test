package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Setup a test server
func setupTestServer() *fiber.App {
	app := fiber.New()
	Setup(app)
	return app
}

func TestRegister_ValidInput(t *testing.T) {
	app := setupTestServer()

	// Valid input data
	payload := `{"name": "John Doe", "email": "john@example.com", "password": "password123"}`

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRegister_InvalidInput(t *testing.T) {
	app := setupTestServer()

	// Invalid input data (missing password)
	payload := `{"name": "John Doe", "email": "john@example.com"}`

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// Similar tests can be written for other endpoints, such as Login, AddProduct, etc.

func TestLogin_ValidInput(t *testing.T) {
	app := setupTestServer()

	// Valid input data
	payload := `{"email": "john@example.com", "password": "password123"}`

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// You can also parse the response body and make assertions on its content
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["message"])
	assert.NotNil(t, response["token"])
}

func TestUser_ValidToken(t *testing.T) {
	app := setupTestServer()

	// Assuming you have a valid token
	validToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2V4YW1wbGUuYXV0aDAuY29tLyIsImF1ZCI6Imh0dHBzOi8vYXBpLmV4YW1wbGUuY29tL2NhbGFuZGFyL3YxLyIsInN1YiI6InVzcl8xMjMiLCJpYXQiOjE0NTg3ODU3OTYsImV4cCI6MTQ1ODg3MjE5Nn0.CA7eaHjIHz5NxeIJoFK9krqaeZrPLwmMmgI_XiQiIkQ"

	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	req.Header.Set("Authorization", validToken)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// You can also parse the response body and make assertions on its content
	var user models.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	assert.NoError(t, err)
	assert.NotEmpty(t, user.Name)
	assert.NotEmpty(t, user.Email)
	// Add more assertions based on your User model
}

func TestUser_InvalidToken(t *testing.T) {
	app := setupTestServer()

	// Assuming you have an invalid token
	invalidToken := "inBearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2V4YW1wbGUuYXV0aDAuY29tLyIsImF1ZCI6Imh0dHBzOi8vYXBpLmV4YW1wbGUuY29tL2NhbGFuZGFyL3YxLyIsInN1YiI6InVzcl8xMjMiLCJpYXQiOjE0NTg3ODU3OTYsImV4cCI6MTQ1ODg3MjE5Nn0.CA7eaHjIHz5NxeIJoFK9krqaeZrPLwmMmgI_XiQiIkQ"

	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	req.Header.Set("Authorization", invalidToken)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// You can also parse the response body and make assertions on its content
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "unauthenticated", response["message"])
}

func TestAddProduct_ValidInput(t *testing.T) {
	app := setupTestServer()

	// Valid input data
	payload := `{"name": "Product 1", "description": "Description 1", "price": 20.5}`

	req := httptest.NewRequest(http.MethodPost, "/user/products", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	// Assuming you have a valid token
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2V4YW1wbGUuYXV0aDAuY29tLyIsImF1ZCI6Imh0dHBzOi8vYXBpLmV4YW1wbGUuY29tL2NhbGFuZGFyL3YxLyIsInN1YiI6InVzcl8xMjMiLCJpYXQiOjE0NTg3ODU3OTYsImV4cCI6MTQ1ODg3MjE5Nn0.CA7eaHjIHz5NxeIJoFK9krqaeZrPLwmMmgI_XiQiIkQ")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// You can also parse the response body and make assertions on its content
	var product models.Product
	err = json.NewDecoder(resp.Body).Decode(&product)
	assert.NoError(t, err)
	assert.Equal(t, "Product 1", product.Name)
	assert.Equal(t, "Description 1", product.Description)
	assert.Equal(t, 20.5, product.Price)
}

func TestAddProduct_InvalidInput(t *testing.T) {
	app := setupTestServer()

	// Invalid input data (missing name)
	payload := `{"description": "Description 1", "price": 20.5}`

	req := httptest.NewRequest(http.MethodPost, "/user/products", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	// Assuming you have a valid token
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2V4YW1wbGUuYXV0aDAuY29tLyIsImF1ZCI6Imh0dHBzOi8vYXBpLmV4YW1wbGUuY29tL2NhbGFuZGFyL3YxLyIsInN1YiI6InVzcl8xMjMiLCJpYXQiOjE0NTg3ODU3OTYsImV4cCI6MTQ1ODg3MjE5Nn0.CA7eaHjIHz5NxeIJoFK9krqaeZrPLwmMmgI_XiQiIkQ")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// You can also parse the response body and make assertions on its content
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Name Key Missing", response["message"])
}

func TestGetProductList_ValidToken(t *testing.T) {
	app := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/user/products", nil)
	// Assuming you have a valid token
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2V4YW1wbGUuYXV0aDAuY29tLyIsImF1ZCI6Imh0dHBzOi8vYXBpLmV4YW1wbGUuY29tL2NhbGFuZGFyL3YxLyIsInN1YiI6InVzcl8xMjMiLCJpYXQiOjE0NTg3ODU3OTYsImV4cCI6MTQ1ODg3MjE5Nn0.CA7eaHjIHz5NxeIJoFK9krqaeZrPLwmMmgI_XiQiIkQ")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// You can also parse the response body and make assertions on its content
	var products []models.Product
	err = json.NewDecoder(resp.Body).Decode(&products)
	assert.NoError(t, err)
	// Add more assertions based on your expected product list
}

func TestGetProductList_InvalidToken(t *testing.T) {
	app := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/user/products", nil)
	// Assuming you have an invalid token
	req.Header.Set("Authorization", "inBearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2V4YW1wbGUuYXV0aDAuY29tLyIsImF1ZCI6Imh0dHBzOi8vYXBpLmV4YW1wbGUuY29tL2NhbGFuZGFyL3YxLyIsInN1YiI6InVzcl8xMjMiLCJpYXQiOjE0NTg3ODU3OTYsImV4cCI6MTQ1ODg3MjE5Nn0.CA7eaHjIHz5NxeIJoFK9krqaeZrPLwmMmgI_XiQiIkQ")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// You can also parse the response body and make assertions on its content
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "unauthenticated", response["message"])
}

func TestGetProductById_ExistingProduct(t *testing.T) {
	app := setupTestServer()

	// Assuming you have an existing product ID
	existingProductID := "1"

	req := httptest.NewRequest(http.MethodGet, "/user/products/"+existingProductID, nil)
	// Assuming you have a valid token
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2V4YW1wbGUuYXV0aDAuY29tLyIsImF1ZCI6Imh0dHBzOi8vYXBpLmV4YW1wbGUuY29tL2NhbGFuZGFyL3YxLyIsInN1YiI6InVzcl8xMjMiLCJpYXQiOjE0NTg3ODU3OTYsImV4cCI6MTQ1ODg3MjE5Nn0.CA7eaHjIHz5NxeIJoFK9krqaeZrPLwmMmgI_XiQiIkQ")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// You can also parse the response body and make assertions on its content
	var product models.Product
	err = json.NewDecoder(resp.Body).Decode(&product)
	assert.NoError(t, err)
	// Add more assertions based on your expected product details
}

func TestGetProductById_NonexistentProduct(t *testing.T) {
	app := setupTestServer()

	// Assuming you have a nonexistent product ID
	nonexistentProductID := "999"

	req := httptest.NewRequest(http.MethodGet, "/user/products/"+nonexistentProductID, nil)
	// Assuming you have a valid token
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2V4YW1wbGUuYXV0aDAuY29tLyIsImF1ZCI6Imh0dHBzOi8vYXBpLmV4YW1wbGUuY29tL2NhbGFuZGFyL3YxLyIsInN1YiI6InVzcl8xMjMiLCJpYXQiOjE0NTg3ODU3OTYsImV4cCI6MTQ1ODg3MjE5Nn0.CA7eaHjIHz5NxeIJoFK9krqaeZrPLwmMmgI_XiQiIkQ")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	// You can also parse the response body and make assertions on its content
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Product Not Found", response["message"])
}

func TestUpdateProduct_ValidInput(t *testing.T) {
	app := setupTestServer()

	// Assuming you have an existing product ID
	existingProductID := "1"

	// Valid input data for update
	payload := `{"name": "Updated Product", "description": "Updated Description", "price": 25.0}`

	req := httptest.NewRequest(http.MethodPut, "/user/products/"+existingProductID, strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	// Assuming you have a valid token
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2V4YW1wbGUuYXV0aDAuY29tLyIsImF1ZCI6Imh0dHBzOi8vYXBpLmV4YW1wbGUuY29tL2NhbGFuZGFyL3YxLyIsInN1YiI6InVzcl8xMjMiLCJpYXQiOjE0NTg3ODU3OTYsImV4cCI6MTQ1ODg3MjE5Nn0.CA7eaHjIHz5NxeIJoFK9krqaeZrPLwmMmgI_XiQiIkQ")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// You can also parse the response body and make assertions on its content
	var updatedProduct models.Product
	err = json.NewDecoder(resp.Body).Decode(&updatedProduct)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Product", updatedProduct.Name)
	assert.Equal(t, "Updated Description", updatedProduct.Description)
	assert.Equal(t, 25.0, updatedProduct.Price)
}

func TestUpdateProduct_InvalidInput(t *testing.T) {
	app := setupTestServer()

	// Assuming you have an existing product ID
	existingProductID := "1"

	// Invalid input data (negative price)
	payload := `{"name": "Updated Product", "description": "Updated Description", "price": -5.0}`

	req := httptest.NewRequest(http.MethodPut, "/user/products/"+existingProductID, strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	// Assuming you have a valid token
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2V4YW1wbGUuYXV0aDAuY29tLyIsImF1ZCI6Imh0dHBzOi8vYXBpLmV4YW1wbGUuY29tL2NhbGFuZGFyL3YxLyIsInN1YiI6InVzcl8xMjMiLCJpYXQiOjE0NTg3ODU3OTYsImV4cCI6MTQ1ODg3MjE5Nn0.CA7eaHjIHz5NxeIJoFK9krqaeZrPLwmMmgI_XiQiIkQ")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// You can also parse the response body and make assertions on its content
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid or missing positive Price", response["message"])
}

// Similar tests can be written for other endpoints, such as AddProduct, GetProductList, etc.
