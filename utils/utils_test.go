package utils

import (
	"fmt"
	"testing"
	"time"

	"github.com/alexanderbkl/mytheresa-promotions/constants"
	"github.com/alexanderbkl/mytheresa-promotions/models"
	"github.com/stretchr/testify/assert"
)

func TestCalculateDiscount(t *testing.T) {
	// Define test cases
	tests := []struct {
		name             string
		product          models.Product
		expectedDiscount *string
		expectedPrice    int
	}{
		{
			name: "No discount",
			product: models.Product{
				SKU:      "000005",
				Name:     "Nathane leather sneakers",
				Category: "sneakers",
				Price:    59000,
			},
			expectedDiscount: nil,
			expectedPrice:    59000,
		},
		{
			name: "Category discount",
			product: models.Product{
				SKU:      "000001",
				Name:     "BV Lean leather ankle boots",
				Category: "boots",
				Price:    89000,
			},
			expectedDiscount: strPointer("30%"),
			expectedPrice:    62300,
		},
		{
			name: "SKU discount",
			product: models.Product{
				SKU:      "000003",
				Name:     "Ashlington leather ankle boots",
				Category: "boots",
				Price:    71000,
			},
			expectedDiscount: strPointer("30%"),
			expectedPrice:    49700,
		},
	}

	fmt.Println("Starting TestCalculateDiscount...")
	for _, tt := range tests {
		// Run each test case
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()
			discount, finalPrice := CalculateDiscount(tt.product)
			// Check if the discount and final price match the expected values
			if assert.Equal(t, tt.expectedDiscount, discount) && assert.Equal(t, tt.expectedPrice, finalPrice) {
				duration := time.Since(start)
				fmt.Printf("%s✔ PASS%s - %s (Time taken: %v)\n", constants.ColorGreen, constants.ColorReset, tt.name, duration)
			} else {
				duration := time.Since(start)
				fmt.Printf("%s✘ FAIL%s - %s (Time taken: %v)\n", constants.ColorRed, constants.ColorReset, tt.name, duration)
			}
		})
	}
	fmt.Println("Completed TestCalculateDiscount.")
}

func strPointer(s string) *string {
	return &s
}
