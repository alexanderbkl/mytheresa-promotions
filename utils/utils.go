// utils/utils.go
package utils

import (
	"fmt"

	"github.com/alexanderbkl/mytheresa-promotions/models"
)

func CalculateDiscount(product models.Product) (*string, int) {
	var discount int

	if product.SKU == "000003" {
		discount = 15
	}
	if product.Category == "boots" && discount < 30 {
		discount = 30
	}

	if discount > 0 {
		discountPercentage := fmt.Sprintf("%d%%", discount)
		finalPrice := product.Price * (100 - discount) / 100
		return &discountPercentage, finalPrice
	}

	return nil, product.Price
}
