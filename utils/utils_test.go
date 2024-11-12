package utils

import (
	"testing"

	"github.com/alexanderbkl/mytheresa-promotions/models"
	"github.com/stretchr/testify/assert"
)

func TestCalculateDiscount(t *testing.T) {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			discount, finalPrice := CalculateDiscount(tt.product)
			assert.Equal(t, tt.expectedDiscount, discount)
			assert.Equal(t, tt.expectedPrice, finalPrice)
		})
	}
}

func strPointer(s string) *string {
	return &s
}
