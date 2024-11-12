// handlers/handlers_test.go
package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexanderbkl/mytheresa-promotions/models"
	"github.com/alexanderbkl/mytheresa-promotions/store"
	"github.com/stretchr/testify/assert"
)

// TestProductsHandler tests the ProductsHandler function by mocking the store and sending a request to the handler
func TestProductsHandler(t *testing.T) {
	// Mock products data for testing purposes
	mockProducts := []models.Product{
		{
			SKU:      "000001",
			Name:     "BV Lean leather ankle boots",
			Category: "boots",
			Price:    89000,
		},
		{
			SKU:      "000005",
			Name:     "Nathane leather sneakers",
			Category: "sneakers",
			Price:    59000,
		},
	}

	// Create a mock store with the mock products data
	mockStore := store.NewMockStore(mockProducts)
	// Create a new handler with the mock store
	handler := ProductsHandler(mockStore)

	// Create a new HTTP request to test the handler
	req, err := http.NewRequest("GET", "/products?category=boots", nil)
	// Check if there are any errors in creating the request
	assert.NoError(t, err)

	// Create a new recorder to record the response from the handler
	rr := httptest.NewRecorder()
	// Serve the HTTP request to the handler
	handler.ServeHTTP(rr, req)

	// Check if the status code of the response is as expected
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse the response body into a map of product responses
	var response map[string][]models.ProductResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	// Check if there are any errors in unmarshalling the response
	assert.NoError(t, err)

	// Check if the response contains the expected number of products
	assert.Len(t, response["products"], 1)
	// Check if the SKU of the first product in the response is as expected
	assert.Equal(t, "000001", response["products"][0].SKU)
}
