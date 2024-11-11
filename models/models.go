// models/models.go
package models

type Product struct {
	SKU      string `json:"sku"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Price    int    `json:"price"`
}

type ProductResponse struct {
	SKU      string      `json:"sku"`
	Name     string      `json:"name"`
	Category string      `json:"category"`
	Price    PriceDetail `json:"price"`
}

type PriceDetail struct {
	Original           int     `json:"original"`
	Final              int     `json:"final"`
	DiscountPercentage *string `json:"discount_percentage"`
	Currency           string  `json:"currency"`
}
