// handlers/handlers_test.go
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexanderbkl/mytheresa-promotions/constants"
	"github.com/alexanderbkl/mytheresa-promotions/models"
	"github.com/alexanderbkl/mytheresa-promotions/store"
	"github.com/stretchr/testify/assert"
)

// TestProductsHandler tests the ProductsHandler function by mocking the store and sending a request to the handler
func TestProductsHandler(t *testing.T) {
	fmt.Println("Starting TestProductsHandler...")

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
	if err != nil {
		fmt.Printf("%s✘ FAIL%s - Error creating request: %v\n", constants.ColorRed, constants.ColorReset, err)
		t.FailNow()
	}
	fmt.Printf("%s✔ PASS%s Created HTTP request - %s\n", constants.ColorGreen, constants.ColorReset, "/products?category=boots")

	// Create a new recorder to record the response from the handler
	rr := httptest.NewRecorder()
	start := time.Now()
	handler.ServeHTTP(rr, req)
	duration := time.Since(start)
	fmt.Printf("%s✔ PASS %sHandler executed - (Time taken: %v)\n", constants.ColorGreen, constants.ColorReset, duration)

	// Check if the status code of the response is as expected
	if assert.Equal(t, http.StatusOK, rr.Code) {
		fmt.Printf("%s✔ PASS %s Status code is OK - %d\n", constants.ColorGreen, constants.ColorReset, rr.Code)
	} else {
		fmt.Printf("%s✘ FAIL%s - Expected status %d, got %d\n", constants.ColorRed, constants.ColorReset, http.StatusOK, rr.Code)
		t.FailNow()
	}

	// Parse the response body into a map of product responses
	var response map[string][]models.ProductResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		fmt.Printf("%s✘ FAIL%s - Error unmarshalling response: %v\n", constants.ColorRed, constants.ColorReset, err)
		t.FailNow()
	}
	fmt.Printf("%s✔ PASS%s Response unmarshalled successfully\n", constants.ColorGreen, constants.ColorReset)

	// Check if the response contains the expected number of products
	if assert.Len(t, response["products"], 1) {
		fmt.Printf("%s✔ PASS %s Correct number of products in response - %d\n", constants.ColorGreen, constants.ColorReset, len(response["products"]))
	} else {
		fmt.Printf("%s✘ FAIL%s - Expected 1 product, got %d\n", constants.ColorRed, constants.ColorReset, len(response["products"]))
		t.FailNow()
	}

	// Check if the SKU of the first product in the response is as expected
	if assert.Equal(t, "000001", response["products"][0].SKU) {
		fmt.Printf("%s✔ PASS%s SKU of first product matches - %s\n", constants.ColorGreen, constants.ColorReset, response["products"][0].SKU)
	} else {
		fmt.Printf("%s✘ FAIL%s - Expected SKU '000001', got %s\n", constants.ColorRed, constants.ColorReset, response["products"][0].SKU)
		t.FailNow()
	}

	fmt.Println("Completed TestProductsHandler.")
}
